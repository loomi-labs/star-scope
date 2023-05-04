package grpc

import (
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/grpc/auth"
	"github.com/loomi-labs/star-scope/grpc/auth/authpb/authpbconnect"
	"github.com/loomi-labs/star-scope/grpc/event"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb/eventpbconnect"
	"github.com/loomi-labs/star-scope/grpc/indexer"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb/indexerpbconnect"
	"github.com/loomi-labs/star-scope/grpc/project"
	"github.com/loomi-labs/star-scope/grpc/project/projectpb/projectpbconnect"
	"github.com/loomi-labs/star-scope/kafka"
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
	IndexerAuthToken     string
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Connect-Protocol-Version, X-grpc-web, X-user-agent")

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
	authInterceptor := auth.NewAuthInterceptor(jwtManager, s.dbManagers.UserManager, auth.AccessibleRoles(), s.config.IndexerAuthToken)

	interceptors := connect.WithInterceptors(authInterceptor)

	mux := http.NewServeMux()
	mux.Handle(authpbconnect.NewAuthServiceHandler(
		auth.NewAuthServiceHandler(s.dbManagers, jwtManager),
		interceptors,
	))
	mux.Handle(indexerpbconnect.NewIndexerServiceHandler(
		indexer.NewIndexerServiceHandler(s.dbManagers),
		interceptors,
	))
	mux.Handle(projectpbconnect.NewProjectServiceHandler(
		project.NewProjectServiceHandler(),
		interceptors,
	))
	mux.Handle(eventpbconnect.NewEventServiceHandler(
		event.NewEventServiceHandler(kafka.NewKafka(s.dbManagers, common.GetEnvX("KAFKA_BROKERS"))),
		interceptors,
	))

	handler := corsHandler(mux) // Wrap the mux router with the CORS handler
	err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%v", s.config.Port),
		h2c.NewHandler(handler, &http2.Server{}),
	)
	if err != nil {
		log.Sugar.Fatal(err)
	}
}
