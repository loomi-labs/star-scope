/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/loomi-labs/star-scope/chain_crawler"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/grpc"
	"github.com/loomi-labs/star-scope/kafka"
	"log"
	"strings"
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

var startIndexEventConsumerCmd = &cobra.Command{
	Use:   "index-event-consumer",
	Short: "Start the index event consumer",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
		if cmd.Flag("fake").Value.String() == "true" {
			log.Println("Creating fake events")
			fakeWalletAddresses := strings.Split(common.GetEnvX("FAKE_WALLET_ADDRESSES"), ",")
			fakeEventCreator := kafka.NewFakeEventCreator(dbManagers, fakeWalletAddresses, kafkaBrokers...)
			fakeEventCreator.CreateFakeEvents()
		} else {
			eventConsumer := kafka.NewKafka(dbManagers, kafkaBrokers...)
			eventConsumer.ProcessIndexedEvents()
		}
	},
}

var startQueryEventConsumerCmd = &cobra.Command{
	Use:   "query-event-consumer",
	Short: "Start the query event consumer",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
		eventConsumer := kafka.NewKafka(dbManagers, kafkaBrokers...)
		eventConsumer.ProcessQueryEvents()
	},
}

var startCrawlerCmd = &cobra.Command{
	Use:   "crawler",
	Short: "Start the chain crwaler",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		chainCrawler := chain_crawler.NewChainCrawler(dbManagers)
		chainCrawler.AddOrUpdateChains()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(startGrpcServerCmd)
	serviceCmd.AddCommand(startIndexEventConsumerCmd)
	serviceCmd.AddCommand(startQueryEventConsumerCmd)
	serviceCmd.AddCommand(startCrawlerCmd)

	startIndexEventConsumerCmd.Flags().BoolP("fake", "f", false, "Create fake events")
}
