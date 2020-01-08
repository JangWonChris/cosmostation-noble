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
	log.Println("Network Stats 1H")

	// query pool
	poolResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/pool")

	var pool types.Pool
	err := json.Unmarshal(types.ReadRespWithHeight(poolResp).Result, &pool)
	if err != nil {
		fmt.Printf("unmarshal pool error - %v\n ", err)
	}

	// query inflation rate
	inflationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/minting/inflation")

	var inflation types.Inflation
	err = json.Unmarshal(inflationResp.Body(), &inflation)
	if err != nil {
		fmt.Printf("unmarshal inflation error - %v\n ", err)
	}

	// Query total supply
	totalSupplyResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/supply/total")

	var coin []sdk.Coin
	err = json.Unmarshal(types.ReadRespWithHeight(totalSupplyResp).Result, &coin)
	if err != nil {
		fmt.Printf("supply/total unmarshal supply total error - %v\n", err)
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalSupplyTokens, _ := strconv.ParseFloat(coin[0].Amount.String(), 64)
	bondedRatio := bondedTokens / totalSupplyTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation.Result, 64)

	// get block time - (last block time - second last block time)
	var blockInfo []schema.BlockInfo
	err = ses.db.Model(&blockInfo).
		Column("time").
		Order("height DESC").
		Limit(2).
		Select()
	if err != nil {
		fmt.Printf("blockInfo database error - %v\n ", err)
	}

	// txs num
	status, _ := ses.rpcClient.Status()
	block, _ := ses.rpcClient.Block(&status.SyncInfo.LatestBlockHeight)

	lastBlocktime := blockInfo[0].Time.UTC()
	secondLastBlocktime := blockInfo[1].Time.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime)

	networkStats := &types.StatsNetwork1H{
		BlockTime:       blockTime.Seconds(),
		TotalSupply:     totalSupplyTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     block.Block.TotalTxs,
		Time:            time.Now(),
	}

	_, err = ses.db.Model(networkStats).Insert()
	if err != nil {
		fmt.Printf("save networkStats1H error - %v\n ", err)
	}
}

// SaveNetworkStats24H saves network statistics 24 hours
func (ses *StatsExporterService) SaveNetworkStats24H() {
	log.Println("Network Stats 24H")

	// query pool
	poolResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/staking/pool")

	var pool types.Pool
	err := json.Unmarshal(types.ReadRespWithHeight(poolResp).Result, &pool)
	if err != nil {
		fmt.Printf("unmarshal pool error - %v\n ", err)
	}

	// query inflation rate
	inflationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/minting/inflation")

	var inflation types.Inflation
	err = json.Unmarshal(inflationResp.Body(), &inflation)
	if err != nil {
		fmt.Printf("unmarshal inflation error - %v\n ", err)
	}

	// Query total supply
	totalSupplyResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/supply/total")

	var coin []sdk.Coin
	err = json.Unmarshal(types.ReadRespWithHeight(totalSupplyResp).Result, &coin)
	if err != nil {
		fmt.Printf("supply/total unmarshal supply total error - %v\n", err)
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalSupplyTokens, _ := strconv.ParseFloat(coin[0].Amount.String(), 64)
	bondedRatio := bondedTokens / totalSupplyTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation.Result, 64)

	// get block time - (last block time - second last block time)
	var blockInfo []schema.BlockInfo
	err = ses.db.Model(&blockInfo).
		Column("time").
		Order("height DESC").
		Limit(2).
		Select()
	if err != nil {
		fmt.Printf("blockInfo database error - %v\n ", err)
	}

	// txs num
	status, _ := ses.rpcClient.Status()
	block, _ := ses.rpcClient.Block(&status.SyncInfo.LatestBlockHeight)

	lastBlocktime := blockInfo[0].Time.UTC()
	secondLastBlocktime := blockInfo[1].Time.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime)

	networkStats := &types.StatsNetwork24H{
		BlockTime:       blockTime.Seconds(),
		TotalSupply:     totalSupplyTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     block.Block.TotalTxs,
		Time:            time.Now(),
	}

	_, err = ses.db.Model(networkStats).Insert()
	if err != nil {
		fmt.Printf("save networkStats24H error - %v\n ", err)
	}
}
