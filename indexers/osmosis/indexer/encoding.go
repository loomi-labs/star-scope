package indexer

import (
	"github.com/osmosis-labs/osmosis/osmoutils/noapptest"
	"github.com/osmosis-labs/osmosis/v15/app/keepers"
)

func GetEncodingConfig() noapptest.TestEncodingConfig {
	return noapptest.MakeTestEncodingConfig(keepers.AppModuleBasics...)
}
