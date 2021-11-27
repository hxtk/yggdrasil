package authn

import (
	"context"
	"fmt"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type key int

var identityKey key

type Identity struct {
	DisplayName string
	UserType    string
	ID          uuid.UUID
	Subject     *pb.SubjectReference
}

func IdentityFromContext(ctx context.Context) (*Identity, error) {
	id, ok := ctx.Value(identityKey).(*Identity)
	if !ok {
		return nil, fmt.Errorf("no identity found")
	}
	return id, nil
}

// TLSAuth is a grpc_auth.AuthFunc for authenticating clients with mutual TLS.
//
// In order to be used, a client certificate must have a URI SAN which refers to
// their object reference in the authorization database. For example:
//
//     discord/users:0cf4934c-8583-4406-a9d8-b2ee88a8c0f9
//     discord/bots:0cf4934c-8583-4406-a9d8-b2ee88a8c0f9
//     users:0cf4934c-8583-4406-a9d8-b2ee88a8c0f9
//     service-accounts:0cf4934c-8583-4406-a9d8-b2ee88a8c0f9
//
// While it is recommended that Client IDs be UUIDs, this is not a requirement.
// For detailed documentation of the specifications for these fields, see [1]
// and understand that the expected format is `object_type ':' object_id`, i.e.,
// the two fields concatenated together with a colon (`:`).
//
// The Certificate's Common Name is used as a display name
//
// [1]: https://github.com/authzed/api/blob/main/authzed/api/v1/core.proto#L38
func TLSAuth(ctx context.Context) (context.Context, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "no peer found")
	}

	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		log.Info("No TLSInfo found")
		return ctx, nil
	}

	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		log.Info("No verified chains found")
		return ctx, nil
	}

	cert := tlsAuth.State.VerifiedChains[0][0]
	userName := cert.Subject.CommonName

	subject, err := getUserURI(cert.URIs)
	if err != nil {
		log.WithError(err).
			WithField("commonName", userName).
			Info("Client certificate unusable.")
		return ctx, nil
	}
	return context.WithValue(ctx, identityKey, &Identity{
		Name:    userName,
		Subject: subject,
	}), nil
}

var _ grpc_auth.AuthFunc = TLSAuth

func init() {
	objectRefRegex = regexp.MustCompile("^([a-z][a-z0-9_]{2,61}[a-z0-9]/[a-z][a-z0-9_]{2,62}[a-z0-9]):([a-zA-Z0-9_][a-zA-Z0-9/_-]{0,127})$")
}

var objectRefRegex *regexp.Regexp

func getUserURI(uris []*url.URL) (*pb.SubjectReference, error) {
	for _, v := range uris {
		if v.Opaque == "" {
			if matches := objectRefRegex.FindSubMatch([]byte(v.Path)); matches != nil {
				return &pb.SubjectReference{
					Object: &pb.ObjectReference{
						ObjectType: string(matches[1]),
						ObjectId:   string(matches[2]),
					},
				}, nil
			}
			continue
		}

		return &pb.SubjectReference{
			Object: &pb.ObjectReference{
				ObjectType: v.Scheme,
				ObjectId:   v.Opaque,
			},
		}, nil
	}
	return nil, fmt.Errorf("no client URI found")
}
