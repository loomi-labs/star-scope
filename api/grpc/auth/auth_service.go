package auth

import (
	"context"
	"encoding/json"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/shifty11/blocklog-backend/database"
	"github.com/shifty11/blocklog-backend/ent"
	pb "github.com/shifty11/blocklog-backend/grpc/auth/v1"
	authconnect "github.com/shifty11/blocklog-backend/grpc/auth/v1/v1connect"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//goland:noinspection GoNameStartsWithPackageName
type AuthService struct {
	authconnect.UnimplementedAuthServiceHandler
	userManager    *database.UserManager
	projectManager *database.ProjectManager
	jwtManager     *JWTManager
}

func NewAuthServiceHandler(dbManagers *database.DbManagers, jwtManager *JWTManager) authconnect.AuthServiceHandler {
	return &AuthService{
		userManager:    dbManagers.UserManager,
		projectManager: dbManagers.ProjectManager,
		jwtManager:     jwtManager,
	}
}

var (
	ErrorLoginFailed             = status.Error(codes.Unauthenticated, "login failed")
	ErrorUserNotFound            = status.Error(codes.NotFound, "user not found")
	ErrorInternal                = status.Error(codes.Internal, "internal error")
	ErrorTokenVerificationFailed = status.Error(codes.Unauthenticated, "token verification failed")
)

func verifySignature(request *pb.KeplrLoginRequest) bool {
	var signature map[string]interface{}
	err := json.Unmarshal([]byte(request.GetSignature()), &signature)
	if err != nil {
		return false
	}

	var signedMessage map[string]interface{}
	err = json.Unmarshal([]byte(request.GetSignedMessage()), &signedMessage)
	if err != nil {
		return false
	}
	// TODO: make a proper verification
	return true

	//pubkeyBytes, err := hex.DecodeString(signature["signature"].(map[string]interface{})["pub_key"].(string))
	//if err != nil {
	//	return false
	//}
	//log.Sugar.Infof("pubkeyBytes: %v", pubkeyBytes)
	//
	//// Create a public key object from the bytes
	////pubKey, err := secp256k1.PU
	////if err != nil {
	////	return false
	////}
	//
	//signatureBytes, err := hex.DecodeString(signature["value"].(string))
	//log.Sugar.Infof("signatureBytes: %v", signatureBytes)
	//return false
}

type SignedMessage struct {
	AccountNumber string `json:"account_number"`
	ChainID       string `json:"chain_id"`
	Fee           struct {
		Amount []any  `json:"amount"`
		Gas    string `json:"gas"`
	} `json:"fee"`
	Memo string `json:"memo"`
	Msgs []struct {
		Type  string `json:"type"`
		Value struct {
			Data   string `json:"data"`
			Signer string `json:"signer"`
		} `json:"value"`
	} `json:"msgs"`
	Sequence string `json:"sequence"`
}

func getWalletAddress(message string) (string, error) {
	var signedMessage SignedMessage
	err := json.Unmarshal([]byte(message), &signedMessage)
	if err != nil {
		return "", err
	}
	return signedMessage.Msgs[0].Value.Signer, nil
}

func (s *AuthService) KeplrLogin(ctx context.Context, request *connect_go.Request[pb.KeplrLoginRequest]) (*connect_go.Response[pb.LoginResponse], error) {
	verifySignature(request.Msg)

	address, err := getWalletAddress(request.Msg.SignedMessage)
	if err != nil {
		log.Sugar.Errorf("error while getting wallet address: %v", err)
		return nil, ErrorLoginFailed
	}

	user, err := s.userManager.QueryByWalletAddress(ctx, address)
	if err != nil && ent.IsNotFound(err) {
		user = s.userManager.CreateOrUpdate(ctx, address, address)
		_, err := s.projectManager.CreateCosmosProject(ctx, user, address)
		if err != nil {
			log.Sugar.Errorf("error while creating cosmos project: %v", err)
			return nil, ErrorLoginFailed
		}
	} else if err != nil {
		log.Sugar.Errorf("error while querying user by wallet address: %v", err)
		return nil, ErrorLoginFailed
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
	return connect_go.NewResponse(&pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}), nil
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
