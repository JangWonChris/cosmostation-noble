package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// SetStatus store the latest status in memory using mutex every 5 seconds.
func SetStatus() error {
	if s == nil {
		return fmt.Errorf("Session is not initialized")
	}

	stakingQueryClient := stakingtypes.NewQueryClient(s.client.GetCLIContext())
	pool, err := stakingQueryClient.Pool(context.Background(), &stakingtypes.QueryPoolRequest{})
	if err != nil {
		zap.L().Error("failed to get staking pool", zap.Error(err))
		return err
	}

	bankQueryClient := banktypes.NewQueryClient(s.client.GetCLIContext())
	coins, err := bankQueryClient.TotalSupply(context.Background(), &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		zap.L().Error("failed to get supply total", zap.Error(err))
		return err
	}

	notBondedTokens, _ := strconv.ParseFloat(pool.Pool.NotBondedTokens.String(), 64)
	bondedTokens, _ := strconv.ParseFloat(pool.Pool.BondedTokens.String(), 64)
	bondedValsNum, _ := s.db.CountValidatorsByStatus(int(stakingtypes.Bonded))
	unbondingValsNum, _ := s.db.CountValidatorsByStatus(int(stakingtypes.Unbonding))
	unbondedValsNum, _ := s.db.CountValidatorsByStatus(int(stakingtypes.Unbonded))
	totalTxsNum := s.db.QueryTotalTransactionNum()

	status, err := s.client.RPC.GetStatus()
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
		TotalSupplyTokens:      *coins,
		TotalCirculatingTokens: *coins, // TODO: should be how we discuss with CoinGecko (Total Supply - Vesting Amount)
		BondedTokens:           bondedTokens,
		NotBondedTokens:        notBondedTokens,
		Timestamp:              status.SyncInfo.LatestBlockTime,
	}
	mu.Unlock()

	return nil
}
