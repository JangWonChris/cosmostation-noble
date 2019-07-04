package services

import (
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

/*
	아래 두 기관에 분배된 토큰 수량 - 236,198,958.12

		All in Bits Inc
		Interchain Foundation

	Total Supply (not_bonded_tokens + bonded_tokens)
		- 1,777,707
		- 21,842,188.81
		+ 992826.76 * 2 (지금까지 총 두번) - 언제 정확히 풀리는지는 알아야 매월마다 합산된 량을 + 해준다.

	vesting end time (UNIX Epoch time) - https://www.epochconverter.com/
	1584140400 - GMT: 2020년 March 13일 Friday PM 11:00:00 (About 8 months left till now) - 45 accounts

	Endtime
*/

var (
	FoundationLockedTokens = 1777707
)

// GetStatus returns ResultStatus, which includes current network status
func GetStatus(RPCClient *client.HTTP, DB *pg.DB, Config *config.Config, w http.ResponseWriter, r *http.Request) error {
	// Query LCD - stake pool to get bonded and unbonded tokens
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

	totalSupplyTokens := bondedTokens + notBondedTokens
	// totalCirculatingTokens := 2

	resultStatus := &models.ResultStatus{
		ChainID:                status.NodeInfo.Network,
		BlockHeight:            status.SyncInfo.LatestBlockHeight,
		BlockTime:              diff.Seconds(),
		TotalTxsNum:            blockInfo.TotalTxs,
		TotalValidatorNum:      unJailedNum + jailedNum,
		TotalSupplyTokens:      totalSupplyTokens,
		TotalCirculatingTokens: bondedTokens + notBondedTokens,
		JailedValidatorNum:     jailedNum,
		UnjailedValidatorNum:   unJailedNum,
		BondedTokens:           bondedTokens,
		NotBondedTokens:        notBondedTokens,
		Time:                   status.SyncInfo.LatestBlockTime,
	}

	u.Respond(w, resultStatus)
	return nil
}
