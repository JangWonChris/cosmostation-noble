package db

import (
	"fmt"

	mdschema "github.com/cosmostation/mintscan-database/schema"
	pg "github.com/go-pg/pg/v10"
)

var (
	limit                                                     = 50
	transferMsgExp                                            string   // QueryTransferTransactionsByAddr에서 사용하는 msgid 집합
	stakingMsgExp                                             string   // QueryTransactionsBetweenAccountAndValidator 에서 사용하는 msgid 집합
	AccountHistoryStmtWithoutTxID, AccountHistoryStmtWithTxID *pg.Stmt // account tx history prepared statement
	TransferStmtWithoutTxID, TransferStmtWithTxID             *pg.Stmt // transfer tx history prepared statement
	VDTXStmtWithoutTxID, VDTXStmtWithTxID                     *pg.Stmt // validator - delegator tx history prepared statement

)

func PrepareTransferMsgExp(msgs ...int) {
	// 미리 준비해놓고 계속 가져다 써야 함 : transferMsgType
	for i := range msgs {
		transferMsgExp += fmt.Sprintf("%d", msgs[i])
		if i != len(msgs)-1 {
			transferMsgExp += ","
		}
	}
	fmt.Println("transfer Msg set :", transferMsgExp)
}

func PrepareStakingMsgExp(msgs ...int) {
	// 미리 준비해놓고 계속 가져다 써야 함 : transferMsgType
	for i := range msgs {
		stakingMsgExp += fmt.Sprintf("%d", msgs[i])
		if i != len(msgs)-1 {
			stakingMsgExp += ","
		}
	}
	fmt.Println("staking Msg set :", stakingMsgExp)
}
func (db *Database) PrepareStmt() {
	var err error
	AccountHistoryStmtWithoutTxID, err = db.Prepare("select t.id, t.chain_info_id, t.block_id, t.chunk, t.timestamp from " + mdschema.GetCommonSchema() + ".transaction t, ( " +
		"select tx_id from " + mdschema.GetCommonSchema() + ".transaction_account ta " +
		"join " + mdschema.GetCommonSchema() + ".account a on a.id = ta.account_id " +
		"where a.address = $1 " +
		"order by ta.tx_id desc " +
		"limit $2 " + // limit
		") sub " +
		"where t.id = sub.tx_id order by t.id desc")
	if err != nil {
		panic(err)
	}
	AccountHistoryStmtWithTxID, err = db.Prepare("select t.id, t.chain_info_id, t.block_id, t.chunk, t.timestamp from " + mdschema.GetCommonSchema() + ".transaction t, ( " +
		"select tx_id from " + mdschema.GetCommonSchema() + ".transaction_account ta " +
		"join " + mdschema.GetCommonSchema() + ".account a on a.id = ta.account_id " +
		"where ta.tx_id < $1 and a.address = $2 " + //txid, address
		"order by ta.tx_id desc " +
		"limit $3 " + // limit
		") sub " +
		"where t.id = sub.tx_id order by t.id desc")
	if err != nil {
		panic(err)
	}

	TransferStmtWithoutTxID, err = db.Prepare("select t.id, t.chain_info_id, t.block_id, t.chunk, t.timestamp from " + mdschema.GetCommonSchema() + ".transaction t, ( " +
		"select tx_id from " + mdschema.GetCommonSchema() + ".transaction_message tm " +
		"join " + mdschema.GetCommonSchema() + ".account a on a.id = tm.account_id " +
		"where a.address = $1 " + //address
		"and tm.msg_id in (" + transferMsgExp + ") " +
		"order by tx_id desc " +
		"limit $2 " + // limit
		") sub " +
		"where t.id = sub.tx_id " +
		"order by t.id desc")
	if err != nil {
		panic(err)
	}
	TransferStmtWithTxID, err = db.Prepare("select t.id, t.chain_info_id, t.block_id, t.chunk, t.timestamp from " + mdschema.GetCommonSchema() + ".transaction t, ( " +
		"select tx_id from " + mdschema.GetCommonSchema() + ".transaction_message tm " +
		"join " + mdschema.GetCommonSchema() + ".account a on a.id = tm.account_id " +
		"where tm.tx_id < $1 " + // tx_id
		"and a.address = $2 " + //address
		"and tm.msg_id in (" + transferMsgExp + ") " +
		"order by tx_id desc " +
		"limit $3 " + // limit
		") sub " +
		"where t.id = sub.tx_id " +
		"order by t.id desc")
	if err != nil {
		panic(err)
	}

	VDTXStmtWithoutTxID, err = db.Prepare("select t.id, t.chain_info_id, t.block_id, t.chunk, t.timestamp from " + mdschema.GetCommonSchema() + ".transaction t, " +
		"(select tm.tx_id from " + mdschema.GetCommonSchema() + ".transaction_message tm, " +
		"( select tx_id from " + mdschema.GetCommonSchema() + ".transaction_message " +
		"where account_id = (select id from " + mdschema.GetCommonSchema() + ".account where address = $1) " + //validator address
		"and msg_id in (" + stakingMsgExp + ") " +
		"order by tx_id desc " + //limit 걸면 안됨
		") tmsub " +
		"where tm.tx_id = tmsub.tx_id " +
		"and account_id = (select id from " + mdschema.GetCommonSchema() + ".account where address = $2) " + // delegator address
		"order by tm.tx_id desc limit $3 " + //limit
		") sub " +
		"where t.id = sub.tx_id " +
		"order by t.id desc ")
	if err != nil {
		panic(err)
	}

	VDTXStmtWithTxID, err = db.Prepare("select t.id, t.chain_info_id, t.block_id, t.chunk, t.timestamp from " + mdschema.GetCommonSchema() + ".transaction t, " +
		"(select tm.tx_id from " + mdschema.GetCommonSchema() + ".transaction_message tm, " +
		"( select tx_id from " + mdschema.GetCommonSchema() + ".transaction_message " +
		"where tx_id < $1 " + // tx_id
		"and account_id = (select id from " + mdschema.GetCommonSchema() + ".account where address = $2) " + //validator address
		"and msg_id in (" + stakingMsgExp + ") " +
		"order by tx_id desc " + //limit 걸면 안됨
		") tmsub " +
		"where tm.tx_id < $3 " + // tx_id
		"and tm.tx_id = tmsub.tx_id " +
		"and account_id = (select id from " + mdschema.GetCommonSchema() + ".account where address = $4) " + // delegator address
		"order by tm.tx_id desc limit $5 " + //limit
		") sub " +
		"where t.id < $6 " + // tx_id
		"and t.id = sub.tx_id " +
		"order by t.id desc ")
	if err != nil {
		panic(err)
	}
}

