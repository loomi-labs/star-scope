package main

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/indexers/base/client"
	"github.com/loomi-labs/star-scope/indexers/base/common"
	"github.com/loomi-labs/star-scope/indexers/base/indexer"
	"github.com/loomi-labs/star-scope/indexers/base/neutron"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"time"
)

func startIndexers(updateChannel chan indexer.SyncStatus) indexerpbconnect.IndexerServiceClient {
	var kafkaBrokers = strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")

	var encodingConfig = indexer.MakeEncodingConfig()
	var grpcClient = client.NewIndexerServiceClient(
		common.GetEnvX("INDEXER_GRPC_ENDPOINT"),
		common.GetEnvX("INDEXER_AUTH_TOKEN"),
	)
	response, err := grpcClient.GetIndexingChains(context.Background(), connect.NewRequest(&emptypb.Empty{}))
	if err != nil {
		log.Sugar.Panicf("Error getting indexing chains: %v", err)
	}

	governanceCrawler := indexer.NewGovernanceCrawler(grpcClient, kafkaBrokers)
	go governanceCrawler.StartCrawling()

	setupCrawler := indexer.NewSetupCrawler(grpcClient, kafkaBrokers)
	go setupCrawler.StartCrawling()

	for _, chain := range response.Msg.GetChains() {
		var config = indexer.Config{
			ChainInfo: indexer.ChainInfo{
				ChainId:      chain.Id,
				Path:         chain.Path,
				RestEndpoint: chain.RestEndpoint,
				Name:         chain.Name,
			},
			KafkaBrokers:   kafkaBrokers,
			EncodingConfig: encodingConfig,
		}
		if chain.HasCustomIndexer {
			config.MessageHandler = indexer.NewCustomMessageHandler(config.ChainInfo, config.EncodingConfig, "http://localhost:50002")
		} else {
			config.MessageHandler = indexer.NewBaseMessageHandler(config.ChainInfo, config.EncodingConfig)
		}

		if chain.Path == "neutron" {
			go neutron.NewNeutronCrawler(
				grpcClient,
				chain.Path,
				"neutron1suhgf5svhu4usrurvxzlgn54ksxmn8gljarjtxqnapv8kjnp4nrstdxvff",
				"neutron1h6828as2z5av0xqtlh4w9m75wxewapk8z9l2flvzc29zeyzhx6fqgp648z",
				kafkaBrokers,
			).StartCrawling()
		} else if chain.Path == "neutron-testnet" {
			go neutron.NewNeutronCrawler(
				grpcClient,
				chain.Path,
				"neutron1suhgf5svhu4usrurvxzlgn54ksxmn8gljarjtxqnapv8kjnp4nrstdxvff",
				"",
				kafkaBrokers,
			).StartCrawling()
		}

		var indx = indexer.NewIndexer(config)
		go indx.StartIndexing(chain, updateChannel)
	}
	return grpcClient
}

func listenForUpdates(grpcClient indexerpbconnect.IndexerServiceClient, updateChannel chan indexer.SyncStatus) {
	const updateBatchTimeout = 5 * time.Second // Time duration to wait for more updates

	var updates = make(map[uint64]*indexerpb.IndexingChain)
	var timer *time.Timer
	var timerExpired <-chan time.Time

	for {
		select {
		case update, ok := <-updateChannel:
			if !ok {
				// Channel is closed, call sendUpdates and exit the function
				log.Sugar.Info("Update channel closed, exiting")
				sendUpdates(grpcClient, updates)
				return
			}

			var handledMessageTypes []string
			for msgType := range update.HandledMessageTypes {
				handledMessageTypes = append(handledMessageTypes, msgType)
			}
			var unhandledMessageTypes []string
			for msgType := range update.UnhandledMessageTypes {
				unhandledMessageTypes = append(unhandledMessageTypes, msgType)
			}
			updates[update.ChainId] = &indexerpb.IndexingChain{
				Id:                    update.ChainId,
				IndexingHeight:        update.Height,
				HandledMessageTypes:   handledMessageTypes,
				UnhandledMessageTypes: unhandledMessageTypes,
			}

			// Start the timer if it's not running
			if timer == nil {
				timer = time.NewTimer(updateBatchTimeout)
				timerExpired = timer.C
			}
		case <-timerExpired:
			sendUpdates(grpcClient, updates) // Send the batch when the timer expires
			timer.Stop()
			timer = nil // Reset the timer
			updates = make(map[uint64]*indexerpb.IndexingChain)
		}
	}
}

func sendUpdates(grpcClient indexerpbconnect.IndexerServiceClient, updates map[uint64]*indexerpb.IndexingChain) {
	if len(updates) > 0 {
		log.Sugar.Debugf("Sending %d updates", len(updates))
		var chains []*indexerpb.IndexingChain
		for _, update := range updates {
			chains = append(chains, update)
		}
		var request = connect.NewRequest(&indexerpb.UpdateIndexingChainsRequest{
			Chains: chains,
		})
		_, err := grpcClient.UpdateIndexingChains(context.Background(), request)
		if err != nil {
			log.Sugar.Errorf("Error updating indexing chains: %v", err)
			time.Sleep(1)
			sendUpdates(grpcClient, updates) // Retry sending the batch
		}
	}
}

func main() {
	defer log.SyncLogger()
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Panic(err)
			return
		}
	}()

	var updateChannel = make(chan indexer.SyncStatus)
	defer close(updateChannel)

	grpcClient := startIndexers(updateChannel)
	listenForUpdates(grpcClient, updateChannel)
}
