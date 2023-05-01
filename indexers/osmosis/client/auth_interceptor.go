package client

import (
	"context"
	"github.com/bufbuild/connect-go"
)

//goland:noinspection GoNameStartsWithPackageName
type AuthInterceptor struct {
	authToken string
}

func NewAuthInterceptor(authToken string) *AuthInterceptor {
	return &AuthInterceptor{authToken: authToken}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		req.Header().Set("Authorization", i.authToken)
		return next(ctx, req)
	}
}

func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(
		ctx context.Context,
		spec connect.Spec,
	) connect.StreamingClientConn {
		conn := next(ctx, spec)
		conn.RequestHeader().Set("Authorization", i.authToken)
		return conn
	}
}

func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		conn.RequestHeader().Set("Authorization", i.authToken)
		return next(ctx, conn)
	}
}
