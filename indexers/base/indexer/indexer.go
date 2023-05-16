package indexer

import (
	"buf.build/gen/go/loomi-labs/star-scope/bufbuild/connect-go/grpc/indexer/indexerpb/indexerpbconnect"
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bufbuild/connect-go"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/loomi-labs/star-scope/indexers/base/client"
	"github.com/loomi-labs/star-scope/indexers/base/common"
	"github.com/loomi-labs/star-scope/indexers/base/kafka"
	"github.com/shifty11/go-logger/log"
	"io"
	"net/http"
	"time"
)

type ChainInfo struct {
	Path         string
	RestEndpoint string
	Name         string
}

type Indexer struct {
	chainInfo      ChainInfo
	grpcClient     indexerpbconnect.IndexerServiceClient
	encodingConfig EncodingConfig
	kafkaProducer  *kafka.KafkaProducer
	messageHandler MessageHandler
}

type MessageHandler interface {
	DecodeTx(tx []byte) (sdktypes.Tx, error)
	HandleMessage(msg sdktypes.Msg, tx []byte) ([]byte, error)
}

type IndexerConfig struct {
	ChainInfo      ChainInfo
	KafkaAddresses []string
	EncodingConfig EncodingConfig
	GrpcAuthToken  string
	GrpcEndpoint   string
}

func NewIndexer(config IndexerConfig, messageHandler MessageHandler) Indexer {
	return Indexer{
		chainInfo:      config.ChainInfo,
		grpcClient:     client.NewGrpcClient(config.GrpcEndpoint, config.GrpcAuthToken),
		encodingConfig: config.EncodingConfig,
		kafkaProducer:  kafka.NewKafkaProducer(config.KafkaAddresses...),
		messageHandler: messageHandler,
	}
}

type TxHelper struct {
	chainInfo      ChainInfo
	encodingConfig EncodingConfig
}

func NewTxHelper(chainInfo ChainInfo, encodingConfig EncodingConfig) TxHelper {
	return TxHelper{
		chainInfo:      chainInfo,
		encodingConfig: encodingConfig,
	}
}

func (h *TxHelper) GetTxResult(tx []byte) (*txtypes.GetTxResponse, error) {
	hash := sha256.Sum256(tx)
	hashString := hex.EncodeToString(hash[:])

	var url = fmt.Sprintf("%v/cosmos/tx/v1beta1/txs/%v", h.chainInfo.RestEndpoint, hashString)
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
	var txResponse txtypes.GetTxResponse
	if err := h.encodingConfig.Codec.UnmarshalJSON(body, &txResponse); err != nil {
		return nil, err
	}
	return &txResponse, nil
}

func (h *TxHelper) GetTxResponse(tx []byte) (*sdktypes.TxResponse, error) {
	resp, err := h.GetTxResult(tx)
	if err != nil {
		return nil, err
	}
	if resp.GetTxResponse().Code == 0 {
		return resp.GetTxResponse(), nil
	}
	return nil, nil
}

func (h *TxHelper) WasTxSuccessful(tx []byte) (bool, error) {
	txResponse, err := h.GetTxResponse(tx)
	if err != nil {
		return false, err
	}
	if txResponse == nil {
		return false, nil
	}
	return len(txResponse.Events) > 0, nil
}

