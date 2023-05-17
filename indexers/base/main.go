package main

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/indexers/base/client"
	"github.com/loomi-labs/star-scope/indexers/base/common"
	"github.com/loomi-labs/star-scope/indexers/base/indexer"
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

	for _, chain := range response.Msg.GetChains() {
		var config = indexer.Config{
			ChainInfo: indexer.ChainInfo{
				Path:         chain.Path,
				RestEndpoint: chain.RpcUrl,
				Name:         chain.Name,
			},
			KafkaBrokers:   kafkaBrokers,
			EncodingConfig: encodingConfig,
		}
		if chain.HasCustomIndexer {
			config.MessageHandler = indexer.NewCustomMessageHandler("http://localhost:50002")
		} else {
			config.MessageHandler = indexer.NewBaseMessageHandler(config.ChainInfo, config.EncodingConfig)
		}
		var unhandledMessageTypes = make(map[string]struct{})
		for _, msgType := range chain.UnhandledMessageTypes {
			unhandledMessageTypes[msgType] = struct{}{}
		}
		var syncStatus = indexer.SyncStatus{
			ChainId:               chain.Id,
			Height:                chain.IndexingHeight,
			LatestHeight:          0,
			UnhandledMessageTypes: unhandledMessageTypes,
		}
		var indx = indexer.NewIndexer(config)
		go indx.StartIndexing(syncStatus, updateChannel)
	}
	return grpcClient
}

func listenForUpdates(grpcClient indexerpbconnect.IndexerServiceClient, updateChannel chan indexer.SyncStatus) {
	const updateBatchTimeout = 5 * time.Second // Time duration to wait for more updates

	var updates = make(map[uint64]*indexerpb.Chain)
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

			var unhandledMessageTypes []string
			for msgType := range update.UnhandledMessageTypes {
				unhandledMessageTypes = append(unhandledMessageTypes, msgType)
			}
			updates[update.ChainId] = &indexerpb.Chain{
				Id:                    update.ChainId,
				IndexingHeight:        update.Height,
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
			updates = make(map[uint64]*indexerpb.Chain)
		}
	}
}

func sendUpdates(grpcClient indexerpbconnect.IndexerServiceClient, updates map[uint64]*indexerpb.Chain) {
	if len(updates) > 0 {
		log.Sugar.Debugf("Sending %d updates", len(updates))
		var chains []*indexerpb.Chain
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
