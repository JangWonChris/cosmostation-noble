package exporter

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/models"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"

	"go.uber.org/zap"
)

// SaveNetworkStats1H saves network statistics every hour.
func (ex *Exporter) SaveNetworkStats1H() {
	// Get the current state of the staking pool.
	poolResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/pool")
	if err != nil {
		zap.S().Errorf("failed to get the current state of the staking pool: %s", err)
		return
	}

	var pool models.Pool
	err = json.Unmarshal(poolResp.Result, &pool)
	if err != nil {
		zap.S().Errorf("failed to unmarshal pool: %s", err)
		return
	}

	// Get the current minting inflation value.
	inflation, err := ex.client.GetInflation()
	if err != nil {
		zap.S().Errorf("failed to get inflation value: %s", err)
		return
	}

	// Get the current total supply.
	totalSupplyResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/supply/total")
	if err != nil {
		zap.S().Errorf("failed to get the current total supply: %s", err)
		return
	}

	var coin []sdk.Coin
	err = json.Unmarshal(totalSupplyResp.Result, &coin)
	if err != nil {
		zap.S().Errorf("failed to unmarshal total supply: %s", err)
		return
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalSupplyTokens, _ := strconv.ParseFloat(coin[0].Amount.String(), 64)
	bondedRatio := bondedTokens / totalSupplyTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation.Result, 64)
	totalTxsNum := ex.db.QueryTotalTransactionNum()

	// Calculate the current block time.
	latestTwoBlocks, err := ex.db.QueryLatestTwoBlocks()
	if len(latestTwoBlocks) <= 1 {
		zap.S().Errorf("failed to query latest two blocks", err)
		return
	}

	lastBlocktime := latestTwoBlocks[0].Timestamp.UTC()
	secondLastBlocktime := latestTwoBlocks[1].Timestamp.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime)

	network := &schema.StatsNetwork1H{
		BlockTime:       blockTime.Seconds(),
		TotalSupply:     totalSupplyTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     totalTxsNum,
	}

	err = ex.db.InsertNetworkStats1H(network)
	if err != nil {
		zap.S().Errorf("failed to save network data: %s", err)
		return
	}

	zap.S().Info("successfully saved NetworkStats1H")
	return
}

// SaveNetworkStats1D saves network statistics every day.
func (ex *Exporter) SaveNetworkStats1D() {
	// Get the current state of the staking pool.
	poolResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/staking/pool")
	if err != nil {
		zap.S().Errorf("failed to get the current state of the staking pool: %s", err)
		return
	}

	var pool models.Pool
	err = json.Unmarshal(poolResp.Result, &pool)
	if err != nil {
		zap.S().Errorf("failed to unmarshal pool", err)
		return
	}

	// Get the current minting inflation value
	inflation, err := ex.client.GetInflation()
	if err != nil {
		zap.S().Errorf("failed to get inflation value", err)
		return
	}

	// Get the current total supply
	totalSupplyResp, err := ex.client.RequestAPIFromLCDWithRespHeight("/supply/total")
	if err != nil {
		zap.S().Errorf("failed to get the current total supply: %s", err)
		return
	}

	var coin []sdk.Coin
	err = json.Unmarshal(totalSupplyResp.Result, &coin)
	if err != nil {
		zap.S().Errorf("failed to unmarshal total supply: %s", err)
		return
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalSupplyTokens, _ := strconv.ParseFloat(coin[0].Amount.String(), 64)
	bondedRatio := bondedTokens / totalSupplyTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation.Result, 64)
	totalTxsNum := ex.db.QueryTotalTransactionNum()

	// Query two latest block numbers to calculate the current block time
	blocks, err := ex.db.QueryLatestTwoBlocks()
	if err != nil {
		zap.S().Errorf("failed to query latest two blocks: %s", err)
		return
	}

	lastBlocktime := blocks[0].Timestamp.UTC()
	secondLastBlocktime := blocks[1].Timestamp.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime)

	network := &schema.StatsNetwork1D{
		BlockTime:       blockTime.Seconds(),
		TotalSupply:     totalSupplyTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     totalTxsNum,
	}

	err = ex.db.InsertNetworkStats1D(network)
	if err != nil {
		zap.S().Errorf("failed to save network data: %s", err)
		return
	}

	zap.S().Info("successfully saved NetworkStats1D")
	return
}