func (i *Indexer) handleBlock(blockResponse *cmtservice.GetBlockByHeightResponse, syncStatus SyncStatus) error {
	var data = blockResponse.GetBlock().GetData()
	var txs = data.GetTxs()
	var cntSkipped = 0
	var cntMsgs = 0
	var protoMsgs = make([][]byte, 0)
	for _, tx := range txs {
		txDecoded, err := i.messageHandler.DecodeTx(tx)
		if err != nil {
			return err
		}
		if txDecoded == nil {
			continue
		}
		cntMsgs += len(txDecoded.GetMsgs())
		for _, anyMsg := range txDecoded.GetMsgs() {
			var protoMsg []byte
			protoMsg, err := i.messageHandler.HandleMessage(anyMsg, tx)
			if err != nil {
				return err
			}
			if protoMsg != nil {
				protoMsgs = append(protoMsgs, protoMsg)
			} else {
				syncStatus.UnhandledMessageTypes[fmt.Sprintf("%T", anyMsg)] = struct{}{}
				cntSkipped++
			}
		}
	}
	if len(protoMsgs) > 0 {
		i.kafkaProducer.Produce(protoMsgs)
	}
	var cntProcessed = cntMsgs - cntSkipped
	var behindText = ""
	var behind = syncStatus.LatestHeight - blockResponse.GetBlock().GetHeader().Height
	if behind > 0 {
		behindText = fmt.Sprintf(" (%v behind latest)", behind)
	}
	log.Sugar.Infof("Block %v%v\tTotal: %v\tSkipped: %v\tProcessed: %v\tKafka msgs: %v",
		blockResponse.GetBlock().GetHeader().Height, behindText, cntMsgs, cntSkipped, cntProcessed, len(protoMsgs))
	return nil
}

type SyncStatus struct {
	Height                int64
	LatestHeight          int64
	UnhandledMessageTypes map[string]struct{}
}

func (i *Indexer) getSyncStatus() SyncStatus {
	log.Sugar.Info("Getting sync status")
	apiResponse, err := i.grpcClient.GetHeight(context.Background(), connect.NewRequest(&indexerpb.GetHeightRequest{ChainName: "Osmosis"}))
	if err != nil {
		log.Sugar.Panic(err)
	}

	var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/latest", i.chainInfo.RestEndpoint)
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
	if err := i.encodingConfig.Codec.UnmarshalJSON(body, &response); err != nil {
		log.Sugar.Panic(err)
	}
	var height = apiResponse.Msg.GetHeight() + 1
	if height == 1 || common.GetEnvAsBool("FORCE_START_FROM_LATEST_BLOCK") {
		height = response.GetBlock().GetHeader().Height
	}
	return SyncStatus{
		LatestHeight:          response.GetBlock().GetHeader().Height,
		Height:                height,
		UnhandledMessageTypes: map[string]struct{}{},
	}
}

func (i *Indexer) updateHeight(syncStatus *SyncStatus) {
	// TODO: send unhandled message types
	var unhandledMessageTypes []string
	for msgType := range syncStatus.UnhandledMessageTypes {
		unhandledMessageTypes = append(unhandledMessageTypes, msgType)
	}

	_, err := i.grpcClient.UpdateSyncStatus(context.Background(),
		connect.NewRequest(
			&indexerpb.UpdateSyncStatusRequest{
				ChainName:             "Osmosis",
				Height:                syncStatus.Height,
				UnhandledMessageTypes: unhandledMessageTypes,
			},
		),
	)
	if err != nil {
		// TODO: handle error based on status code and retry if necessary
		panic(err)
	}
	syncStatus.Height++
}

func (i *Indexer) StartIndexing() {
	var syncStatus = i.getSyncStatus()
	var catchUp = syncStatus.Height < syncStatus.LatestHeight
	log.Sugar.Infof("Start indexing at height: %v", syncStatus.Height)
	for true {
		var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/%v", i.chainInfo.RestEndpoint, syncStatus.Height)
		var blockResponse cmtservice.GetBlockByHeightResponse
		status, err := GetAndDecode(url, i.encodingConfig, &blockResponse)
		if err != nil {
			if status == 400 {
				log.Sugar.Debugf("Block does not yet exist: %v", syncStatus.Height)
			} else if status > 400 && status < 500 {
				log.Sugar.Warnf("Failed to get block: %v %v", status, err)
			} else if status >= 500 {
				log.Sugar.Warnf("Failed to get block: %v %v", status, err)
			} else {
				log.Sugar.Panicf("Failed to get block: %v %v", status, err)
			}
		} else {
			err := i.handleBlock(&blockResponse, syncStatus)
			if err != nil {
				log.Sugar.Errorf("Failed to handle block (try again): %v", err)
				time.Sleep(200 * time.Millisecond)
			} else {
				i.updateHeight(&syncStatus)
			}
		}
		if !catchUp {
			time.Sleep(1 * time.Second)
		} else {
			catchUp = syncStatus.Height < syncStatus.LatestHeight
		}
	}
}
