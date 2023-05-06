/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/grpc"
	"github.com/loomi-labs/star-scope/kafka"
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
		eventConsumer := kafka.NewKafka(dbManagers, common.GetEnvX("KAFKA_BROKERS"))
		eventConsumer.ProcessIndexedEvents()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(startGrpcServerCmd)
	serviceCmd.AddCommand(startEventConsumerCmd)
}
