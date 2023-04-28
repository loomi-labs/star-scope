package auth

import (
	"context"
	"encoding/hex"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/shifty11/blocklog-backend/database"
	pb "github.com/shifty11/blocklog-backend/grpc/auth/v1"
	authconnect "github.com/shifty11/blocklog-backend/grpc/auth/v1/v1connect"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//goland:noinspection GoNameStartsWithPackageName
type AuthService struct {
	authconnect.UnimplementedAuthServiceHandler
	userManager *database.UserManager
	jwtManager  *JWTManager
}

func NewAuthServiceHandler(userManager *database.UserManager, jwtManager *JWTManager) authconnect.AuthServiceHandler {
	return &AuthService{
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

func verifySignature(msg []byte, sig []byte, pubkey string) bool {
	// Parse the public key from the hex-encoded string
	pubkeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return false
	}

	// Create a public key object from the bytes
	pubKey, err := secp256k1.PubKeyFromBytes(pubkeyBytes)
	if err != nil {
		return false
	}

	// Create a signature object from the bytes
	sigObj := ed25519.Signature(sig)

	// Verify the signature using the public key
	return pubKey.VerifySignature(msg, sigObj)
}

func (s *AuthService) KeplrLogin(ctx context.Context, request *connect_go.Request[pb.KeplrLoginRequest]) (*connect_go.Response[pb.LoginResponse], error) {
	// TODO: verify signature

	return connect_go.NewResponse(&pb.LoginResponse{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}), nil

	//user, err := s.userManager.QueryByWalletAddress(ctx, request.Msg.GetWalletAddress())
	//if err != nil {
	//	return nil, ErrorUserNotFound
	//}
	//
	//accessToken, err := s.jwtManager.GenerateToken(user, AccessToken)
	//if err != nil {
	//	log.Sugar.Errorf("Could not generate accessToken for user %v (%v): %v", user.Name, user.ID, err)
	//	return nil, ErrorLoginFailed
	//}
	//
	//refreshToken, err := s.jwtManager.GenerateToken(user, RefreshToken)
	//if err != nil {
	//	log.Sugar.Errorf("Could not generate refreshToken for user %v (%v): %v", user.Name, user.ID, err)
	//	return nil, ErrorInternal
	//}
	//return connect_go.NewResponse(&pb.LoginResponse{
	//	AccessToken:  accessToken,
	//	RefreshToken: refreshToken,
	//}), nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, request *connect_go.Request[pb.RefreshAccessTokenRequest]) (*connect_go.Response[pb.RefreshAccessTokenResponse], error) {
	claims, err := s.jwtManager.Verify(request.Msg.GetRefreshToken())
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

	return connect_go.NewResponse(&pb.RefreshAccessTokenResponse{AccessToken: accessToken}), nil
}
