package auth

import (
	"context"
	"encoding/json"
	"github.com/shifty11/blocklog-backend/database"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//goland:noinspection GoNameStartsWithPackageName
type AuthServer struct {
	UnimplementedAuthServiceServer
	userManager *database.UserManager
	jwtManager  *JWTManager
}

type DiscordIdentity struct {
	ID       json.Number `json:"id"`
	Username string      `json:"username"`
}

func NewAuthServer(
	userManager *database.UserManager,
	jwtManager *JWTManager,
) AuthServiceServer {
	return &AuthServer{
		userManager: userManager,
		jwtManager:  jwtManager,
	}
}

var (
	ErrorLoginFailed             = status.Error(codes.Unauthenticated, "login failed")
	ErrorUserNotFound            = status.Error(codes.NotFound, "user not found")
	ErrorInternal                = status.Error(codes.Internal, "internal error")
	ErrorTokenVerificationFailed = status.Error(codes.Unauthenticated, "token verification failed")
)

func (s *AuthServer) KeplrLogin(ctx context.Context, request *KeplrLoginRequest) (*LoginResponse, error) {
	// TODO: verify signature

	user, err := s.userManager.QueryByWalletAddress(ctx, request.GetWalletAddress())
	if err != nil {
		return nil, ErrorUserNotFound
	}

	accessToken, err := s.jwtManager.GenerateToken(user, AccessToken)
	if err != nil {
		log.Sugar.Errorf("Could not generate accessToken for user %v (%v): %v", user.Name, user.ID, err)
		return nil, ErrorLoginFailed
	}

	refreshToken, err := s.jwtManager.GenerateToken(user, RefreshToken)
	if err != nil {
		log.Sugar.Errorf("Could not generate refreshToken for user %v (%v): %v", user.Name, user.ID, err)
		return nil, ErrorInternal
	}
	return &LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthServer) RefreshAccessToken(ctx context.Context, request *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error) {
	claims, err := s.jwtManager.Verify(request.RefreshToken)
	if err != nil {
		return nil, ErrorTokenVerificationFailed
	}

	entUser, err := s.userManager.QueryById(ctx, claims.UserId)
	if err != nil {
		log.Sugar.Errorf("Could not find user %v (%v): %v", claims.UserId, claims.UserId, err)
		return nil, ErrorUserNotFound
	}

	accessToken, err := s.jwtManager.GenerateToken(entUser, AccessToken)
	if err != nil {
		log.Sugar.Errorf("Could not generate accessToken for user %v (%v): %v", entUser.Name, entUser.ID, err)
		return nil, ErrorInternal
	}

	return &RefreshAccessTokenResponse{AccessToken: accessToken}, nil
}
