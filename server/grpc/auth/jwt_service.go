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

func AccessibleRoles() map[string][]Role {
	const path = "/starscope.grpc"
	const authService = path + ".AuthService/"
	const indexerService = path + ".IndexerService/"
	const userService = path + ".UserService/"
	const eventService = path + ".EventService/"

	roles := map[string][]Role{
		authService + "KeplrLogin":         {Unauthenticated, User, Admin},
		authService + "RefreshAccessToken": {Unauthenticated, User, Admin},
		userService + "GetUser":            {User, Admin},
		userService + "ListChannels":       {User, Admin},
		eventService + "EventStream":       {User, Admin},
		indexerService + "GetHeight":       {Token},
		indexerService + "UpdateHeight":    {Token},
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