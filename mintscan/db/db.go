package db

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/model"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*pg.DB
}

// Connect opens a database connections with the given database connection info from config.
func Connect(cfg config.DBConfig) *Database {
	db := pg.Connect(&pg.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		User:         cfg.User,
		Password:     cfg.Password,
		Database:     cfg.Table,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Disable pluralization
	orm.SetTableNameInflector(func(s string) string {
		return s
	})

	return &Database{db}
}

// Ping returns a database connection handle or an error if the connection fails.
func (db *Database) Ping() error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

// --------------------
// Query
// --------------------

// QueryLatestBlockHeight queries the latest block height in database.
// return 0 if there is not row in result set and -1 for any type of database errors.
func (db *Database) QueryLatestBlockHeight() (int64, error) {
	var block schema.Block
	err := db.Model(&block).
		Order("height DESC").
		Limit(1).
		Select()

	if err == pg.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return -1, err
	}

	return block.Height, nil
}

// QueryBlocks queries blocks information with given parameters.
func (db *Database) QueryBlocks(before, after int, limit int) (blocks []schema.Block, err error) {
	switch {
	case before > 0:
		err = db.Model(&blocks).
			Where("height < ?", before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after >= 0:
		err = db.Model(&blocks).
			Where("height > ?", after).
			Limit(limit).
			Order("id ASC").
			Select()
	default:
		err = db.Model(&blocks).
			Limit(limit).
			Order("id DESC").
			Select()
	}

	if err != nil {
		return []schema.Block{}, err
	}

	return blocks, nil
}

// QueryBlocksByProposer queries blocks by proposer
func (db *Database) QueryBlocksByProposer(address string, before, after, limit int) (blocks []schema.Block, err error) {
	switch {
	case before > 0:
		err = db.Model(&blocks).
			Where("proposer = ? AND height < ?", address, before).
			Limit(limit).
			Order("height DESC").
			Select()
	case after > 0:
		err = db.Model(&blocks).
			Where("proposer = ? AND height > ?", address, after).
			Limit(limit).
			Order("height ASC").
			Select()
	default:
		err = db.Model(&blocks).
			Where("proposer = ?", address).
			Limit(limit).
			Order("id DESC").
			Select()
	}

	if err != nil {
		return []schema.Block{}, err
	}

	return blocks, nil
}

// QueryLastestTwoBlocks queries lastest two blocks for blocktime calculation.
func (db *Database) QueryLastestTwoBlocks() (blocks []schema.Block, err error) {
	err = db.Model(&blocks).
		Order("height DESC").
		Limit(2).
		Select()

	if err != nil {
		return []schema.Block{}, err
	}

	return blocks, nil
}

// QueryMissingBlocksDetail queries how many missing blocks a validator misses in detail.
func (db *Database) QueryMissingBlocksDetail(address string, latestHeight int64, count int) (misses []schema.MissDetail, err error) {
	err = db.Model(&misses).
		Where("address = ? AND height BETWEEN ? AND ?", address, int(latestHeight)-count, latestHeight).
		Limit(count).
		Order("height DESC").
		Select()

	if err != nil {
		return []schema.MissDetail{}, err
	}

	return misses, nil
}

// QueryMissingBlocks queries a range of missing blocks a validator misses.
func (db *Database) QueryMissingBlocks(address string, limit int) (misses []schema.Miss, err error) {
	err = db.Model(&misses).
		Where("address = ?", address).
		Limit(limit).
		Order("start_height DESC").
		Select()

	if err != nil {
		return []schema.Miss{}, err
	}

	return misses, nil
}

// QueryProposals returns proposals.
func (db *Database) QueryProposals() (proposals []schema.Proposal, err error) {
	err = db.Model(&proposals).Select()
	if err != nil {
		return []schema.Proposal{}, err
	}

	return proposals, nil
}

// QueryProposal returns a proposal.
func (db *Database) QueryProposal(id string) (proposal schema.Proposal, err error) {
	err = db.Model(&proposal).
		Where("id = ?", id).
		Select()

	if err != nil {
		return schema.Proposal{}, err
	}

	return proposal, nil
}

// QueryDeposits returns all deposit information.
func (db *Database) QueryDeposits(id string) (deposits []schema.Deposit, err error) {
	err = db.Model(&deposits).
		Where("proposal_id = ?", id).
		Order("id DESC").
		Select()

	if err != nil {
		return []schema.Deposit{}, err
	}

	return deposits, nil
}

// QueryVotes returns all vote information.
func (db *Database) QueryVotes(id string) (votes []schema.Vote, err error) {
	err = db.Model(&votes).
		Where("proposal_id = ?", id).
		Order("id DESC").
		Select()

	if err != nil {
		return []schema.Vote{}, err
	}

	return votes, nil
}

// QueryVoteOptions queries all vote options for the proposal
func (db *Database) QueryVoteOptions(id string) (int, int, int, int) {
	votes := make([]schema.Vote, 0)

	yes, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, model.YES).
		Count()

	no, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, model.NO).
		Count()

	noWithVeto, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, model.NOWITHVETO).
		Count()

	abstain, _ := db.Model(&votes).
		Where("proposal_id = ? AND option = ?", id, model.ABSTAIN).
		Count()

	return yes, no, noWithVeto, abstain
}

