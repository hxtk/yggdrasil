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
	Name     string
	UserType string
	ID       uuid.UUID
	Subject  *pb.SubjectReference
}

func IdentityFromContext(ctx context.Context) (*Identity, error) {
	id, ok := ctx.Value(identityKey).(*Identity)
	if !ok {
		return nil, fmt.Errorf("no identity found")
	}
	return id, nil
}

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
