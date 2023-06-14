/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/crawler/chain_crawler"
	"github.com/loomi-labs/star-scope/crawler/governance_crawler"
	"github.com/loomi-labs/star-scope/crawler/validator_crawler"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/grpc"
	"github.com/loomi-labs/star-scope/kafka"
	"github.com/loomi-labs/star-scope/kafka_internal"
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

var startWalletEventConsumerCmd = &cobra.Command{
	Use:   "wallet-event-consumer",
	Short: "Start the wallet event consumer",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
		if cmd.Flag("fake").Value.String() == "true" {
			log.Println("Creating fake events")
			fakeWalletAddresses := strings.Split(common.GetEnvX("FAKE_WALLET_ADDRESSES"), ",")
			fakeEventCreator := kafka.NewFakeEventCreator(dbManagers, fakeWalletAddresses, kafkaBrokers...)
			fakeEventCreator.CreateFakeEvents()
		} else {
			eventConsumer := kafka.NewKafka(dbManagers, kafkaBrokers)
			eventConsumer.ProcessWalletEvents()
		}
	},
}

var startChainEventConsumerCmd = &cobra.Command{
	Use:   "chain-event-consumer",
	Short: "Start the chain event consumer",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
		eventConsumer := kafka.NewKafka(dbManagers, kafkaBrokers)
		eventConsumer.ProcessChainEvents()
	},
}

var startContractEventConsumerCmd = &cobra.Command{
	Use:   "contract-event-consumer",
	Short: "Start the contract event consumer",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
		eventConsumer := kafka.NewKafka(dbManagers, kafkaBrokers)
		eventConsumer.ProcessContractEvents()
	},
}

var startChainCrawlerCmd = &cobra.Command{
	Use:   "chain-crawler",
	Short: "Start the chain crawler",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		chainCrawler := chain_crawler.NewChainCrawler(dbManagers)
		chainCrawler.StartCrawling()
	},
}

var startGovernanceCrawlerCmd = &cobra.Command{
	Use:   "gov-crawler",
	Short: "Start the governance crawler",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
		kafkaInternal := kafka_internal.NewKafkaInternal(kafkaBrokers)
		chainCrawler := governance_crawler.NewGovernanceCrawler(dbManagers, kafkaInternal)
		chainCrawler.StartCrawling()
	},
}

var startValidatorCrawlerCmd = &cobra.Command{
	Use:   "validator-crawler",
	Short: "Start the validator crawler",
	Run: func(cmd *cobra.Command, args []string) {
		dbManagers := database.NewDefaultDbManagers()
		kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
		kafkaInternal := kafka_internal.NewKafkaInternal(kafkaBrokers)
		crawler := validator_crawler.NewValidatorCrawler(dbManagers, kafkaInternal)
		crawler.StartCrawling()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(startGrpcServerCmd)
	serviceCmd.AddCommand(startWalletEventConsumerCmd)
	serviceCmd.AddCommand(startChainEventConsumerCmd)
	serviceCmd.AddCommand(startContractEventConsumerCmd)

	serviceCmd.AddCommand(startChainCrawlerCmd)
	serviceCmd.AddCommand(startGovernanceCrawlerCmd)
	serviceCmd.AddCommand(startValidatorCrawlerCmd)

	startWalletEventConsumerCmd.Flags().BoolP("fake", "f", false, "Create fake events")
}
