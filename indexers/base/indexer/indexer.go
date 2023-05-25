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
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"time"
)

type ChainInfo struct {
	ChainId      uint64
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
		kafkaProducer:  kafka.NewKafkaProducer(kafka.IndexEventsTopic, config.KafkaBrokers...),
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

func (h *TxHelper) GetTxResponse(tx []byte) (*sdktypes.TxResponse, *timestamppb.Timestamp, error) {
	resp, err := h.GetTxResult(tx)
	if err != nil {
		return nil, nil, err
	}
	if resp.GetTxResponse().Code == 0 {
		timestamp, err := h.GetTxTimestamp(resp.GetTxResponse())
		if err != nil {
			log.Sugar.Errorf("Error getting tx timestamp: %s", err)
			return resp.GetTxResponse(), nil, nil
		}
		return resp.GetTxResponse(), timestamppb.New(timestamp), nil
	}
	return nil, nil, nil
}

func (h *TxHelper) WasTxSuccessful(tx []byte) (bool, *timestamppb.Timestamp, error) {
	txResponse, timestamp, err := h.GetTxResponse(tx)
	if err != nil {
		return false, nil, err
	}
	if txResponse == nil {
		return false, nil, nil
	}
	return len(txResponse.Events) > 0, timestamp, nil
}

func (h *TxHelper) GetTxTimestamp(txResponse *sdktypes.TxResponse) (time.Time, error) {
	return time.Parse(time.RFC3339, txResponse.Timestamp)
}

func (i *Indexer) handleBlock(blockResponse *cmtservice.GetBlockByHeightResponse, syncStatus SyncStatus) (*indexerpb.HandleTxsResponse, error) {
	var data = blockResponse.GetBlock().GetData()
	var txs = data.GetTxs()

	var result, err = i.txHandler.HandleTxs(txs)
	if err != nil {
		return nil, err
	}

	if len(result.GetProtoMessages()) > 0 {
		i.kafkaProducer.Produce(result.ProtoMessages)
	}

	if result.HandledMessageTypes != nil {
		for _, msgType := range result.HandledMessageTypes {
			syncStatus.HandledMessageTypes[msgType] = struct{}{}
		}
	}

	if result.UnhandledMessageTypes != nil {
		for _, msgType := range result.UnhandledMessageTypes {
			syncStatus.UnhandledMessageTypes[msgType] = struct{}{}
		}
	}
	return result, nil
}

type SyncStatus struct {
	ChainId               uint64
	Height                uint64
	LatestHeight          uint64
	HandledMessageTypes   map[string]struct{}
	UnhandledMessageTypes map[string]struct{}
}

func (i *Indexer) getLatestHeight(syncStatus *SyncStatus) {
	// get latest height if we don't have it or if we are more than 100 blocks behind
	if syncStatus.LatestHeight == 0 || syncStatus.LatestHeight < syncStatus.Height-100 {
		log.Sugar.Debugf("Getting latest height of chain %v", i.chainInfo.Name)

		var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/latest", i.chainInfo.RestEndpoint)
		resp, err := http.Get(url)
		if err != nil {
			log.Sugar.Panic(err)
		}
		if resp.StatusCode != 200 {
			log.Sugar.Warnf("Failed to get latest block: %v", resp.StatusCode)
			time.Sleep(5 * time.Second)
			i.getLatestHeight(syncStatus)
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
}

func (i *Indexer) updateSyncStatus(syncStatus *SyncStatus, updateChannel chan SyncStatus) {
	updateChannel <- *syncStatus
	syncStatus.Height++
}

const startSleepTime = 5 * time.Second
const incrementSleepTime = 500 * time.Millisecond
const minSleepTime = 200 * time.Millisecond
const maxSleepTime = 10 * time.Second

func setSleepTime(sleepTime time.Duration, attempts int) time.Duration {
	if attempts == 1 {
		sleepTime -= incrementSleepTime
		if sleepTime < minSleepTime {
			sleepTime = minSleepTime
		}
	} else if attempts > 2 {
		sleepTime += incrementSleepTime
		if sleepTime > maxSleepTime {
			sleepTime = maxSleepTime
		}
	}
	return sleepTime
}

func (i *Indexer) logStats(result *indexerpb.HandleTxsResponse, syncStatus SyncStatus, getBlockRequestDuration time.Duration, handleBlockDuration time.Duration, sleepTime time.Duration, catchUp bool) {
	var behindText = ""
	var height = syncStatus.Height - 1
	if syncStatus.LatestHeight > height {
		behindText = fmt.Sprintf(" (%v behind latest)", syncStatus.LatestHeight-height)
	}
	if result != nil {
		var total = result.CountProcessed + result.CountSkipped
		var sleepTimeText = "0"
		if !catchUp {
			sleepTimeText = sleepTime.String()
		}
		log.Sugar.Infof("%-15sBlock %v%v\tTotal: %v\tSkipped: %v\tProcessed: %v\tKafka msgs: %v\tGet block: %v\tHandle block: %v\tSleep: %v",
			i.chainInfo.Name, height,
			behindText, total, result.CountSkipped, result.CountProcessed, len(result.ProtoMessages), getBlockRequestDuration, handleBlockDuration, sleepTimeText)
	}
}

func (i *Indexer) StartIndexing(syncStatus SyncStatus, updateChannel chan SyncStatus) {
	i.getLatestHeight(&syncStatus)
	var sleepTime = startSleepTime
	var catchUp = syncStatus.Height < syncStatus.LatestHeight
	var attempts = 1
	log.Sugar.Infof("%-15sStart indexing at height: %v", i.chainInfo.Name, syncStatus.Height)
	for true {
		var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/%v", i.chainInfo.RestEndpoint, syncStatus.Height)
		var blockResponse cmtservice.GetBlockByHeightResponse
		var startGetBlockRequest = time.Now()
		status, err := GetAndDecode(url, i.encodingConfig, &blockResponse)
		var getBlockRequestDuration = time.Since(startGetBlockRequest)
		var handleBlockDuration time.Duration
		var result *indexerpb.HandleTxsResponse
		if err != nil {
			if status == 400 {
				log.Sugar.Debugf("%-15sBlock does not yet exist: %v", i.chainInfo.Name, syncStatus.Height)
				attempts++
			} else if status > 400 && status < 500 {
				log.Sugar.Warnf("%-15sFailed to get block: %v %v", i.chainInfo.Name, status, err)
			} else if status >= 500 {
				log.Sugar.Warnf("%-15sFailed to get block: %v %v", i.chainInfo.Name, status, err)
			} else {
				log.Sugar.Errorf("%-15sFailed to get block: %v %v", i.chainInfo.Name, status, err)
			}
		} else {
			var startHandleBlock = time.Now()
			result, err = i.handleBlock(&blockResponse, syncStatus)
			handleBlockDuration = time.Since(startHandleBlock)
			if err != nil {
				log.Sugar.Errorf("%-15sFailed to handle block %v (try again): %v", i.chainInfo.Name, syncStatus.Height, err)
				time.Sleep(200 * time.Millisecond)
			} else {
				i.updateSyncStatus(&syncStatus, updateChannel)
				i.getLatestHeight(&syncStatus)
				attempts = 1
			}
		}
		if !catchUp {
			sleepTime = setSleepTime(sleepTime, attempts)
			time.Sleep(sleepTime)
		} else {
			catchUp = syncStatus.Height < syncStatus.LatestHeight
		}
		i.logStats(result, syncStatus, getBlockRequestDuration, handleBlockDuration, sleepTime, catchUp)
	}
}
