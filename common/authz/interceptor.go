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

type namer interface {
	GetName() string
}

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

		if _, ok := req.(namer); !ok {
			return nil, status.Error(codes.Internal, "Could not determine resource name from request.")
		}

		decision, err := ra.authzClient.CheckPermission(ctx, &pb.CheckPermissionRequest{
			Resource: &pb.ObjectReference{
				ObjectType: perms.GetResourceType(),
				ObjectId:   req.(namer).GetName(),
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
			return nil, status.Error(codes.Unavailable, "Permission check failed.")
		}

		if decision.GetPermissionship() != pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
			return nil, status.Error(codes.PermissionDenied, "User does not have permission to perform this action.")
		}

		return handler(context.WithValue(ctx, zedTokenKey, decision.CheckedAt), req)
	}
}

func ZedTokenFromContext(ctx context.Context) (*pb.ZedToken, error) {
	zt, ok := ctx.Value(zedTokenKey).(*pb.ZedToken)
	if !ok {
		return nil, fmt.Errorf("no ZedToken was found in this context")
	}
	return zt, nil
}
