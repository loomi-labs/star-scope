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
	"sync"
	"time"
)

func startIndexer(grpcClient indexerpbconnect.IndexerServiceClient, updateChannel chan indexer.SyncStatus, stopChannels map[uint64]chan uint64, stopChannelsMutex *sync.Mutex) {
	var kafkaBrokers = strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
	var encodingConfig = indexer.MakeEncodingConfig()
	response, err := grpcClient.GetIndexingChains(context.Background(), connect.NewRequest(&emptypb.Empty{}))
	if err != nil {
		log.Sugar.Errorf("Error getting indexing chains: %v", err)
	} else {
		for _, chain := range response.Msg.GetChains() {
			if _, ok := stopChannels[chain.Id]; ok {
				// Chain is already being indexed
				continue
			}
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

			var indx = indexer.NewIndexer(config)
			var stopChannel = make(chan uint64)
			stopChannelsMutex.Lock()
			stopChannels[chain.Id] = stopChannel
			stopChannelsMutex.Unlock()
			go indx.StartIndexing(chain, updateChannel, stopChannel)
		}
	}
}

func startChainFetchInterval(grpcClient indexerpbconnect.IndexerServiceClient, updateChannel chan indexer.SyncStatus, stopChannels map[uint64]chan uint64, stopChannelsMutex *sync.Mutex) {
	if len(stopChannels) == 0 {
		startIndexer(grpcClient, updateChannel, stopChannels, stopChannelsMutex)
	}
	var timer = time.NewTicker(5 * time.Minute)
	for {
		<-timer.C
		startIndexer(grpcClient, updateChannel, stopChannels, stopChannelsMutex)
	}
}

func startIndexers() {
	var updateChannel = make(chan indexer.SyncStatus)
	defer close(updateChannel)

	var grpcClient = client.NewIndexerServiceClient(
		common.GetEnvX("INDEXER_GRPC_ENDPOINT"),
		common.GetEnvX("INDEXER_AUTH_TOKEN"),
	)

	stopChannels := make(map[uint64]chan uint64)
	stopChannelsMutex := sync.Mutex{}

	go startChainFetchInterval(grpcClient, updateChannel, stopChannels, &stopChannelsMutex)
	listenForUpdates(grpcClient, updateChannel, stopChannels, &stopChannelsMutex)
}

func listenForUpdates(grpcClient indexerpbconnect.IndexerServiceClient, updateChannel chan indexer.SyncStatus, stopChannels map[uint64]chan uint64, stopChannelsMutex *sync.Mutex) {
	const updateBatchTimeout = 30 * time.Second // Time duration to wait for more updates

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
			disabledChainIds := sendUpdates(grpcClient, updates) // Send the batch when the timer expires
			timer.Stop()
			timer = nil // Reset the timer
			updates = make(map[uint64]*indexerpb.IndexingChain)
			for _, chainId := range disabledChainIds {
				stopChannelsMutex.Lock()
				if stopChannel, ok := stopChannels[chainId]; ok {
					stopChannel <- chainId
					delete(stopChannels, chainId)
				}
				stopChannelsMutex.Unlock()
			}
		}
	}
}

func sendUpdates(grpcClient indexerpbconnect.IndexerServiceClient, updates map[uint64]*indexerpb.IndexingChain) []uint64 {
	if len(updates) > 0 {
		log.Sugar.Debugf("Sending %d updates", len(updates))
		var chains []*indexerpb.IndexingChain
		for _, update := range updates {
			chains = append(chains, update)
		}
		var request = connect.NewRequest(&indexerpb.UpdateIndexingChainsRequest{
			Chains: chains,
		})
		response, err := grpcClient.UpdateIndexingChains(context.Background(), request)
		if err != nil {
			log.Sugar.Errorf("Error updating indexing chains: %v", err)
			time.Sleep(1)
			sendUpdates(grpcClient, updates) // Retry sending the batch
		}
		return response.Msg.GetDisabledChainIds()
	}
	return nil
}

func main() {
	defer log.SyncLogger()
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Panic(err)
			return
		}
	}()

	startIndexers()
}
