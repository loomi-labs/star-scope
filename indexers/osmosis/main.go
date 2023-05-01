package main

import (
	"github.com/shifty11/blocklog-backend/indexers/osmosis/indexer"
	"github.com/shifty11/go-logger/log"
)

func main() {
	defer log.SyncLogger()
	indexer.StartIndexing("https://rest.cosmos.directory/osmosis")
}
