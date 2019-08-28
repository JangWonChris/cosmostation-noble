package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"

	resty "gopkg.in/resty.v1"
)

// SaveNetworkStats1H
func (ses *StatsExporterService) SaveNetworkStats1H() {
	log.Println("Network Stats 1H")

	// query pool
	var pool types.Pool
	resp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/pool")
	if err != nil {
		fmt.Printf("Query /staking/pool error - %v\n ", err)
	}

	err = json.Unmarshal(resp.Body(), &pool)
	if err != nil {
		fmt.Printf("Unmarshal pool error - %v\n ", err)
	}

	// query inflation rate
	var inflation types.Inflation
	inflationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/minting/inflation")
	err = json.Unmarshal(inflationResp.Body(), &inflation)
	if err != nil {
		fmt.Printf("Unmarshal inflation error - %v\n ", err)
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalBondedTokens := bondedTokens + notBondedTokens
	bondedRatio := bondedTokens / totalBondedTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation.Result, 64)

	// get block time - (last block time - second last block time)
	var blockInfo []types.BlockInfo
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
		TotalSupply:     totalBondedTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     block.Block.TotalTxs,
		Time:            time.Now(),
	}

	// Save
	_, err = ses.db.Model(networkStats).Insert()
	if err != nil {
		fmt.Printf("save networkStats error - %v\n ", err)
	}
}

// SaveNetworkStats24H
func (ses *StatsExporterService) SaveNetworkStats24H() {
	log.Println("Network Stats 1H")

	// query pool
	var pool types.Pool
	resp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/pool")
	if err != nil {
		fmt.Printf("Query /staking/pool error - %v\n ", err)
	}

	err = json.Unmarshal(resp.Body(), &pool)
	if err != nil {
		fmt.Printf("Unmarshal pool error - %v\n ", err)
	}

	// query inflation rate
	var inflation string
	inflationResp, _ := resty.R().Get(ses.config.Node.LCDURL + "/minting/inflation")
	err = json.Unmarshal(inflationResp.Body(), &inflation)
	if err != nil {
		fmt.Printf("Unmarshal inflation error - %v\n ", err)
	}

	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	totalBondedTokens := bondedTokens + notBondedTokens
	bondedRatio := bondedTokens / totalBondedTokens * 100
	inflationRatio, _ := strconv.ParseFloat(inflation, 64)

	// get block time - (last block time - second last block time)
	var blockInfo []types.BlockInfo
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
		TotalSupply:     totalBondedTokens,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		BondedRatio:     bondedRatio,
		InflationRatio:  inflationRatio,
		TotalTxsNum:     block.Block.TotalTxs,
		Time:            time.Now(),
	}

	// Save
	_, err = ses.db.Model(networkStats).Insert()
	if err != nil {
		fmt.Printf("save networkStats error - %v\n ", err)
	}
}