// QueryTransactionsByAddr returns all transactions that are created by an account.
func (db *Database) QueryTransactionsByAddr(txID int64, accAddr string) ([]mdschema.Transaction, error) {
	var txs []mdschema.Transaction
	var err error

	switch {
	case txID < 0:
		return txs, nil
	case txID == 0:
		_, err = AccountHistoryStmtWithoutTxID.Query(&txs, accAddr, limit)
	case txID > 0:
		_, err = AccountHistoryStmtWithTxID.Query(&txs, txID, accAddr, limit)
	}

	if err != nil {
		if err == pg.ErrNoRows {
			return txs, nil
		}
		return txs, err
	}

	return txs, nil
}

// QueryTransferTransactionsByAddr queries Send / MultiSend transactions that are made by an account
func (db *Database) QueryTransferTransactionsByAddr(txID int64, accAddr string) ([]mdschema.Transaction, error) {
	var txs []mdschema.Transaction
	var err error

	switch {
	case txID < 0:
		return txs, nil
	case txID == 0:
		_, err = TransferStmtWithoutTxID.Query(&txs, accAddr, limit)
	case txID > 0:
		_, err = TransferStmtWithTxID.Query(&txs, txID, accAddr, limit)
	}

	if err != nil {
		if err == pg.ErrNoRows {
			return txs, nil
		}
		return txs, err
	}

	return txs, nil
}

// QueryTransactionsBetweenAccountAndValidator queries transactions that are made between an account and his delegated validator
func (db *Database) QueryTransactionsBetweenAccountAndValidator(txID int64, accAddr, valAddr string) ([]mdschema.Transaction, error) {
	var txs []mdschema.Transaction
	var err error

	switch {
	case txID < 0:
		return txs, nil
	case txID == 0:
		_, err = VDTXStmtWithoutTxID.Query(&txs, valAddr, accAddr, limit)
	case txID > 0:
		_, err = VDTXStmtWithTxID.Query(&txs, txID, valAddr, txID, accAddr, limit, txID)
	}

	if err != nil {
		if err == pg.ErrNoRows {
			return txs, nil
		}
		return txs, err
	}

	return txs, nil
}
