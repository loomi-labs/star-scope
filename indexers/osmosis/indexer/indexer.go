package indexer

import (
	"buf.build/gen/go/rapha/blocklog/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/rapha/blocklog/protocolbuffers/go/grpc/indexer/indexerpb"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/bufbuild/connect-go"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcChannel "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/golang/protobuf/proto"
	"github.com/osmosis-labs/osmosis/osmoutils/noapptest"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	"github.com/shifty11/blocklog-backend/indexers/osmosis/client"
	indexEvent "github.com/shifty11/blocklog-backend/indexers/osmosis/index_event"
	"github.com/shifty11/go-logger/log"
	"github.com/tendermint/tendermint/abci/types"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type ChainInfo struct {
	ChainName string
}
type Indexer struct {
	baseUrl    string
	chainInfo  ChainInfo
	grpcClient indexerpbconnect.IndexerServiceClient
	//rmqClient      rmq.Connection
	encodingConfig noapptest.TestEncodingConfig
}

func NewIndexer(baseUrl string) Indexer {
	return Indexer{
		baseUrl: baseUrl,
		chainInfo: ChainInfo{
			ChainName: "osmosis",
		},
		grpcClient: client.GetClient(),
		//rmqClient:      openRmqConnection(),
		encodingConfig: GetEncodingConfig(),
	}
}

func (i *Indexer) handleBlock(blockResponse *cmtservice.GetBlockByHeightResponse) {
	log.Sugar.Debugf("handleBlock: %v", blockResponse.GetBlock().GetHeader().Height)
	var data = blockResponse.GetBlock().GetData()
	var txs = data.GetTxs()
	for _, tx := range txs {
		txDecoded, err := i.encodingConfig.TxConfig.TxDecoder()(tx)
		if err != nil {
			//log.Sugar.Error(err)
			log.Sugar.Info("Failed to decode txDecoded")
			continue
		}
		for _, anyMsg := range txDecoded.GetMsgs() {
			switch msg := anyMsg.(type) {
			case *banktypes.MsgSend:
				i.handleMsgSend(msg, tx)
			case *banktypes.MsgMultiSend:
				i.handleMsgMultiSend(msg, tx)
			case *ibcChannel.MsgRecvPacket:
				i.handleMsgRecvPacket(msg, tx)
			case *lockuptypes.MsgBeginUnlockingAll:
				i.handleMsgBeginUnlockingAll(msg)
			case *lockuptypes.MsgBeginUnlocking:
				i.handleMsgBeginUnlocking(msg, tx)
			default:
				log.Sugar.Debugf("Unknown message type")
			}
		}
	}
}

func (i *Indexer) getTxResult(tx []byte) txtypes.GetTxResponse {
	hash := sha256.Sum256(tx)
	hashString := hex.EncodeToString(hash[:])

	var url = fmt.Sprintf("https://rest.cosmos.directory/osmosis/cosmos/tx/v1beta1/txs/%v", hashString)
	resp, err := http.Get(url)
	if err != nil {
		log.Sugar.Panic(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Sugar.Panic(err)
	}
	encodingConfig := GetEncodingConfig()
	var txResponse txtypes.GetTxResponse
	if err := encodingConfig.Codec.UnmarshalJSON(body, &txResponse); err != nil {
		log.Sugar.Panic(err)
	}
	return txResponse
}

func (i *Indexer) getTxEvents(tx []byte) []types.Event {
	var resp = i.getTxResult(tx)
	if resp.GetTxResponse().Code == 0 {
		return resp.GetTxResponse().Events
	}
	return nil
}

func (i *Indexer) handleCoinReceivedEvent(events []types.Event) {
	txEvent := &indexEvent.TxEvent{
		ChainName: i.chainInfo.ChainName,
		Event: &indexEvent.TxEvent_CoinReceived{
			CoinReceived: &indexEvent.CoinReceivedEvent{},
		},
	}
	for _, event := range events {
		switch event.Type {
		case "coin_spent":
			for _, attribute := range event.Attributes {
				switch string(attribute.GetKey()) {
				case "spender":
					txEvent.GetCoinReceived().Sender = string(attribute.GetValue())
				}
			}
		case "coin_received":
			for _, attribute := range event.Attributes {
				switch string(attribute.GetKey()) {
				case "receiver":
					txEvent.WalletAddress = string(attribute.GetValue())
				case "amount":
					pattern := `^(\d+)(ibc\/[a-fA-F0-9]+|[a-zA-Z]+)$`
					re := regexp.MustCompile(pattern)
					matches := re.FindStringSubmatch(string(attribute.GetValue()))
					if len(matches) == 3 {
						amount, err := strconv.Atoi(matches[1])
						if err != nil {
							log.Sugar.Errorf("Failed to parse amount: %v", string(attribute.GetValue()))
							break
						}
						txEvent.GetCoinReceived().Amount = uint64(amount)
						txEvent.GetCoinReceived().Coin = matches[2]
					} else {
						log.Sugar.Errorf("Failed to parse amount: %v", string(attribute.GetValue()))
						break
					}
				}
			}
		}
	}
	if txEvent.GetWalletAddress() == "" ||
		txEvent.GetCoinReceived().GetAmount() == 0 ||
		txEvent.GetCoinReceived().GetCoin() == "" ||
		txEvent.GetCoinReceived().GetSender() == "" {
		log.Sugar.Errorf("Failed to parse coin_received event: %v", txEvent.String())
	} else {
		i.encodeAndPublish(txEvent)
	}
}

func (i *Indexer) handleMsgSend(msg *banktypes.MsgSend, tx []byte) {
	log.Sugar.Infof("MsgSend: %v", msg.String())
	i.handleCoinReceivedEvent(i.getTxEvents(tx))
}

func (i *Indexer) handleMsgMultiSend(msg *banktypes.MsgMultiSend, tx []byte) {
	log.Sugar.Infof("MsgMultiSend: %v", msg.String())
}

func (i *Indexer) handleMsgRecvPacket(msg *ibcChannel.MsgRecvPacket, tx []byte) {
	log.Sugar.Infof("MsgRecvPacket: %v", msg.String())
	i.handleCoinReceivedEvent(i.getTxEvents(tx))
}

func (i *Indexer) handleMsgBeginUnlockingAll(msg *lockuptypes.MsgBeginUnlockingAll) {
	log.Sugar.Infof("MsgBeginUnlockingAll: %v", msg.String())
}

func (i *Indexer) handleMsgBeginUnlocking(msg *lockuptypes.MsgBeginUnlocking, tx []byte) {
	log.Sugar.Infof("MsgBeginUnlocking: %v", msg.String())
	txEvent := &indexEvent.TxEvent{
		ChainName: i.chainInfo.ChainName,
		Event: &indexEvent.TxEvent_OsmosisPoolUnlock{
			OsmosisPoolUnlock: &indexEvent.OsmosisPoolUnlockEvent{},
		},
	}
	for _, event := range i.getTxEvents(tx) {
		log.Sugar.Infof("Event: %v", event.Type)
		for _, attribute := range event.Attributes {
			log.Sugar.Infof("Key: %v, Value: %v", string(attribute.GetKey()), string(attribute.GetValue()))
		}
		switch event.Type {
		case "begin_unlock":
			for _, attribute := range event.Attributes {
				switch string(attribute.GetKey()) {
				case "owner":
					txEvent.WalletAddress = string(attribute.GetValue())
				case "duration":
					duration, err := parseDuration(string(attribute.GetValue()))
					if err != nil {
						log.Sugar.Errorf("Failed to parse duration: %v", err)
						break
					}
					txEvent.GetOsmosisPoolUnlock().Duration = duration
				case "unlock_time":
					ts, err := parseTime(string(attribute.GetValue()))
					if err != nil {
						log.Sugar.Errorf("Failed to parse time: %v", err)
						break
					}
					txEvent.GetOsmosisPoolUnlock().UnlockTime = ts
				}
			}

		}
	}
	if txEvent.GetOsmosisPoolUnlock().GetDuration().IsValid() &&
		txEvent.GetOsmosisPoolUnlock().GetUnlockTime().IsValid() &&
		txEvent.GetWalletAddress() != "" {
		i.encodeAndPublish(txEvent)
	} else {
		log.Sugar.Errorf("Failed to parse begin_unlock event: %v", txEvent.String())
	}
}

func (i *Indexer) encodeAndPublish(msg *indexEvent.TxEvent) {
	encoded, err := proto.Marshal(msg)
	if err != nil {
		log.Sugar.Error(err)
	}
	log.Sugar.Infof("Publishing: %v", encoded)
}

type SyncStatus struct {
	Height       int64
	LatestHeight int64
}

func (i *Indexer) getSyncStatus(baseUrl string, encodingConfig noapptest.TestEncodingConfig, apiClient indexerpbconnect.IndexerServiceClient) SyncStatus {
	log.Sugar.Info("Getting sync status")
	apiResponse, err := apiClient.GetHeight(context.Background(), connect.NewRequest(&indexerpb.GetHeightRequest{ChainName: "Osmosis"}))
	if err != nil {
		log.Sugar.Panic(err)
	}

	var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/latest", baseUrl)
	resp, err := http.Get(url)
	if err != nil {
		log.Sugar.Panic(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Sugar.Panic(err)
	}
	var response cmtservice.GetLatestBlockResponse
	if err := encodingConfig.Codec.UnmarshalJSON(body, &response); err != nil {
		log.Sugar.Panic(err)
	}
	var height = apiResponse.Msg.GetHeight() + 1
	if height == 1 {
		height = response.GetBlock().GetHeader().Height
	}
	return SyncStatus{
		LatestHeight: response.GetBlock().GetHeader().Height,
		Height:       height,
	}
}

func (i *Indexer) StartIndexing() {
	var syncStatus = i.getSyncStatus(i.baseUrl, i.encodingConfig, i.grpcClient)
	log.Sugar.Infof("Starting indexing at height: %v", syncStatus.Height)
	for true {
		var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/%v", i.baseUrl, syncStatus.Height)
		var blockResponse cmtservice.GetBlockByHeightResponse
		status, err := GetAndDecode(url, i.encodingConfig, &blockResponse)
		if err != nil {
			// TODO: handle error based on status code
			if status == 400 {
				log.Sugar.Infof("Block does not yet exist: %v", syncStatus.Height)
			} else {
				log.Sugar.Panicf("Failed to get block: %v %v", status, err)
			}
		} else {
			i.handleBlock(&blockResponse)
			_, err := i.grpcClient.UpdateHeight(context.Background(),
				connect.NewRequest(
					&indexerpb.UpdateHeightRequest{ChainName: "Osmosis", Height: syncStatus.Height},
				),
			)
			if err != nil {
				log.Sugar.Panic(err)
			}
			syncStatus.Height++
		}
		if syncStatus.Height >= syncStatus.LatestHeight {
			time.Sleep(1 * time.Second)
		}
	}
}
