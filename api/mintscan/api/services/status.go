package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	dbtypes "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	u "github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/go-pg/pg"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetStatus returns ResultStatus, which includes current network status
func GetStatus(RPCClient *client.HTTP, DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Query LCD - stake pool to get bonded and unbonded tokens
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(Config.Node.LCDURL + "/staking/pool")

	// Unmarshal Pool struct
	var pool *models.Pool
	err := json.Unmarshal(resp.Body(), &pool)
	if err != nil {
		fmt.Printf("staking/pool unmarshal pool error - %v\n", err)
	}

	// Convert tokens from string to float64 type
	notBondedTokens, err := strconv.ParseFloat(pool.NotBondedTokens, 64)
	bondedTokens, err := strconv.ParseFloat(pool.BondedTokens, 64)

	// Get a number of unjailed validators
	var unjailedValidators dbtypes.ValidatorInfo
	unJailedNum, _ := DB.Model(&unjailedValidators).
		Where("jailed = ?", false).
		Count()

	// Get a number of jailed validators
	var jailedValidators dbtypes.ValidatorInfo
	jailedNum, _ := DB.Model(&jailedValidators).
		Where("jailed = ?", true).
		Count()

	// Total Txs Num
	var blockInfo dbtypes.BlockInfo
	_ = DB.Model(&blockInfo).
		Column("total_txs").
		Order("height DESC").
		Limit(1).
		Select()

	// Query status
	status, _ := RPCClient.Status()

	// Query the lastly saved block time
	var lastBlockTime []dbtypes.BlockInfo
	_ = DB.Model(&lastBlockTime).
		Column("time").
		Order("height DESC").
		Limit(2).
		Select()

	// Latest block time and its previous block time
	lastBlocktime := lastBlockTime[0].Time.UTC()
	secondLastBlocktime := lastBlockTime[1].Time.UTC()

	// * 실질적으로 status.SyncInfo.LatestBlockTime.UTC()로 비교를 해야 되지만 현재로써는 마지막, 두번째마지막으로 비교
	// Get the block time that is taken from the previous block
	diff := lastBlocktime.Sub(secondLastBlocktime)

	resultStatus := &models.ResultStatus{
		ChainID:                status.NodeInfo.Network,
		BlockHeight:            status.SyncInfo.LatestBlockHeight,
		BlockTime:              diff.Seconds(),
		TotalTxsNum:            blockInfo.TotalTxs,
		TotalValidatorNum:      unJailedNum + jailedNum,
		UnjailedValidatorNum:   unJailedNum,
		JailedValidatorNum:     jailedNum,
		TotalCirculatingTokens: bondedTokens + notBondedTokens,
		BondedTokens:           bondedTokens,
		NotBondedTokens:        notBondedTokens,
		Time:                   status.SyncInfo.LatestBlockTime,
	}

	u.Respond(w, resultStatus)
	return nil
}
