package indexer

import (
	"buf.build/gen/go/loomi-labs/star-scope/protocolbuffers/go/grpc/indexer/indexerpb"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
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
	encodingConfig EncodingConfig
	kafkaProducer  *kafka.KafkaProducer
	txHandler      TxHandler
}

type TxHandler interface {
	HandleTxs(txs [][]byte) (*indexerpb.HandleTxsResponse, error)
}

type Config struct {
	ChainInfo      ChainInfo
	KafkaBrokers   []string
	EncodingConfig EncodingConfig
	MessageHandler TxHandler
}

func NewIndexer(config Config) Indexer {
	return Indexer{
		chainInfo:      config.ChainInfo,
		encodingConfig: config.EncodingConfig,
		kafkaProducer:  kafka.NewKafkaProducer(config.KafkaBrokers...),
		txHandler:      config.MessageHandler,
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

func (h *TxHelper) GetTxTimestamp(txResponse *sdktypes.TxResponse) (time.Time, error) {
	return time.Parse(time.RFC3339, txResponse.Timestamp)
}

func (i *Indexer) handleBlock(blockResponse *cmtservice.GetBlockByHeightResponse, syncStatus SyncStatus) error {
	var data = blockResponse.GetBlock().GetData()
	var txs = data.GetTxs()

	var result, err = i.txHandler.HandleTxs(txs)
	if err != nil {
		return err
	}

	if len(result.GetProtoMessages()) > 0 {
		i.kafkaProducer.Produce(result.ProtoMessages)
	}

	if result.UnhandledMessageTypes != nil {
		for _, msgType := range result.UnhandledMessageTypes {
			syncStatus.UnhandledMessageTypes[msgType] = struct{}{}
		}
	}

	var behindText = ""
	var behind = int(syncStatus.LatestHeight) - int(blockResponse.GetBlock().GetHeader().Height)
	if behind > 0 {
		behindText = fmt.Sprintf(" (%v behind latest)", behind)
	}
	var total = result.CountProcessed + result.CountSkipped
	log.Sugar.Infof("Block %v%v\tTotal: %v\tSkipped: %v\tProcessed: %v\tKafka msgs: %v",
		blockResponse.GetBlock().GetHeader().Height, behindText, total, result.CountSkipped, result.CountProcessed, len(result.ProtoMessages))
	return nil
}

type SyncStatus struct {
	ChainId               uint64
	Height                uint64
	LatestHeight          uint64
	UnhandledMessageTypes map[string]struct{}
}

func (i *Indexer) getLatestHeight(syncStatus *SyncStatus) {
	log.Sugar.Infof("Getting latest height of chain %v", i.chainInfo.Name)

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
	syncStatus.LatestHeight = uint64(response.GetBlock().GetHeader().Height)
	if syncStatus.Height == 0 {
		syncStatus.Height = syncStatus.LatestHeight
	}
}

func (i *Indexer) updateSyncStatus(syncStatus *SyncStatus, updateChannel chan SyncStatus) {
	var unhandledMessageTypes []string
	for msgType := range syncStatus.UnhandledMessageTypes {
		unhandledMessageTypes = append(unhandledMessageTypes, msgType)
	}
	updateChannel <- *syncStatus
	syncStatus.Height++
}

func (i *Indexer) StartIndexing(syncStatus SyncStatus, updateChannel chan SyncStatus) {
	i.getLatestHeight(&syncStatus)
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
			err = i.handleBlock(&blockResponse, syncStatus)
			if err != nil {
				log.Sugar.Errorf("Failed to handle block %v (try again): %v", syncStatus.Height, err)
				time.Sleep(200 * time.Millisecond)
			} else {
				i.updateSyncStatus(&syncStatus, updateChannel)
			}
		}
		if !catchUp {
			time.Sleep(1 * time.Second)
		} else {
			catchUp = syncStatus.Height < syncStatus.LatestHeight
		}
	}
}
