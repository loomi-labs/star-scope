package main

import (
	"github.com/loomi-labs/star-scope/indexers/osmosis/common"
	"github.com/loomi-labs/star-scope/indexers/osmosis/indexer"
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

	var indx = indexer.NewIndexer(
		common.GetEnvX("REST_ENDPOINT"),
		strings.Split(common.GetEnvX("KAFKA_BROKERS"), ","),
	)
	indx.StartIndexing()
}
