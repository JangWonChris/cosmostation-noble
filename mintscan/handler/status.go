package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"

	"go.uber.org/zap"
)

var (
	mu           *sync.RWMutex
	latestStatus *model.ResultStatus
)

func init() {
	mu = new(sync.RWMutex)
	latestStatus = new(model.ResultStatus)
}

// GetStatus returns ResultStatus, which includes current network status
// TODO: 1. Circulating Supply is equal to Total Supply Tokens? If it is, remove which one?
// 		 2. Does this API need to request RPC rather than REST APIs?
func GetStatus(rw http.ResponseWriter, r *http.Request) {
	mu.RLock()
	result := *latestStatus
	mu.RUnlock()
	model.Respond(rw, &result)
	return
}

// SetStatus store the latest status in memory periodically
func SetStatus() error {
	if s == nil {
		return fmt.Errorf("Session is not initialized")
	}

	poolResp, err := s.client.HandleResponseHeight("/staking/pool")
	if err != nil {
		zap.L().Error("failed to get staking pool", zap.Error(err))
		return err
	}

	var pool model.Pool
	err = json.Unmarshal(poolResp.Result, &pool)
	if err != nil {
		zap.L().Error("failed to unmarshal pool", zap.Error(err))
		return err
	}

	tsResp, err := s.client.HandleResponseHeight("/supply/total")
	if err != nil {
		zap.L().Error("failed to get supply total", zap.Error(err))
		return err
	}

	var coins []model.Coin
	err = json.Unmarshal(tsResp.Result, &coins)
	if err != nil {
		zap.L().Error("failed to unmarshal coin", zap.Error(err))
		return err
	}

	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	bondedValsNum, _ := s.db.CountValidatorNumByStatus(model.BondedValidatorStatus)
	unbondingValsNum, _ := s.db.CountValidatorNumByStatus(model.UnbondingValidatorStatus)
	unbondedValsNum, _ := s.db.CountValidatorNumByStatus(model.UnbondedValidatorStatus)
	totalTxsNum := s.db.QueryTotalTransactionNum()

	status, err := s.client.GetStatus()
	if err != nil {
		zap.L().Error("failed to get chain status", zap.Error(err))
		return err
	}

	// Query two latest blocks to calculate block time.
	latestTwoBlocks, _ := s.db.QueryLastestTwoBlocks()
	if len(latestTwoBlocks) <= 1 {
		zap.L().Debug("failed to query two latest blocks", zap.Any("blocks", latestTwoBlocks))
		return err
	}

	lastBlocktime := latestTwoBlocks[0].Timestamp.UTC()
	secondLastBlocktime := latestTwoBlocks[1].Timestamp.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime).Seconds()

	mu.Lock()
	latestStatus = &model.ResultStatus{
		ChainID:                status.NodeInfo.Network,
		BlockHeight:            status.SyncInfo.LatestBlockHeight,
		BlockTime:              blockTime,
		TotalTxsNum:            totalTxsNum,
		TotalValidatorNum:      bondedValsNum + unbondingValsNum + unbondedValsNum,
		JailedValidatorNum:     unbondingValsNum + unbondedValsNum,
		UnjailedValidatorNum:   bondedValsNum,
		TotalSupplyTokens:      coins,
		TotalCirculatingTokens: coins, // TODO: should be how we discuss with CoinGecko (Total Supply - Vesting Amount)
		BondedTokens:           bondedTokens,
		NotBondedTokens:        notBondedTokens,
		Timestamp:              status.SyncInfo.LatestBlockTime,
	}
	mu.Unlock()

	return nil
}
