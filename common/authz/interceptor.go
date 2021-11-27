package authz

import (
	"context"
	"fmt"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hxtk/yggdrasil/common/authn"
	"github.com/hxtk/yggdrasil/common/authz/v1alpha1"
)

type key int

var zedTokenKey key

type Registrar interface {
	RegisterPermissions(perms map[string]*v1alpha1.PermissionsRule)
}

// ResourceAuthorizer is a resource-oriented authorization interceptor.
//
// It is designed to work with resource-oriented gRPC APIs compliant with
// Google's Resource-Oriente API Design Guide [1]. In particular, request
// types intended to be used with this interceptor MUST include either a
// `name` or `parent` field. Operations on existing resources SHOULD
// use the `name` field to uniquely identify the specific operand resource,
// but operations which create a resource MAY alternatively use a `parent`
// field to identify the parent resource. If both fields are present then
// the `name` field will take precedence.
//
// The permissions for a particular resource and endpoint are determined
// by annotations in the protobuf which may be found in the `v1alpha1`
// subpackage.
//
// These annotations specify the namespace or type of the resource
// and the permission one needs in order to use it as an operand to the
// gRPC method or service on which the annotation is specified.
//
// The authorization is evaluated using AuthZed SpiceDB [1], which
// has an architecture heavily inspired by Google's Zanzibar API [3].
//
// [1]: https://cloud.google.com/apis/design/
// [2]: https://docs.authzed.com/
// [3]: https://research.google/pubs/pub48190/
type ResourceAuthorizer struct {
	authzClient    pb.PermissionsServiceClient
	rpcPermissions map[string]*v1alpha1.PermissionsRule
}

func (ra *ResourceAuthorizer) RegisterPermissions(perms map[string]*v1alpha1.PermissionsRule) {
	if ra.rpcPermissions == nil {
		ra.rpcPermissions = perms
		return
	}

	for k, v := range perms {
		ra.rpcPermissions[k] = v
	}
}

// UnaryServerInterceptor returns an interceptor for checking authorization.
//
// Authorization checks are tied to a particular moment in time in order to
// balance the concerns of minimizing latency while ensuring that users are
// unable to access resources newer than their most recent authorization
// decision. This is Zanzibar's answer to extracting strongly-consistent
// authorization policy enforcement from a weakly-consistent distributed
// data store.
//
// SpiceDB uses "ZedTokens" to indicate the point in time at which a
// decision is valid [1]. This is an opaque token that, for optimal latency,
// SHOULD be cached and reused where practical. To that end, it is embedded
// in the request context and may be used by the request implementation
// in order to effect more fine-grained authorization policy decisions.
//
// [1]: https://docs.authzed.com/reference/api-consistency
func (ra *ResourceAuthorizer) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		identity, err := authn.IdentityFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "No client identity found.")
		}

		perms, ok := ra.rpcPermissions[info.FullMethod]
		if !ok {
			return nil, status.Error(codes.Internal, "Could not determine required permissions.")
		}

		name, err := getResourceName(req)
		if err != nil {
			return nil, err
		}

		decision, err := ra.authzClient.CheckPermission(ctx, &pb.CheckPermissionRequest{
			Resource: &pb.ObjectReference{
				ObjectType: perms.GetResourceType(),
				ObjectId:   name,
			},
			Permission: perms.GetPermission(),
			Subject: &pb.SubjectReference{
				Object: &pb.ObjectReference{
					ObjectType: "users",
					ObjectId:   identity.UserID.String(),
				},
			},
		})

		if err != nil {
			return nil, status.Error(codes.Unavailable, "Permission check could not be processed.")
		}

		if decision.GetPermissionship() != pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
			return nil, status.Error(codes.PermissionDenied, "User does not have permission to perform this action.")
		}

		return handler(context.WithValue(ctx, zedTokenKey, decision.CheckedAt), req)
	}
}

type namer interface {
	GetName() string
}

type parenter interface {
	GetParent() string
}

func getResourceName(req interface{}) (string, error) {
	switch v := req.(type) {
	case namer:
		return v.GetName(), nil
	case parenter:
		return v.GetParent(), nil
	default:
		return "", status.Error(codes.Internal, "Could not determine resource name from request.")
	}
}

// ZedTokenFromContext returns a ZedToken embedded in the context.
//
// A ZedToken is embedded in the context at resource authorization time by the
// interceptor function, and should be reused in downstream authorization checks
// for the same request in order to minimize latency while ensuring consistent
// results.
func ZedTokenFromContext(ctx context.Context) (*pb.ZedToken, error) {
	zt, ok := ctx.Value(zedTokenKey).(*pb.ZedToken)
	if !ok {
		return nil, fmt.Errorf("no ZedToken was found in this context")
	}
	return zt, nil
}
