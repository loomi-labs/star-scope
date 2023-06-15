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

type Indexer struct {
	chain          *indexerpb.IndexingChain
	encodingConfig EncodingConfig
	kafkaProducer  *kafka.KafkaProducer
	txHandler      TxHandler
}

type TxHandler interface {
	HandleTxs(txs [][]byte) (*indexerpb.HandleTxsResponse, error)
}

func NewIndexer(chain *indexerpb.IndexingChain, encodingConfig EncodingConfig, kafkaBrokers []string, txHandler TxHandler) Indexer {
	return Indexer{
		chain:          chain,
		encodingConfig: encodingConfig,
		kafkaProducer:  kafka.NewKafkaProducer(kafka.WalletEvents, kafkaBrokers...),
		txHandler:      txHandler,
	}
}

type TxHelper struct {
	chainInfo      *indexerpb.IndexingChain
	encodingConfig EncodingConfig
}

func NewTxHelper(chainInfo *indexerpb.IndexingChain, encodingConfig EncodingConfig) TxHelper {
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
	recordedAt            time.Time
	lastBlockTimestamps   []time.Time
	latestHeightUpdatedAt time.Time
}

func getLatestHeight(encodingConfig EncodingConfig, chainInfo *indexerpb.IndexingChain) (uint64, error) {
	log.Sugar.Debugf("Getting latest height of chain %v", chainInfo.Name)

	var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/latest", chainInfo.RestEndpoint)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, errors.New(fmt.Sprintf("Status code: %v", resp.StatusCode))
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
	return uint64(response.GetBlock().GetHeader().Height), nil
}

func (i *Indexer) newSyncStatus(chain *indexerpb.IndexingChain) SyncStatus {
	var handledMessageTypes = make(map[string]struct{})
	for _, msgType := range chain.HandledMessageTypes {
		handledMessageTypes[msgType] = struct{}{}
	}
	var unhandledMessageTypes = make(map[string]struct{})
	for _, msgType := range chain.UnhandledMessageTypes {
		unhandledMessageTypes[msgType] = struct{}{}
	}

	// try to get latest height 5 times and then give up
	var latestHeight uint64
	var err error
	for j := 0; j < 5; j++ {
		latestHeight, err = getLatestHeight(i.encodingConfig, i.chain)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Sugar.Panicf("Error getting latest height for chain %v: %s", chain.Name, err)
	}
	var height = chain.IndexingHeight
	if height == 0 {
		height = latestHeight
	}

	return SyncStatus{
		ChainId:               chain.Id,
		Height:                height,
		LatestHeight:          latestHeight,
		HandledMessageTypes:   handledMessageTypes,
		UnhandledMessageTypes: unhandledMessageTypes,
		latestHeightUpdatedAt: time.Now(),
	}
}

func (s *SyncStatus) recordBlock(blockResponse cmtservice.GetBlockByHeightResponse) {
	if blockResponse.GetBlock() == nil {
		log.Sugar.Panicf("Block response is nil")
	}
	s.lastBlockTimestamps = append(s.lastBlockTimestamps, blockResponse.GetBlock().GetHeader().Time)
	if len(s.lastBlockTimestamps) > 10 {
		s.lastBlockTimestamps = s.lastBlockTimestamps[1:]
	}
	s.recordedAt = time.Now()
}

const startSleepTime = 5 * time.Second
const minSleepTime = 200 * time.Millisecond

func (s *SyncStatus) getSleepDuration() time.Duration {
	if s.LatestHeight > s.Height {
		// if we are far behind, don't sleep at all
		if s.LatestHeight > s.Height+10 {
			return 0
		}
		return minSleepTime
	}
	if len(s.lastBlockTimestamps) < 2 {
		return startSleepTime
	}
	var totalDuration time.Duration
	for i := 1; i < len(s.lastBlockTimestamps); i++ {
		totalDuration += s.lastBlockTimestamps[i].Sub(s.lastBlockTimestamps[i-1])
	}
	var avgDuration = totalDuration / time.Duration(len(s.lastBlockTimestamps)-1)
	return avgDuration - time.Since(s.recordedAt)
}

