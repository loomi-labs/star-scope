package chain_crawler

import (
	"context"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/types"
	"github.com/robfig/cron/v3"
	"github.com/shifty11/go-logger/log"
)

type ChainCrawler struct {
	chainManager *database.ChainManager
}

func NewChainCrawler(dbManagers *database.DbManagers) *ChainCrawler {
	return &ChainCrawler{
		chainManager: dbManagers.ChainManager,
	}
}

func (c *ChainCrawler) chainNeedsUpdate(entChain *ent.Chain, chainInfo *types.ChainData) bool {
	return entChain.ChainID != chainInfo.ChainId ||
		entChain.Name != chainInfo.Name ||
		entChain.PrettyName != chainInfo.PrettyName ||
		entChain.Path != chainInfo.Path ||
		entChain.Bech32Prefix != chainInfo.Bech32Prefix ||
		entChain.Image != chainInfo.Image
}

func (c *ChainCrawler) isChainValid(chainInfo *types.ChainData) bool {
	return chainInfo.NetworkType == "mainnet" && chainInfo.ChainId != "" && chainInfo.Path != "" && chainInfo.Bech32Prefix != ""
}

func (c *ChainCrawler) AddOrUpdateChains() {
	log.Sugar.Debug("Updating chains")
	var chainInfo types.ChainInfo
	_, err := common.GetJson("https://chains.cosmos.directory/", 5, &chainInfo)
	if err != nil {
		log.Sugar.Errorf("while getting chainData info: %v", err)
	}

	var chains = c.chainManager.QueryAll(context.Background())

	for _, chainData := range chainInfo.Chains {
		if !c.isChainValid(&chainData) {
			continue
		}
		var found = false
		for _, entChain := range chains {
			if entChain.Name == chainData.Name {
				found = true
				if c.chainNeedsUpdate(entChain, &chainData) {
					_, err := c.chainManager.UpdateChainInfo(context.Background(), entChain, &chainData)
					if err != nil {
						log.Sugar.Errorf("while updating chain: %v", err)
					}
				}
				break
			}
		}
		if !found && chainData.NetworkType == "mainnet" {
			_, err := c.chainManager.Create(context.Background(), &chainData)
			if err != nil {
				log.Sugar.Errorf("while creating chain: %v", err)
			}
		}
	}
}

func (c *ChainCrawler) ScheduleCrawl() {
	log.Sugar.Info("Scheduling chain crawl")
	cr := cron.New()
	_, err := cr.AddFunc("0 9 * * *", func() { c.AddOrUpdateChains() }) // every day at 9:00
	if err != nil {
		log.Sugar.Errorf("while executing 'AddOrUpdateChains' via cron: %v", err)
	}
	cr.Start()
}
