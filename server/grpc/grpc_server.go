package grpc

import (
	"fmt"
	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/grpc/auth"
	"github.com/loomi-labs/star-scope/grpc/auth/authpb/authpbconnect"
	"github.com/loomi-labs/star-scope/grpc/event"
	"github.com/loomi-labs/star-scope/grpc/event/eventpb/eventpbconnect"
	"github.com/loomi-labs/star-scope/grpc/indexer"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb/indexerpbconnect"
	"github.com/loomi-labs/star-scope/grpc/settings"
	"github.com/loomi-labs/star-scope/grpc/settings/settingspb/settingspbconnect"
	"github.com/loomi-labs/star-scope/grpc/user"
	"github.com/loomi-labs/star-scope/grpc/user/userpb/userpbconnect"
	"github.com/loomi-labs/star-scope/kafka"
	"github.com/loomi-labs/star-scope/kafka_internal"
	"github.com/ravener/discord-oauth2"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/net/http2"
	"golang.org/x/oauth2"
	"strings"

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

func (s GRPCServer) Run() {
	log.Sugar.Infof("Starting GRPC server on port %v", s.config.Port)
	jwtManager := auth.NewJWTManager([]byte(s.config.JwtSecretKey), s.config.AccessTokenDuration, s.config.RefreshTokenDuration)
	authInterceptor := auth.NewAuthInterceptor(jwtManager, s.dbManagers.UserManager, auth.AccessibleRoles(), s.config.IndexerAuthToken)
	kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")

	telegramToken := common.GetEnvX("TELEGRAM_BOT_TOKEN")
	discordClientId := common.GetEnvX("DISCORD_CLIENT_ID")
	discordClientSecret := common.GetEnvX("DISCORD_CLIENT_SECRET")
	discordWebAppUrl := common.GetEnvX("DISCORD_WEB_APP_URL")
	var discordConfig = &oauth2.Config{
		RedirectURL:  discordWebAppUrl,
		ClientID:     discordClientId,
		ClientSecret: discordClientSecret,
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}

	interceptors := connect.WithInterceptors(authInterceptor)

	reflector := grpcreflect.NewStaticReflector(auth.ServiceNames()...)
	mux := http.NewServeMux()

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	mux.Handle(authpbconnect.NewAuthServiceHandler(
		auth.NewAuthServiceHandler(s.dbManagers, jwtManager, kafka_internal.NewKafkaInternal(kafkaBrokers), telegramToken, discordConfig),
		interceptors,
	))
	mux.Handle(indexerpbconnect.NewIndexerServiceHandler(
		indexer.NewIndexerServiceHandler(s.dbManagers),
		interceptors,
	))
	mux.Handle(userpbconnect.NewUserServiceHandler(
		user.NewUserServiceHandler(s.dbManagers),
		interceptors,
	))
	mux.Handle(userpbconnect.NewUserSetupServiceHandler(
		user.NewUserSetupServiceHandler(s.dbManagers),
		interceptors,
	))
	mux.Handle(eventpbconnect.NewEventServiceHandler(
		event.NewEventServiceHandler(s.dbManagers, kafka.NewKafka(s.dbManagers, kafkaBrokers)),
		interceptors,
	))
	mux.Handle(settingspbconnect.NewSettingsServiceHandler(
		settings.NewSettingsServiceHandler(s.dbManagers),
		interceptors,
	))

	err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%v", s.config.Port),
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Sugar.Fatal(err)
	}
}
