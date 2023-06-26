package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/user"
	"time"
)

type Role string

var (
	Unauthenticated Role = "unauthenticated"
	User                 = Role(user.RoleUser.String())
	Admin                = Role(user.RoleAdmin.String())
	Token           Role = "token"
)

type TokenType string

const (
	AccessToken  TokenType = "AccessToken"
	RefreshToken TokenType = "RefreshToken"
)

const (
	basePath         = "/starscope.grpc"
	authService      = basePath + ".AuthService/"
	userService      = basePath + ".UserService/"
	userSetupService = basePath + ".UserSetupService/"
	eventService     = basePath + ".EventService/"
	indexerService   = basePath + ".IndexerService/"
)

func ServiceNames() []string {
	return []string{
		authService,
		userService,
		userSetupService,
		eventService,
		indexerService,
	}
}

func AccessibleRoles() map[string][]Role {
	roles := map[string][]Role{
		authService + "KeplrLogin":              {Unauthenticated, User, Admin},
		authService + "RefreshAccessToken":      {Unauthenticated, User, Admin},
		authService + "DiscordLogin":            {Unauthenticated, User, Admin},
		authService + "TelegramLogin":           {Unauthenticated, User, Admin},
		authService + "ConnectDiscord":          {User, Admin},
		authService + "ConnectTelegram":         {User, Admin},
		userService + "GetUser":                 {User, Admin},
		userService + "ListChannels":            {User, Admin},
		userService + "DeleteAccount":           {User, Admin},
		userService + "ListDiscordChannels":     {User, Admin},
		userService + "DeleteDiscordChannel":    {User, Admin},
		userService + "ListTelegramChats":       {User, Admin},
		userService + "DeleteTelegramChat":      {User, Admin},
		userSetupService + "GetCurrentStep":     {User, Admin},
		userSetupService + "FinishStep":         {User, Admin},
		eventService + "EventStream":            {User, Admin},
		eventService + "ListEvents":             {User, Admin},
		eventService + "ListChains":             {User, Admin},
		eventService + "ListEventsCount":        {User, Admin},
		eventService + "MarkEventRead":          {User, Admin},
		eventService + "GetWelcomeMessage":      {User, Admin},
		indexerService + "GetIndexingChains":    {Token},
		indexerService + "UpdateIndexingChains": {Token},
	}
	return roles
}

type Claims struct {
	jwt.RegisteredClaims
	UserId int  `json:"user_id"`
	Role   Role `json:"role,omitempty"`
}

type JWTManager struct {
	jwtSecretKey         []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(jwtSecretKey []byte, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		jwtSecretKey:         jwtSecretKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (manager *JWTManager) GenerateToken(user *ent.User, tokenType TokenType) (string, error) {
	expirationTime := time.Now().Add(manager.accessTokenDuration)
	if tokenType == RefreshToken {
		expirationTime = time.Now().Add(manager.refreshTokenDuration)
	}

	claims := &Claims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Role: Role(user.Role.String()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(manager.jwtSecretKey)
}

func (manager *JWTManager) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return manager.jwtSecretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
