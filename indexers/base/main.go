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
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	GRPCClient        indexerpbconnect.IndexerServiceClient
	UpdateChannel     chan indexer.SyncStatus
	StopChannels      map[uint64]chan struct{}
	StopChannelsMutex *sync.Mutex
	KakfaBrokers      []string
	EncodingConfig    indexer.EncodingConfig
}

func newConfig() *Config {
	kafkaBrokers := strings.Split(common.GetEnvX("KAFKA_BROKERS"), ",")
	encodingConfig := indexer.MakeEncodingConfig()
	grpcClient := client.NewIndexerServiceClient(
		common.GetEnvX("INDEXER_GRPC_ENDPOINT"),
		common.GetEnvX("INDEXER_AUTH_TOKEN"),
	)
	return &Config{
		GRPCClient:        grpcClient,
		UpdateChannel:     make(chan indexer.SyncStatus),
		StopChannels:      make(map[uint64]chan struct{}),
		StopChannelsMutex: &sync.Mutex{},
		KakfaBrokers:      kafkaBrokers,
		EncodingConfig:    encodingConfig,
	}
}

func (c *Config) Close() {
	close(c.UpdateChannel)
}

func startIndexer(chain *indexerpb.IndexingChain, config *Config, g *errgroup.Group) {
	config.StopChannelsMutex.Lock()
	stopChannel, ok := config.StopChannels[chain.Id]
	config.StopChannelsMutex.Unlock()
	if ok {
		// Chain is already being indexed
		return
	}

	stopChannel = make(chan struct{})
	config.StopChannelsMutex.Lock()
	config.StopChannels[chain.Id] = stopChannel
	config.StopChannelsMutex.Unlock()

	var messageHandler indexer.TxHandler
	if chain.HasCustomIndexer {
		messageHandler = indexer.NewCustomMessageHandler(chain, config.EncodingConfig, "http://localhost:50002")
	} else {
		messageHandler = indexer.NewBaseMessageHandler(chain, config.EncodingConfig)
	}

	indx := indexer.NewIndexer(chain, config.EncodingConfig, config.KakfaBrokers, messageHandler)
	g.Go(func() error {
		indx.StartIndexing(config.UpdateChannel, stopChannel)
		return nil
	})
}

func startChainFetchInterval(config *Config, g *errgroup.Group) {
	g.Go(func() error {
		startIndexerRequest := func() error {
			response, err := config.GRPCClient.GetIndexingChains(context.Background(), connect.NewRequest(&emptypb.Empty{}))
			if err != nil {
				return err
			}

			for _, chain := range response.Msg.GetChains() {
				startIndexer(chain, config, g)
			}

			return nil
		}

		// Start indexing immediately
		if err := startIndexerRequest(); err != nil {
			return err
		}

		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-ticker.C:
				if err := startIndexerRequest(); err != nil {
					log.Sugar.Errorf("Error starting indexer: %v", err)
				}
			}
		}
	})
}

func listenForUpdates(config *Config) {
	const updateBatchTimeout = 30 * time.Second // Time duration to wait for more updates

	updates := make(map[uint64]*indexerpb.IndexingChain)
	timer := time.NewTimer(updateBatchTimeout)
	timerExpired := timer.C

	for {
		select {
		case update, ok := <-config.UpdateChannel:
			if !ok {
				// Channel is closed, call sendUpdates and exit the function
				log.Sugar.Info("Update channel closed, exiting")
				sendUpdates(config, updates)
				return
			}

			handledMessageTypes := make([]string, 0, len(update.HandledMessageTypes))
			for msgType := range update.HandledMessageTypes {
				handledMessageTypes = append(handledMessageTypes, msgType)
			}
			unhandledMessageTypes := make([]string, 0, len(update.UnhandledMessageTypes))
			for msgType := range update.UnhandledMessageTypes {
				unhandledMessageTypes = append(unhandledMessageTypes, msgType)
			}
			updates[update.ChainId] = &indexerpb.IndexingChain{
				Id:                    update.ChainId,
				IndexingHeight:        update.Height,
				HandledMessageTypes:   handledMessageTypes,
				UnhandledMessageTypes: unhandledMessageTypes,
			}
		case <-timerExpired:
			disabledChainIds := sendUpdates(config, updates) // Send the batch when the timer expires
			timer.Stop()
			timer = time.NewTimer(updateBatchTimeout)
			timerExpired = timer.C
			updates = make(map[uint64]*indexerpb.IndexingChain)
			for _, chainID := range disabledChainIds {
				if stopChannel, ok := config.StopChannels[chainID]; ok {
					config.StopChannelsMutex.Lock()
					close(stopChannel)
					delete(config.StopChannels, chainID)
					config.StopChannelsMutex.Unlock()
				}
			}
		}
	}
}

func sendUpdates(config *Config, updates map[uint64]*indexerpb.IndexingChain) []uint64 {
	if len(updates) > 0 {
		log.Sugar.Debugf("Sending %d updates", len(updates))
		chains := make([]*indexerpb.IndexingChain, 0, len(updates))
		for _, update := range updates {
			chains = append(chains, update)
		}
		request := connect.NewRequest(&indexerpb.UpdateIndexingChainsRequest{
			Chains: chains,
		})
		response, err := config.GRPCClient.UpdateIndexingChains(context.Background(), request)
		if err != nil {
			log.Sugar.Errorf("Error updating indexing chains: %v", err)
			return nil
		}
		return response.Msg.GetDisabledChainIds()
	}
	return nil
}

func main() {
	defer log.SyncLogger()

	config := newConfig()
	defer config.Close()

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		startChainFetchInterval(config, g)
		return nil
	})
	g.Go(func() error {
		listenForUpdates(config)
		return nil
	})

	// Wait for interrupt signal to gracefully shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)
	select {
	case <-interrupt:
		// Received an interrupt signal, initiate graceful shutdown
		log.Sugar.Info("Received interrupt signal, shutting down")
		cancel()
	case <-ctx.Done():
		// Context canceled, exiting the program
	}

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		log.Sugar.Fatalf("Error during goroutine execution: %v", err)
	}
}
