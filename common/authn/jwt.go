package main

import (
	"context"
	"crypto/rsa"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTAuth struct {
	keys []*rsa.PublicKey
}

func (a *JWTAuth) Authenticate(ctx context.Context) (context.Context, error) {
	tokenStr, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Error validating authorization.")
	}

	for _, key := range a.keys {
		token, err = jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
			return key, nil
		})
	}
}
