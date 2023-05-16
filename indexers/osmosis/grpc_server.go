package main

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"fmt"
	"github.com/shifty11/go-logger/log"
	"golang.org/x/net/http2"

	"golang.org/x/net/http2/h2c"
	"net/http"
)

type Config struct {
	Port int
}

//goland:noinspection GoNameStartsWithPackageName
type GRPCServer struct {
	config *Config
}

func NewGRPCServer(
	config *Config,
) *GRPCServer {
	return &GRPCServer{
		config: config,
	}
}

func (s GRPCServer) Run() {
	log.Sugar.Infof("Starting GRPC server on port %v", s.config.Port)
	mux := http.NewServeMux()

	mux.Handle(indexerpbconnect.NewTxHandlerServiceHandler(
		NewTxHandlerServiceHandler(),
	))

	err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%v", s.config.Port),
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Sugar.Fatal(err)
	}
}