// QueryValidators returns all validators.
func (db *Database) QueryValidators() (validators []*schema.Validator, err error) {
	err = db.Model(&validators).
		Order("id ASC").
		Select()

	if err != nil {
		return []*schema.Validator{}, err
	}

	return validators, nil
}

// QueryValidatorByID returns a validator by querying with validator id.
// Validator id is determined by their voting power when chain exporter aggregates validator power event data.
func (db *Database) QueryValidatorByID(address string) (int, error) {
	var peh schema.PowerEventHistory
	err := db.Model(&peh).
		Column("id_validator").
		Where("proposer = ?", address).
		Limit(1).
		Select()

	// return 0 when there is no row in result set
	if err == pg.ErrNoRows {
		return 0, err
	}

	// return -1 for any type of errors
	if err != nil {
		return -1, err
	}

	return peh.IDValidator, nil
}

// QueryValidatorByValAddr returns a validator by querying with validator operator address.
func (db *Database) QueryValidatorByValAddr(valAddr string) (validator schema.Validator, err error) {
	err = db.Model(&validator).
		Where("operator_address = ?", valAddr).
		Limit(1).
		Select()

	if err == pg.ErrNoRows {
		return schema.Validator{}, fmt.Errorf("no rows in block table: %s", err)
	}

	if err != nil {
		return schema.Validator{}, err
	}

	return validator, nil
}

// QueryValidatorsByStatus returns a validator by querying with bonding status.
func (db *Database) QueryValidatorsByStatus(status int) (validators []*schema.Validator, err error) {
	err = db.Model(&validators).
		Where("status = ?", status).
		Order("id ASC").
		Select()

	if err == pg.ErrNoRows {
		return []*schema.Validator{}, nil
	}

	if err != nil {
		return []*schema.Validator{}, err
	}

	return validators, nil
}

// QueryValidatorBondedInfo returns a validator's bonded information.
func (db *Database) QueryValidatorBondedInfo(address string) (peh schema.PowerEventHistory, err error) {
	msgType := "create_validator"

	err = db.Model(&peh).
		Where("proposer = ? AND msg_type = ?", address, msgType).
		Limit(1).
		Select()

	if err != nil {
		return schema.PowerEventHistory{}, err
	}

	return peh, nil
}

