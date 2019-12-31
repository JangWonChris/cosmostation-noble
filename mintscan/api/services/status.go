package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

/*
	[Vesting]

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

// GetStatus returns ResultStatus, which includes current network status
func GetStatus(config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	resp, _ := resty.R().Get(config.Node.LCDURL + "/staking/pool")

	var pool models.Pool
	err := json.Unmarshal(types.ReadRespWithHeight(resp).Result, &pool)
	if err != nil {
		fmt.Printf("staking/pool unmarshal pool error - %v\n", err)
		log.Info().Str(models.Service, models.LogStatus).Str(models.Method, "GetStatus").Err(err).Msg("unmarshal pool error")
	}

	// Query total supply
	totalSupplyResp, _ := resty.R().Get(config.Node.LCDURL + "/supply/total")

	var coin []models.Coin
	err = json.Unmarshal(types.ReadRespWithHeight(totalSupplyResp).Result, &coin)
	if err != nil {
		fmt.Printf("supply/total unmarshal supply total error - %v\n", err)
	}

	notBondedTokens, _ := strconv.ParseFloat(pool.NotBondedTokens, 64)
	bondedTokens, _ := strconv.ParseFloat(pool.BondedTokens, 64)
	totalSupplyTokens, _ := strconv.ParseFloat(coin[0].Amount, 64)

	// a number of unjailed validators
	var unjailedValidators types.ValidatorInfo
	unJailedNum, _ := db.Model(&unjailedValidators).
		Where("status = ?", 2).
		Count()

	// a number of jailed validators
	var jailedValidators types.ValidatorInfo
	jailedNum, _ := db.Model(&jailedValidators).
		Where("status = ? OR status = ?", 0, 1).
		Count()

	// total txs num
	var blockInfo types.BlockInfo
	_ = db.Model(&blockInfo).
		Column("total_txs").
		Order("height DESC").
		Limit(1).
		Select()

	// query status
	status, _ := rpcClient.Status()

	// query the lastly saved block time
	var lastBlockTime []types.BlockInfo
	_ = db.Model(&lastBlockTime).
		Column("time").
		Order("height DESC").
		Limit(2).
		Select()

	// latest block time and its previous block time
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
		TotalSupplyTokens:      totalSupplyTokens,
		TotalCirculatingTokens: bondedTokens + notBondedTokens,
		JailedValidatorNum:     jailedNum,
		UnjailedValidatorNum:   unJailedNum,
		BondedTokens:           bondedTokens,
		NotBondedTokens:        notBondedTokens,
		Time:                   status.SyncInfo.LatestBlockTime,
	}

	utils.Respond(w, resultStatus)
	return nil
}
