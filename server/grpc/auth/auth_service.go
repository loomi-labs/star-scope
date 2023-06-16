package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/auth/authpb"
	"github.com/loomi-labs/star-scope/grpc/auth/authpb/authpbconnect"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type AuthService struct {
	authpbconnect.UnimplementedAuthServiceHandler
	userManager          *database.UserManager
	chainManager         *database.ChainManager
	eventListenerManager *database.EventListenerManager
	jwtManager           *JWTManager
	kafkaInternal        kafka_internal.KafkaInternal
	telegramToken        string
	discordOAuth2Config  *oauth2.Config
}

func NewAuthServiceHandler(
	dbManagers *database.DbManagers,
	jwtManager *JWTManager,
	kafkaInternal kafka_internal.KafkaInternal,
	telegramToken string,
	discordOAuth2Config *oauth2.Config,
) authpbconnect.AuthServiceHandler {
	return &AuthService{
		userManager:          dbManagers.UserManager,
		chainManager:         dbManagers.ChainManager,
		eventListenerManager: dbManagers.EventListenerManager,
		jwtManager:           jwtManager,
		kafkaInternal:        kafkaInternal,
		telegramToken:        telegramToken,
		discordOAuth2Config:  discordOAuth2Config,
	}
}

var (
	ErrorLoginFailed             = status.Error(codes.Unauthenticated, "login failed")
	ErrorLoginExpired            = status.Error(codes.Unauthenticated, "login credentials expired")
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

func (s *AuthService) login(user *ent.User) (*connect.Response[authpb.LoginResponse], error) {
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

func (s *AuthService) KeplrLogin(ctx context.Context, request *connect.Request[authpb.KeplrLoginRequest]) (*connect.Response[authpb.LoginResponse], error) {
	if !verifySignature(request.Msg.GetKeplrResponse()) {
		return nil, ErrorLoginFailed
	}

	walletAddress, err := getWalletAddress(request.Msg.GetKeplrResponse())
	if err != nil {
		log.Sugar.Errorf("error while getting wallet address: %v", err)
		return nil, ErrorLoginFailed
	}

	user, err := s.userManager.QueryByWalletAddress(ctx, walletAddress)
	if err != nil && ent.IsNotFound(err) {
		err := s.userManager.WithTx(ctx, func(tx *ent.Tx) error {
			user, err = s.userManager.CreateByWalletAddress(ctx, tx, walletAddress)
			if err != nil {
				return err
			}
			chains := s.chainManager.QueryEnabled(ctx)
			els, err := s.eventListenerManager.CreateBulk(ctx, tx, user, chains, walletAddress)
			if err != nil {
				return err
			}
			go NewSetupCrawler(s.kafkaInternal).fetchUnstakingEvents(els)
			return nil
		})
		if err != nil {
			log.Sugar.Errorf("error while creating user by wallet address: %v", err)
			return nil, ErrorLoginFailed
		}
	} else if err != nil {
		log.Sugar.Errorf("error while querying user by wallet address: %v", err)
		return nil, ErrorLoginFailed
	}

	return s.login(user)
}

func (s *AuthService) secretKey1() []byte {
	sha := sha256.New()
	sha.Write([]byte(s.telegramToken))
	secretKey := sha.Sum(nil)
	return secretKey
}

func (s *AuthService) secretKey2() []byte {
	h1 := hmac.New(sha256.New, []byte("WebAppData"))
	h1.Write([]byte(s.telegramToken))
	secretKey := h1.Sum(nil)
	return secretKey
}

func (s *AuthService) isValid(dataStr string, secretKey []byte, hash string) bool {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataStr))
	hh := h.Sum(nil)
	resultHash := hex.EncodeToString(hh)
	return resultHash == hash
}

func (s *AuthService) TelegramLogin(ctx context.Context, request *connect.Request[authpb.TelegramLoginRequest]) (*connect.Response[authpb.LoginResponse], error) {
	var msg = request.Msg
	if !s.isValid(msg.DataStr, s.secretKey1(), msg.Hash) && !s.isValid(msg.DataStr, s.secretKey2(), msg.Hash) {
		return nil, ErrorLoginFailed
	}

	if time.Now().Sub(time.Unix(msg.AuthDate, 0)) > time.Hour {
		return nil, ErrorLoginExpired
	}

	user, err := s.userManager.QueryByTelegram(ctx, msg.UserId)
	if err != nil {
		log.Sugar.Errorf("error while querying user by telegram chat id: %v", err)
		return nil, ErrorLoginFailed
	}

	return s.login(user)
}

type DiscordIdentity struct {
	ID       json.Number `json:"id"`
	Username string      `json:"username"`
}

func (s *AuthService) DiscordLogin(ctx context.Context, request *connect.Request[authpb.DiscordLoginRequest]) (*connect.Response[authpb.LoginResponse], error) {
	token, err := s.discordOAuth2Config.Exchange(context.Background(), request.Msg.GetCode())
	if err != nil {
		log.Sugar.Infof("Error exchanging code for token: %v", err)
		return nil, ErrorLoginFailed
	}

	res, err := s.discordOAuth2Config.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		log.Sugar.Infof("Error getting user (%v): %v", res.StatusCode, err)
		return nil, ErrorLoginFailed
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Sugar.Infof("Error reading response body: %v", err)
		return nil, ErrorInternal
	}

	var identity DiscordIdentity
	err = json.Unmarshal(body, &identity)
	if err != nil {
		log.Sugar.Infof("Error unmarshalling response body: %v", err)
		return nil, ErrorInternal
	}

	id, err := identity.ID.Int64()
	if err != nil {
		log.Sugar.Infof("Error converting id to int64: %v", err)
		return nil, ErrorInternal
	}
	user, err := s.userManager.QueryByDiscordChannel(ctx, id)
	if err != nil {
		return nil, ErrorUserNotFound
	}

	return s.login(user)
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
