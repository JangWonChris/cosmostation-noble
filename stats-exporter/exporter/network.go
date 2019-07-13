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

// SaveNetworkStats
func (ses *StatsExporterService) SaveNetworkStats() {
	log.Println("Network Stats")

	// Query LCD
	resp, err := resty.R().Get(ses.config.Node.LCDURL + "/staking/pool")
	if err != nil {
		fmt.Printf("Staking Pool LCD resty - %v\n", err)
	}

	// Parse Proposal struct
	var pool types.Pool
	err = json.Unmarshal(resp.Body(), &pool)
	if err != nil {
		fmt.Printf("Proposal unmarshal error - %v\n", err)
	}

	bondedTokens, _ := strconv.ParseInt(pool.BondedTokens, 10, 64)
	notBondedTokens, _ := strconv.ParseInt(pool.NotBondedTokens, 10, 64)
	totalBondedTokens := bondedTokens + notBondedTokens
	bondedRatio := (float64(bondedTokens) / float64(totalBondedTokens)) * 100

	// Block Time
	var blockInfo []types.BlockInfo
	err = ses.db.Model(&blockInfo).
		Column("time").
		Order("height DESC").
		Limit(2).
		Select()
	if err != nil {
		fmt.Printf("BlockInfo DB error - %v\n", err)
	}

	// Latest block time and its previous block time
	lastBlocktime := blockInfo[0].Time.UTC()
	secondLastBlocktime := blockInfo[1].Time.UTC()
	blockTime := lastBlocktime.Sub(secondLastBlocktime)

	// NetworkStats
	networkStats := &types.NetworkStats{
		BlockTime1H:       blockTime.Seconds(),
		BondedTokens1H:    bondedTokens,
		BondedRatio1H:     bondedRatio,
		NotBondedTokens1H: notBondedTokens,
		LastUpdated:       time.Now(),
	}

	// Save
	_, err = ses.db.Model(networkStats).Insert()
	if err != nil {
		fmt.Printf("Save NetworkStats error - %v\n", err)
	}
}
