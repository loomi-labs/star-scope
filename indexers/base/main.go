package main

import (
	"github.com/loomi-labs/star-scope/indexers/base/common"
	"github.com/loomi-labs/star-scope/indexers/base/indexer"
	"github.com/shifty11/go-logger/log"
	"strings"
)

func main() {
	defer log.SyncLogger()
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Panic(err)
			return
		}
	}()

	var encodingConfig = indexer.MakeEncodingConfig()
	var chainPath = common.GetEnvX("CHAIN_PATH")
	var restEndpoint = common.GetEnvX("REST_ENDPOINT") + "/" + chainPath
	var config = indexer.IndexerConfig{
		ChainInfo:      indexer.ChainInfo{Path: chainPath, RestEndpoint: restEndpoint},
		KafkaAddresses: strings.Split(common.GetEnvX("KAFKA_BROKERS"), ","),
		EncodingConfig: indexer.EncodingConfig{
			InterfaceRegistry: encodingConfig.InterfaceRegistry,
			Codec:             encodingConfig.Codec,
			TxConfig:          encodingConfig.TxConfig,
		},
		GrpcAuthToken: common.GetEnvX("INDEXER_AUTH_TOKEN"),
		GrpcEndpoint:  common.GetEnvX("INDEXER_GRPC_ENDPOINT"),
	}
	var indx = indexer.NewIndexer(config, indexer.NewBaseMessageHandler(config.ChainInfo, config.EncodingConfig))
	indx.StartIndexing()
}
