package db

import "github.com/cosmostation/cosmostation-cosmos/mintscan/api/schema"

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
