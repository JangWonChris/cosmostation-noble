package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	resty "gopkg.in/resty.v1"
)

// SaveNetworkStats1H saves network statistics every hour
func (ses *StatsExporterService) SaveNetworkStats1H() {
	// requests the current state of the staking pool
	poolResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/pool")

	var pool types.Pool
	err := json.Unmarshal(types.ReadRespWithHeight(poolResp).Result, &pool)
	if err != nil {
		fmt.Printf("failed to unmarshal Pool: %v \n ", err)
	}

	// requests the current minting inflation value
	inflationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/minting/inflation")

	var inflation types.Inflation
	err = json.Unmarshal(inflationResp.Body(), &inflation)
	if err != nil {
		fmt.Printf("failed to unmarshal Inflation: %v \n ", err)
	}

	// requests the current total supply
	totalSupplyResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/supply/total")

	var coin []sdk.Coin
	err = json.Unmarshal(types.ReadRespWithHeight(totalSupplyResp).Result, &coin)
	if err != nil {
		fmt.Printf("failed to unmarshal Coin: %v \n", err)
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalSupplyTokens, _ := strconv.ParseFloat(coin[0].Amount.String(), 64)
	bondedRatio := bondedTokens / totalSupplyTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation.Result, 64)

	// query two latest block numbers to calculate the current block time
	blockInfo, _ := ses.db.QueryLatestBlocks(2)

	lastBlocktime := blockInfo[0].Time.UTC()
	secondLastBlocktime := blockInfo[1].Time.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime)

	status, _ := ses.rpcClient.Status()
	block, _ := ses.rpcClient.Block(&status.SyncInfo.LatestBlockHeight)

	networkStats := &schema.StatsNetwork1H{
		BlockTime:       blockTime.Seconds(),
		TotalSupply:     totalSupplyTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     block.Block.TotalTxs,
		Time:            time.Now(),
	}

	result, _ := ses.db.InsertNetworkStats1H(*networkStats)
	if result {
		log.Println("succesfully saved NetworkStats 1H")
	}
}

// SaveNetworkStats24H saves network statistics 24 hours
func (ses *StatsExporterService) SaveNetworkStats24H() {
	// requests the current state of the staking pool
	poolResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/pool")

	var pool types.Pool
	err := json.Unmarshal(types.ReadRespWithHeight(poolResp).Result, &pool)
	if err != nil {
		fmt.Printf("failed to unmarshal Pool: %v \n ", err)
	}

	// requests the current minting inflation value
	inflationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/minting/inflation")

	var inflation types.Inflation
	err = json.Unmarshal(inflationResp.Body(), &inflation)
	if err != nil {
		fmt.Printf("failed to unmarshal Inflation: %v \n ", err)
	}

	// requests the current total supply
	totalSupplyResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/supply/total")

	var coin []sdk.Coin
	err = json.Unmarshal(types.ReadRespWithHeight(totalSupplyResp).Result, &coin)
	if err != nil {
		fmt.Printf("failed to unmarshal Coin: %v \n", err)
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalSupplyTokens, _ := strconv.ParseFloat(coin[0].Amount.String(), 64)
	bondedRatio := bondedTokens / totalSupplyTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation.Result, 64)

	// query two latest block numbers to calculate the current block time
	blockInfo, _ := ses.db.QueryLatestBlocks(2)

	lastBlocktime := blockInfo[0].Time.UTC()
	secondLastBlocktime := blockInfo[1].Time.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime)

	status, _ := ses.rpcClient.Status()
	block, _ := ses.rpcClient.Block(&status.SyncInfo.LatestBlockHeight)

	networkStats := &schema.StatsNetwork24H{
		BlockTime:       blockTime.Seconds(),
		TotalSupply:     totalSupplyTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     block.Block.TotalTxs,
		Time:            time.Now(),
	}

	result, _ := ses.db.InsertNetworkStats24H(*networkStats)
	if result {
		log.Println("succesfully saved NetworkStats 24H")
	}
}
