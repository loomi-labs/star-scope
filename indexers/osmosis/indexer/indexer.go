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
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibcChannel "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/osmosis-labs/osmosis/osmoutils/noapptest"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	"github.com/shifty11/blocklog-backend/indexers/osmosis/client"
	"github.com/shifty11/go-logger/log"
	"io"
	"net/http"
	"time"
)

func handleBlock(blockResponse *cmtservice.GetBlockByHeightResponse, config noapptest.TestEncodingConfig) {
	log.Sugar.Infof("handleBlock: %v", blockResponse.GetBlock().GetHeader().Height)
	var data = blockResponse.GetBlock().GetData()
	var txs = data.GetTxs()
	for _, tx := range txs {
		txDecoded, err := config.TxConfig.TxDecoder()(tx)
		if err != nil {
			//log.Sugar.Error(err)
			log.Sugar.Info("Failed to decode txDecoded")
			continue
		}
		for _, anyMsg := range txDecoded.GetMsgs() {
			switch msg := anyMsg.(type) {
			case *banktypes.MsgSend:
				handleMsgSend(msg)
			case *banktypes.MsgMultiSend:
				handleMsgMultiSend(msg)
			case *transfertypes.MsgTransfer:
				handleTransferMsg(msg, tx)
			case *ibcChannel.MsgRecvPacket:
				handleIbcMsg(msg)
			case *lockuptypes.MsgBeginUnlockingAll:
				handleMsgBeginUnlockingAll(msg)
			case *lockuptypes.MsgBeginUnlocking:
				handleMsgBeginUnlocking(msg, tx)
			default:
				log.Sugar.Infof("Unknown message type")
			}
		}
	}
}

func getTxResult(tx []byte) txtypes.GetTxResponse {
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

func handleMsgSend(msg *banktypes.MsgSend) {
	log.Sugar.Infof("MsgSend: %v", msg.String())
}

func handleMsgMultiSend(msg *banktypes.MsgMultiSend) {
	log.Sugar.Infof("MsgMultiSend: %v", msg.String())
}

func handleTransferMsg(msg *transfertypes.MsgTransfer, tx []byte) {
	log.Sugar.Infof("MsgTransfer: %v", msg.String())
}

func handleIbcMsg(msg *ibcChannel.MsgRecvPacket) {
	log.Sugar.Infof("MsgRecvPacket: %v", msg.String())
}

func handleMsgBeginUnlockingAll(msg *lockuptypes.MsgBeginUnlockingAll) {
	log.Sugar.Infof("MsgBeginUnlockingAll: %v", msg.String())
}

func handleMsgBeginUnlocking(msg *lockuptypes.MsgBeginUnlocking, tx []byte) {
	log.Sugar.Infof("MsgBeginUnlocking: %v", msg.String())
	var resp = getTxResult(tx)
	if resp.GetTxResponse().Code == 0 {
		for _, event := range resp.GetTxResponse().Events {
			if event.Type == "begin_unlock" {
				for _, attribute := range event.Attributes {
					if string(attribute.GetKey()) == "unlock_time" {
						log.Sugar.Infof("unlock_time: %v", string(attribute.GetValue()))
					}
				}
			}
		}
	}
}

type SyncStatus struct {
	Height       int64
	LatestHeight int64
}

func getSyncStatus(baseUrl string, encodingConfig noapptest.TestEncodingConfig, apiClient indexerpbconnect.IndexerServiceClient) SyncStatus {
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

func StartIndexing(baseUrl string) {
	encodingConfig := GetEncodingConfig()
	apiClient := client.GetClient()

	var syncStatus = getSyncStatus(baseUrl, encodingConfig, apiClient)
	for true {
		var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/%v", baseUrl, syncStatus.Height)
		var blockResponse cmtservice.GetBlockByHeightResponse
		status, err := GetAndDecode(url, encodingConfig, &blockResponse)
		if err != nil {
			// TODO: handle error based on status code
			if status == 400 {
				log.Sugar.Infof("Block does not yet exist: %v", syncStatus.Height)
			} else {
				log.Sugar.Panicf("Failed to get block: %v %v", status, err)
			}
		} else {
			handleBlock(&blockResponse, encodingConfig)
			_, err := apiClient.UpdateHeight(context.Background(), connect.NewRequest(&indexerpb.UpdateHeightRequest{ChainName: "Osmosis", Height: syncStatus.Height}))
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