func (s *SyncStatus) updateSyncStatus(updateChannel chan SyncStatus, encodingConfig EncodingConfig, chainInfo *indexerpb.IndexingChain) {
	updateChannel <- *s
	s.Height++
	if s.latestHeightUpdatedAt.Add(1 * time.Minute).Before(time.Now()) {
		var latestHeight, err = getLatestHeight(encodingConfig, chainInfo)
		if err != nil {
			log.Sugar.Errorf("Error getting latest height for chain %v: %s", chainInfo.Name, err)
		} else {
			s.LatestHeight = latestHeight
			s.latestHeightUpdatedAt = time.Now()
		}
	}
}

func (i *Indexer) logStats(result *indexerpb.HandleTxsResponse, syncStatus SyncStatus, getBlockRequestDuration time.Duration, handleBlockDuration time.Duration, sleepDuration time.Duration) {
	var behindText = ""
	var height = syncStatus.Height - 1
	if syncStatus.LatestHeight > height {
		behindText = fmt.Sprintf(" (%v behind latest)", syncStatus.LatestHeight-height)
	}
	if result != nil {
		var total = result.CountProcessed + result.CountSkipped
		sleepTrunc := sleepDuration.Truncate(time.Millisecond * 100)
		log.Sugar.Infof("%-15sBlock %v%v\tTotal: %v\tSkipped: %v\tProcessed: %v\tKafka msgs: %v\tGet block: %v\tHandle block: %v\tSleep: %s",
			i.chain.Name, height,
			behindText, total, result.CountSkipped, result.CountProcessed, len(result.ProtoMessages), getBlockRequestDuration, handleBlockDuration, sleepTrunc.String())
	}
}

func (i *Indexer) StartIndexing(updateChannel chan SyncStatus, stopChannel chan struct{}) {
	var syncStatus = i.newSyncStatus(i.chain)
	log.Sugar.Infof("%-15sStart indexing at height: %v", i.chain.Name, syncStatus.Height)
	for {
		select {
		case <-stopChannel:
			log.Sugar.Infof("%-15sStop indexing at height: %v", i.chain.Name, syncStatus.Height)
			return
		default:
			var url = fmt.Sprintf("%v/cosmos/base/tendermint/v1beta1/blocks/%v", i.chain.RestEndpoint, syncStatus.Height)
			var blockResponse cmtservice.GetBlockByHeightResponse
			var startGetBlockRequest = time.Now()
			status, err := GetAndDecode(url, i.encodingConfig, &blockResponse)
			var getBlockRequestDuration = time.Since(startGetBlockRequest)
			var handleBlockDuration time.Duration
			var result *indexerpb.HandleTxsResponse
			if err != nil {
				if status == 400 {
					log.Sugar.Debugf("%-15sBlock does not yet exist: %v", i.chain.Name, syncStatus.Height)
				} else if status > 400 && status < 500 {
					log.Sugar.Warnf("%-15sFailed to get block: %v %v", i.chain.Name, status, err)
				} else if status >= 500 {
					log.Sugar.Warnf("%-15sFailed to get block: %v %v", i.chain.Name, status, err)
				} else {
					log.Sugar.Errorf("%-15sFailed to get block: %v %v", i.chain.Name, status, err)
				}
				time.Sleep(200 * time.Millisecond)
			} else {
				syncStatus.recordBlock(blockResponse)
				var startHandleBlock = time.Now()
				result, err = i.handleBlock(&blockResponse, syncStatus)
				handleBlockDuration = time.Since(startHandleBlock)
				if err != nil {
					log.Sugar.Errorf("%-15sFailed to handle block %v (try again): %v", i.chain.Name, syncStatus.Height, err)
					time.Sleep(200 * time.Millisecond)
				} else {
					syncStatus.updateSyncStatus(updateChannel, i.encodingConfig, i.chain)
				}
			}
			var sleepDuration = syncStatus.getSleepDuration()
			i.logStats(result, syncStatus, getBlockRequestDuration, handleBlockDuration, sleepDuration)
			time.Sleep(sleepDuration)
		}
	}
}