// QueryValidatorVotingPowerEventHistory returns validator's power events with given parameters.
func (db *Database) QueryValidatorVotingPowerEventHistory(validatorID, before, after, limit int) (peh []schema.PowerEventHistory, err error) {
	switch {
	case before > 0:
		err = db.Model(&peh).
			Where("id_validator = ? AND height < ?", validatorID, before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after > 0:
		err = db.Model(&peh).
			Where("id_validator = ? AND height > ?", validatorID, after).
			Limit(limit).
			Order("id ASC").
			Select()
	default:
		err = db.Model(&peh).
			Where("id_validator = ?", validatorID).
			Limit(limit).
			Order("id DESC").
			Select()
	}

	if err == pg.ErrNoRows {
		return []schema.PowerEventHistory{}, nil
	}

	if err != nil {
		return []schema.PowerEventHistory{}, err
	}

	return peh, nil
}

// QueryValidatorByAny queries validator information by any type of input address.
func (db *Database) QueryValidatorByAny(address string) (val schema.Validator, err error) {
	switch {
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ValidatorPubPrefix()): // Bech32 prefix for validator public key
		err = db.Model(&val).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ValidatorAddrPrefix()): // Bech32 prefix for validator address
		err = db.Model(&val).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32AccountAddrPrefix()): // Bech32 prefix for account address
		err = db.Model(&val).
			Where("address = ?", address).
			Limit(1).
			Select()
	case len(address) == 40: // Validator consensus address in hex
		address := strings.ToUpper(address)
		err = db.Model(&val).
			Where("proposer = ?", address).
			Limit(1).
			Select()
	default:
		err = db.Model(&val).
			Where("moniker = ?", address). // Validator moniker
			Limit(1).
			Select()
	}

	if err == pg.ErrNoRows {
		return schema.Validator{}, nil
	}

	if err != nil {
		return schema.Validator{}, err
	}

	return val, nil
}

