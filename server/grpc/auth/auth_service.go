package auth

import (
	"context"
	"encoding/json"
	"github.com/bufbuild/connect-go"
	"github.com/cosmos/btcutil/bech32"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/auth/authpb"
	"github.com/loomi-labs/star-scope/grpc/auth/authpb/authpbconnect"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//goland:noinspection GoNameStartsWithPackageName
type AuthService struct {
	authpbconnect.UnimplementedAuthServiceHandler
	userManager          *database.UserManager
	chainManager         *database.ChainManager
	eventListenerManager *database.EventListenerManager
	jwtManager           *JWTManager
}

func NewAuthServiceHandler(dbManagers *database.DbManagers, jwtManager *JWTManager) authpbconnect.AuthServiceHandler {
	return &AuthService{
		userManager:          dbManagers.UserManager,
		chainManager:         dbManagers.ChainManager,
		eventListenerManager: dbManagers.EventListenerManager,
		jwtManager:           jwtManager,
	}
}

var (
	ErrorLoginFailed             = status.Error(codes.Unauthenticated, "login failed")
	ErrorUserNotFound            = status.Error(codes.NotFound, "user not found")
	ErrorInternal                = status.Error(codes.Internal, "internal error")
	ErrorTokenVerificationFailed = status.Error(codes.Unauthenticated, "token verification failed")
)

func verifySignature(message string) bool {
	var keplrResponse map[string]interface{}
	err := json.Unmarshal([]byte(message), &keplrResponse)
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

type KeplrResponse struct {
	Signed struct {
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
	} `json:"signed"`
	Signature struct {
		PubKey struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"pub_key"`
		Signature string `json:"signature"`
	} `json:"signature"`
}

func getWalletAddress(message string) (string, error) {
	var keplrResponse KeplrResponse
	err := json.Unmarshal([]byte(message), &keplrResponse)
	if err != nil {
		return "", err
	}
	return keplrResponse.Signed.Msgs[0].Value.Signer, nil
}

func (s *AuthService) KeplrLogin(ctx context.Context, request *connect.Request[authpb.KeplrLoginRequest]) (*connect.Response[authpb.LoginResponse], error) {
	if !verifySignature(request.Msg.GetKeplrResponse()) {
		return nil, ErrorLoginFailed
	}

	address, err := getWalletAddress(request.Msg.GetKeplrResponse())
	if err != nil {
		log.Sugar.Errorf("error while getting wallet address: %v", err)
		return nil, ErrorLoginFailed
	}

	user, err := s.userManager.QueryByWalletAddress(ctx, address)
	if err != nil && ent.IsNotFound(err) {
		user = s.userManager.CreateOrUpdate(ctx, address, address)

		hrp, _, err := bech32.Decode(address, 1023)
		if err != nil {
			log.Sugar.Errorf("error while decoding bech32: %v", err)
			return nil, ErrorLoginFailed
		}

		chain, err := s.chainManager.QueryByBech32Prefix(ctx, hrp)
		if err != nil {
			log.Sugar.Errorf("error while querying chain by bech32 prefix: %v", err)
			return nil, ErrorLoginFailed
		}
		_, err = s.eventListenerManager.Create(ctx, user, chain, address)
		if err != nil {
			log.Sugar.Errorf("error while creating event listener: %v", err)
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
	return connect.NewResponse(&authpb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}), nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, request *connect.Request[authpb.RefreshAccessTokenRequest]) (*connect.Response[authpb.RefreshAccessTokenResponse], error) {
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

	return connect.NewResponse(&authpb.RefreshAccessTokenResponse{AccessToken: accessToken}), nil
}
