package auth

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

//goland:noinspection GoNameStartsWithPackageName
type AuthInterceptor struct {
	jwtManager      *JWTManager
	userManager     *database.UserManager
	accessibleRoles map[string][]Role
	authToken       string
}

func NewAuthInterceptor(jwtManager *JWTManager, userManager *database.UserManager, accessibleRoles map[string][]Role, authToken string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager, accessibleRoles: accessibleRoles, userManager: userManager, authToken: authToken}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		debugInfo := "--> unary interceptor: " + req.Spec().Procedure

		ctx, err := i.authorize(ctx, req.Header(), req.Spec().Procedure)
		if err != nil {
			log.Sugar.Debug(debugInfo + " access denied!")
			return nil, err
		}
		log.Sugar.Debug(debugInfo)
		return next(ctx, req)
	}
}

func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(
		ctx context.Context,
		spec connect.Spec,
	) connect.StreamingClientConn {
		conn := next(ctx, spec)
		debugInfo := "--> stream client interceptor: " + conn.Spec().Procedure
		log.Sugar.Debug(debugInfo)
		return conn
	}
}

func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		debugInfo := "--> stream handler interceptor: " + conn.Spec().Procedure

		ctx, err := i.authorize(ctx, conn.RequestHeader(), conn.Spec().Procedure)
		if err != nil {
			log.Sugar.Debug(debugInfo + " access denied!")
			return err
		}
		//wrapped := grpcmiddleware.WrapServerStream(stream)
		//wrapped.WrappedContext = ctx
		log.Sugar.Debug(debugInfo)
		return next(ctx, conn)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, header http.Header, procedure string) (context.Context, error) {
	accessibleRoles, ok := i.accessibleRoles[procedure]
	if !ok {
		log.Sugar.Errorf("no access roles are defined for %s", procedure)
		return ctx, status.Errorf(codes.Unauthenticated, "unauthorized access")
	}
	if slices.Contains(accessibleRoles, Unauthenticated) {
		return ctx, nil
	}

	authorization := header.Get("Authorization")
	if authorization == "" {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	if authorization == i.authToken {
		return ctx, nil
	}

	accessToken := strings.Replace(authorization, "Bearer ", "", 1)
	claims, err := i.jwtManager.Verify(accessToken)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			entUser, err := i.userManager.QueryById(ctx, claims.UserId)
			if err != nil {
				return nil, status.Error(codes.Internal, "user not found")
			}
			return context.WithValue(ctx, common.ContextKeyUser, entUser), nil
		}
	}

	return ctx, status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
