package grpc

import (
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/shifty11/blocklog-backend/database"
	"github.com/shifty11/blocklog-backend/grpc/auth"
	authconnect "github.com/shifty11/blocklog-backend/grpc/auth/v1/v1connect"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"time"
)

type Config struct {
	Port                 int
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	JwtSecretKey         string
}

//goland:noinspection GoNameStartsWithPackageName
type GRPCServer struct {
	config     *Config
	dbManagers *database.DbManagers
}

func NewGRPCServer(
	config *Config,
	dbManagers *database.DbManagers,
) *GRPCServer {
	return &GRPCServer{
		config:     config,
		dbManagers: dbManagers,
	}
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, connect-protocol-version")

		// If the request method is OPTIONS, just return with no content
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func (s GRPCServer) Run() {
	log.Sugar.Infof("Starting GRPC server on port %v", s.config.Port)
	jwtManager := auth.NewJWTManager([]byte(s.config.JwtSecretKey), s.config.AccessTokenDuration, s.config.RefreshTokenDuration)
	authInterceptor := auth.NewAuthInterceptor(jwtManager, s.dbManagers.UserManager, auth.AccessibleRoles())

	interceptors := connect.WithInterceptors(authInterceptor)

	path, handler := authconnect.NewAuthServiceHandler(
		auth.NewAuthServiceHandler(s.dbManagers.UserManager, jwtManager),
		interceptors,
	)

	mux := http.NewServeMux()
	mux.Handle(path, handler)
	handler = corsHandler(mux) // Wrap the mux router with the CORS handler
	err := http.ListenAndServe(
		fmt.Sprintf("localhost:%v", s.config.Port),
		h2c.NewHandler(handler, &http2.Server{}),
	)
	if err != nil {
		log.Sugar.Fatal(err)
	}
}
