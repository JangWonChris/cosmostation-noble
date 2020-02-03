package db

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"
)

// Query param that is used when querying transactions that are made by an account
const (
	QueryTxsParamFromAddress      = "messages->0->'value'->>'from_address' = "
	QueryTxsParamToAddress        = "messages->0->'value'->>'to_address' = "
	QueryTxsParamInputsAddress    = "messages->0->'value'->'inputs'->0->>'address' = "
	QueryTxsParamOutpusAddress    = "messages->0->'value'->'outputs'->0->>'address' = "
	QueryTxsParamDelegatorAddress = "messages->0->'value'->>'delegator_address' = "
	QueryTxsParamAddress          = "messages->0->'value'->>'address' = "
	QueryTxsParamProposer         = "messages->0->'value'->>'proposer' = "
	QueryTxsParamDepositer        = "messages->0->'value'->>'depositor' = "
	QueryTxsParamVoter            = "messages->0->'value'->>'voter' = "
	QueryTxsParamValidatorAddress = "messages->0->>'type' = 'cosmos-sdk/MsgWithdrawValidatorCommission' AND messages->0->'value'->>'validator_address' = "
)

// QueryTxsByAddr queries transactions that are made by an account
func (db *Database) QueryTxsByAddr(address string, operAddr string, limit int, offset int, before int, after int) ([]schema.TransactionInfo, error) {
	var txs []schema.TransactionInfo

	// Make sure to use brackets that surround each local operator.
	// Otherwise it returns incorrect data
	params := "(" + QueryTxsParamFromAddress + "'" + address + "'" + " OR " +
		QueryTxsParamToAddress + "'" + address + "'" + " OR " +
		QueryTxsParamInputsAddress + "'" + address + "'" + " OR " +
		QueryTxsParamOutpusAddress + "'" + address + "'" + " OR " +
		QueryTxsParamDelegatorAddress + "'" + address + "'" + " OR " +
		QueryTxsParamAddress + "'" + address + "'" + " OR " +
		QueryTxsParamProposer + "'" + address + "'" + " OR " +
		QueryTxsParamDepositer + "'" + address + "'" + " OR " +
		QueryTxsParamVoter + "'" + address + "'" + " OR " +
		QueryTxsParamValidatorAddress + "'" + operAddr + "'" + ")"

	switch {
	case before > 0:
		params += " AND (height < ?)"
		_ = db.Model(&txs).
			Where(params, before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after >= 0:
		params += " AND (height > ?)"
		_ = db.Model(&txs).
			Where(params, after).
			Limit(limit).
			Order("id ASC").
			Select()
	case offset >= 0:
		_ = db.Model(&txs).
			Where(params).
			Limit(limit).
			Offset(offset).
			Order("id DESC").
			Select()
	}

	return txs, nil
}

// QueryTransactions queries transactions
func (db *Database) QueryTransactions(height int64) ([]schema.TransactionInfo, error) {
	var txInfos []schema.TransactionInfo
	_ = db.Model(&txInfos).
		Column("tx_hash").
		Where("height = ?", height).
		Select()

	return txInfos, nil
}

// QueryTotalTxsNum queries total number of transactions
func (db *Database) QueryTotalTxsNum() int64 {
	var blockInfo schema.BlockInfo
	_ = db.Model(&blockInfo).
		Column("total_txs").
		Order("height DESC").
		Limit(1).
		Select()

	return blockInfo.TotalTxs
}
