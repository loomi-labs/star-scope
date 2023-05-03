/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/shifty11/blocklog-backend/common"
	"github.com/shifty11/blocklog-backend/consumer"
	"github.com/shifty11/blocklog-backend/database"
	"github.com/shifty11/blocklog-backend/grpc"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service commands",
}

// startGrpcServerCmd represents the grpc command
var startGrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		jwtSecretKey := common.GetEnvX("JWT_SECRET_KEY")
		accessTokenDurationStr := common.GetEnvX("ACCESS_TOKEN_DURATION")
		accessTokenDuration, err := time.ParseDuration(accessTokenDurationStr)
		if err != nil {
			log.Panicf("Invalid access token duration: %s", err)
		}
		refreshTokenDurationStr := common.GetEnvX("REFRESH_TOKEN_DURATION")
		refreshTokenDuration, err := time.ParseDuration(refreshTokenDurationStr)
		if err != nil {
			log.Panicf("Invalid refresh token duration: %s", err)
		}
		indexerAuthToken := common.GetEnvX("INDEXER_AUTH_TOKEN")
		var config = &grpc.Config{
			JwtSecretKey:         jwtSecretKey,
			AccessTokenDuration:  accessTokenDuration,
			RefreshTokenDuration: refreshTokenDuration,
			Port:                 50001,
			IndexerAuthToken:     indexerAuthToken,
		}

		dbManagers := database.NewDefaultDbManagers()
		server := grpc.NewGRPCServer(config, dbManagers)
		server.Run()
	},
}

var startEventConsumerCmd = &cobra.Command{
	Use:   "event-consumer",
	Short: "Start the event consumer",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		eventConsumer := consumer.NewKafkaConsumer(dbManagers, "localhost:9092")
		eventConsumer.StartConsuming()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(startGrpcServerCmd)
	serviceCmd.AddCommand(startEventConsumerCmd)
}