// QueryTransactions queries transactions with pagination params, such as limit, before, after, and offset
func (db *Database) QueryTransactions(before int, after int, limit int) (txs []schema.Transaction, err error) {
	switch {
	case before > 0:
		err = db.Model(&txs).
			Where("id < ?", before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after > 0:
		err = db.Model(&txs).
			Where("id > ?", after).
			Limit(limit).
			Order("id ASC").
			Select()
	default:
		err = db.Model(&txs).
			Limit(limit).
			Order("id DESC").
			Select()
	}

	if err == pg.ErrNoRows {
		return []schema.Transaction{}, nil
	}

	if err != nil {
		return []schema.Transaction{}, fmt.Errorf("unexpected database error: %s", err)
	}

	return txs, nil
}

// QueryTransactionByID returns transaction information with given id.
func (db *Database) QueryTransactionByID(id int64) (tx schema.Transaction, err error) {
	err = db.Model(&tx).
		Where("id = ?", id).
		Limit(1).
		Select()

	if err != nil {
		return schema.Transaction{}, err
	}

	return tx, nil
}

// QueryTransactionByTxHash returns transaction information with given tx hash.
func (db *Database) QueryTransactionByTxHash(txHashStr string) (tx schema.Transaction, err error) {
	err = db.Model(&tx).
		Where("tx_hash = ?", txHashStr).
		Limit(1).
		Select()

	if err != nil {
		return schema.Transaction{}, err
	}

	return tx, nil
}

// QueryTransactionsByBlockHeight returns transactions that are included in a single block.
func (db *Database) QueryTransactionsByBlockHeight(height int64) (txs []schema.Transaction, err error) {
	err = db.Model(&txs).
		Column("tx_hash").
		Where("height = ?", height).
		Select()

	if err != nil {
		return []schema.Transaction{}, err
	}

	return txs, nil
}

// QueryTransactionsByAddr returns all transactions that are created by an account.
func (db *Database) QueryTransactionsByAddr(accAddr, valAddr string, before, after, limit int) (txs []schema.Transaction, err error) {
	// Make sure to use brackets that surround each local operator, otherwise it will return incorrect data.
	params := "(" + QueryTxParamFromAddress + "'" + accAddr + "'" + " OR " +
		QueryTxParamToAddress + "'" + accAddr + "'" + " OR " +
		QueryTxParamInputsAddress + "'" + accAddr + "'" + " OR " +
		QueryTxParamOutpusAddress + "'" + accAddr + "'" + " OR " +
		QueryTxParamDelegatorAddress + "'" + accAddr + "'" + " OR " +
		QueryTxParamAddress + "'" + accAddr + "'" + " OR " +
		QueryTxParamProposer + "'" + accAddr + "'" + " OR " +
		QueryTxParamDepositer + "'" + accAddr + "'" + " OR " +
		QueryTxParamVoter + "'" + accAddr + "'" + " OR " +
		QueryTxParamValidatorCommission + " AND " + QueryTxParamValidatorAddress + "'" + valAddr + "'" + ")"

	switch {
	case before > 0:
		params += " AND (id < ?)"
		err = db.Model(&txs).
			Where("id < ?", before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after > 0:
		params += " AND (id > ?)"
		err = db.Model(&txs).
			Where("id > ?", after).
			Limit(limit).
			Order("id ASC").
			Select()
	default:
		err = db.Model(&txs).
			Where(params).
			Limit(limit).
			Order("id DESC").
			Select()
	}

	if err != nil {
		return []schema.Transaction{}, err
	}

	return txs, nil
}

// QueryTransferTransactionsByAddr queries Send / MultiSend transactions that are made by an account
func (db *Database) QueryTransferTransactionsByAddr(accAddr, denom string, before, after, limit int) (txs []schema.Transaction, err error) {
	params := "(" + QueryTxParamFromAddress + "'" + accAddr + "'" + " AND " + QueryTxParamDenom + "'" + denom + "')" + " OR " +
		"(" + QueryTxParamToAddress + "'" + accAddr + "'" + " AND " + QueryTxParamDenom + "'" + denom + "')" + " OR " +
		"(" + QueryTxParamInputsAddress + "'" + accAddr + "'" + " AND " + QueryTxParamDenom + "'" + denom + "')" + " OR " +
		"(" + QueryTxParamOutpusAddress + "'" + accAddr + "'" + " AND " + QueryTxParamDenom + "'" + denom + "')"

	switch {
	case before > 0:
		params += " AND (id < ?)"
		err = db.Model(&txs).
			Where(params, before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after > 0:
		params += " AND (id > ?)"
		err = db.Model(&txs).
			Where(params, after).
			Limit(limit).
			Order("id ASC").
			Select()
	default:
		err = db.Model(&txs).
			Where(params).
			Limit(limit).
			Order("id DESC").
			Select()
	}

	if err != nil {
		return []schema.Transaction{}, err
	}

	return txs, nil
}

// QueryTransactionsBetweenAccountAndValidator queries transactions that are made between an account and his delegated validator
func (db *Database) QueryTransactionsBetweenAccountAndValidator(address, valAddr string, before, after, limit int) (txs []schema.Transaction, err error) {
	params := "(" + QueryTxParamValidatorAddress + "'" + valAddr + "'" + " OR " +
		QueryTxParamValidatorDstAddress + "'" + valAddr + "'" + " OR " +
		QueryTxParamValidatorSrcAddress + "'" + valAddr + "')" + " AND " +
		"(" + QueryTxParamDelegatorAddress + "'" + address + "'" + ")"

	switch {
	case before > 0:
		params += " AND (id < ?)"
		err = db.Model(&txs).
			Where(params, before).
			Limit(limit).
			Order("id DESC").
			Select()
	case after > 0:
		params += " AND (id > ?)"
		err = db.Model(&txs).
			Where(params, after).
			Limit(limit).
			Order("id ASC").
			Select()
	default:
		err = db.Model(&txs).
			Where(params).
			Limit(limit).
			Order("id DESC").
			Select()
	}

	if err != nil {
		return []schema.Transaction{}, err
	}

	return txs, nil
}

// QueryTotalTransactionNum queries total number of transactions
func (db *Database) QueryTotalTransactionNum() int {
	var tx schema.Transaction
	_ = db.Model(&tx).
		Order("id DESC").
		Limit(1).
		Select()

	return int(tx.ID)
}

// QueryValidatorStats1D returns validator statistics from 1 day validator stats table.
func (db *Database) QueryValidatorStats1D(address string, limit int) ([]schema.StatsValidators1D, error) {
	statsValidators24H := make([]schema.StatsValidators1D, 0)
	_ = db.Model(&statsValidators24H).
		Where("proposer = ?", address).
		Order("id DESC").
		Limit(limit).
		Select()

	return statsValidators24H, nil
}

// QueryPriceFromMarketStat5M returns market data from 5 minutes makret stats table.
func (db *Database) QueryPriceFromMarketStat5M() (data schema.StatsMarket5M, err error) {
	err = db.Model(&data).
		Order("id DESC").
		Limit(1).
		Select()

	if err != nil {
		return schema.StatsMarket5M{}, err
	}

	return data, nil
}

// QueryPricesFromMarketStat1H returns market statistics from 1 hour makret stats table.
func (db *Database) QueryPricesFromMarketStat1H(limit int) (stats []schema.StatsMarket1H, err error) {
	err = db.Model(&stats).
		Order("id DESC").
		Limit(limit).
		Select()

	if err != nil {
		return []schema.StatsMarket1H{}, err
	}

	return stats, nil
}

// QueryNetworkStats1H returns network statistics from 1 hour network stats table.
func (db *Database) QueryNetworkStats1H(limit int) ([]schema.StatsNetwork1H, error) {
	var networkStats []schema.StatsNetwork1H
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(limit).
		Select()

	if err != nil {
		return networkStats, err
	}

	return networkStats, nil
}

// QueryNetworkStats1D returns 1 day network statistics.
func (db *Database) QueryNetworkStats1D(limit int) (stats []schema.StatsNetwork1D, err error) {
	err = db.Model(&stats).
		Order("id DESC").
		Limit(limit).
		Select()

	if err == pg.ErrNoRows {
		return []schema.StatsNetwork1D{}, nil
	}

	if err != nil {
		return []schema.StatsNetwork1D{}, err
	}

	return stats, nil
}

// QueryBondedRateIn1D return bonded rate in network from 1 day network stats table.
func (db *Database) QueryBondedRateIn1D() ([]schema.StatsNetwork1D, error) {
	var networkStats []schema.StatsNetwork1D
	err := db.Model(&networkStats).
		Order("id DESC").
		Limit(2).
		Select()

	if err != nil {
		return networkStats, err
	}

	return networkStats, nil
}

// --------------------
// Count
// --------------------

// CountProposedBlocks counts how many proposed blocks made by a proposer.
func (db *Database) CountProposedBlocks(proposer string) (int, error) {
	var block schema.Block
	count, err := db.Model(&block).
		Where("proposer = ?", proposer).
		Count()

	if err != nil {
		return -1, err
	}

	return count, nil
}

// CountMissingBlocks counts how many missing blocks a validator misses in detail and return total missing blocks count.
func (db *Database) CountMissingBlocks(address string, latestHeight int, count int) (int, error) {
	var misses []schema.MissDetail
	count, err := db.Model(&misses).
		Where("address = ? AND height BETWEEN ? AND ?", address, latestHeight-count, latestHeight).
		Count()

	if err != nil {
		return -1, err
	}

	return count, nil
}

// CountValidatorNumByStatus counts a number of validators by their bonding status.
func (db *Database) CountValidatorNumByStatus(status int) (int, error) {
	var val schema.Validator
	num, err := db.Model(&val).
		Where("status = ?", status).
		Count()

	if err != nil {
		return -1, err
	}

	return num, nil
}

// CountValidatorPowerEvents counts validator's power event history transactions.
func (db *Database) CountValidatorPowerEvents(proposer string) (int, error) {
	var peh schema.PowerEventHistory
	num, err := db.Model(&peh).
		Where("proposer = ?", proposer).
		Count()

	if err != nil {
		return -1, err
	}

	return num, nil
}

// CountMarketStats1H counts network statistics.
func (db *Database) CountMarketStats1H() (int, error) {
	var market schema.StatsMarket1H
	num, err := db.Model(&market).Count()
	if err != nil {
		return -1, err
	}

	return num, nil
}

// CountNetworkStats1H counts network statistics.
func (db *Database) CountNetworkStats1H() (int, error) {
	var network schema.StatsNetwork1H
	num, err := db.Model(&network).Count()
	if err != nil {
		return -1, err
	}

	return num, nil
}

// --------------------
// Exist
// --------------------
