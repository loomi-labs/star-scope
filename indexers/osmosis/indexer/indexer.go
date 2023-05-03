package indexer

import (
	"buf.build/gen/go/rapha/blocklog/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/rapha/blocklog/protocolbuffers/go/grpc/indexer/indexerpb"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bufbuild/connect-go"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcChannel "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/osmosis-labs/osmosis/osmoutils/noapptest"
	lockuptypes "github.com/osmosis-labs/osmosis/v15/x/lockup/types"
	"github.com/shifty11/blocklog-backend/indexers/osmosis/client"
	"github.com/shifty11/go-logger/log"
	"github.com/tendermint/tendermint/abci/types"
	"io"
	"net/http"
	"time"
)

type ChainInfo struct {
	ChainName string
}

type Indexer struct {
	baseUrl        string
	chainInfo      ChainInfo
	grpcClient     indexerpbconnect.IndexerServiceClient
	encodingConfig noapptest.TestEncodingConfig
	kafkaProducer  *KafkaProducer
}

func NewIndexer(baseUrl string, kafkaAddresses []string) Indexer {
	return Indexer{
		baseUrl: baseUrl,
		chainInfo: ChainInfo{
			ChainName: "osmosis",
		},
		grpcClient:     client.GetClient(),
		encodingConfig: GetEncodingConfig(),
		kafkaProducer:  NewKafkaProducer(kafkaAddresses...),
	}
}

func (i *Indexer) handleBlock(blockResponse *cmtservice.GetBlockByHeightResponse) {
	var data = blockResponse.GetBlock().GetData()
	var txs = data.GetTxs()
	var cntSkipped = 0
	var cntMsgs = 0
	var protoMsgs = make([][]byte, 0)
	for _, tx := range txs {
		txDecoded, err := i.encodingConfig.TxConfig.TxDecoder()(tx)
		if err != nil {
			log.Sugar.Errorf("Error decoding tx: %v", err)
			break
		}
		cntMsgs += len(txDecoded.GetMsgs())
		for _, anyMsg := range txDecoded.GetMsgs() {
			var protoMsg []byte
			switch msg := anyMsg.(type) {
			case *banktypes.MsgSend:
				protoMsg, err = i.handleMsgSend(msg, tx)
			case *banktypes.MsgMultiSend:
				i.handleMsgMultiSend(msg, tx, blockResponse.GetBlock().GetHeader().Height)
			case *ibcChannel.MsgRecvPacket:
				protoMsg, err = i.handleMsgRecvPacket(msg, tx)
			case *lockuptypes.MsgBeginUnlockingAll:
				i.handleMsgBeginUnlockingAll(msg, tx, blockResponse.GetBlock().GetHeader().Height)
			case *lockuptypes.MsgBeginUnlocking:
				protoMsg, err = i.handleMsgBeginUnlocking(msg, tx)
			default:
				cntSkipped++
			}
			if err != nil {
				log.Sugar.Errorf("Error handling msg: %v", err)
				break
			}
			if protoMsg != nil {
				protoMsgs = append(protoMsgs, protoMsg)
			}
		}
		if err != nil {
			break
		}
	}
	if len(protoMsgs) > 0 {
		i.kafkaProducer.Produce(protoMsgs)
	}
	var cntProcessed = cntMsgs - cntSkipped
	log.Sugar.Debugf("Block %v\tTotal: %v\tSkipped: %v\tProcessed: %v\tKafka msgs: %v",
		blockResponse.GetBlock().GetHeader().Height, cntMsgs, cntSkipped, cntProcessed, len(protoMsgs))
}

func (i *Indexer) getTxResult(tx []byte) (*txtypes.GetTxResponse, error) {
	hash := sha256.Sum256(tx)
	hashString := hex.EncodeToString(hash[:])

	var url = fmt.Sprintf("%v/cosmos/tx/v1beta1/txs/%v", i.baseUrl, hashString)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code: %v", resp.StatusCode))
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	encodingConfig := GetEncodingConfig()
	var txResponse txtypes.GetTxResponse
	if err := encodingConfig.Codec.UnmarshalJSON(body, &txResponse); err != nil {
		return nil, err
	}
	return &txResponse, nil
}

func (i *Indexer) getTxEvents(tx []byte) ([]types.Event, error) {
	resp, err := i.getTxResult(tx)
	if err != nil {
		return nil, err
	}
	if resp.GetTxResponse().Code == 0 {
		//log.Sugar.Debugf("Tx %v was successful", resp.GetTxResponse().TxHash)
		return resp.GetTxResponse().Events, nil
	}
	return nil, nil
}

func (i *Indexer) wasTxSuccessful(tx []byte) (bool, error) {
	events, err := i.getTxEvents(tx)
	if err != nil {
		return false, err
	}
	return len(events) > 0, nil
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
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Failed to get latest block: %v", resp.StatusCode))
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

func (i *Indexer) updateHeight(syncStatus *SyncStatus) {
	_, err := i.grpcClient.UpdateHeight(context.Background(),
		connect.NewRequest(
			&indexerpb.UpdateHeightRequest{ChainName: "Osmosis", Height: syncStatus.Height},
		),
	)
	if err != nil {
		panic(err)
	}
	syncStatus.Height++
}

func (i *Indexer) StartIndexing() {
	var syncStatus = i.getSyncStatus(i.baseUrl, i.encodingConfig, i.grpcClient)
	log.Sugar.Infof("Start indexing at height: %v", syncStatus.Height)
	for true {
		var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/%v", i.baseUrl, syncStatus.Height)
		var blockResponse cmtservice.GetBlockByHeightResponse
		status, err := GetAndDecode(url, i.encodingConfig, &blockResponse)
		if err != nil {
			// TODO: handle error based on status code
			if status == 400 {
				log.Sugar.Infof("Block does not yet exist: %v", syncStatus.Height)
			} else if status >= 500 {
				log.Sugar.Warnf("Failed to get block: %v %v", status, err)
			} else {
				log.Sugar.Panicf("Failed to get block: %v %v", status, err)
			}
		} else {
			i.handleBlock(&blockResponse)
			i.updateHeight(&syncStatus)
		}
		if syncStatus.Height >= syncStatus.LatestHeight {
			time.Sleep(1 * time.Second)
		}
	}
}
