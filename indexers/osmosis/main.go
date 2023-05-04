package main

import (
	"github.com/loomi-labs/star-scope/indexers/osmosis/indexer"
	"github.com/shifty11/go-logger/log"
)

func main() {
	defer log.SyncLogger()
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Panic(err)
			return
		}
	}()

	var indx = indexer.NewIndexer("https://rest.cosmos.directory/osmosis", []string{"localhost:9092"})
	indx.StartIndexing()
}
