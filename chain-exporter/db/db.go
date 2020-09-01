package db

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	// columnLength is the column length of varchar type in every table.
	// This needs to be considered again to set it to what specific length is needed, but right now set it to 99999.
	columnLength = 99999
)

const (
	// Define PostgreSQL database indexes to improve the speed of data retrieval operations on a database tables.
	indexAccountAddress          = "CREATE INDEX accunt_account_address_idx ON account USING btree(account_address);"
	indexBlockHeight             = "CREATE INDEX block_height_idx ON block USING btree(height);"
	indexValidatorRank           = "CREATE INDEX validator_rank_idx ON validator USING btree(rank);"
	indexValidatorStatus         = "CREATE INDEX validator_status_idx ON validator USING btree(status);"
	indexPowerEventHistoryHeight = "CREATE INDEX power_event_history_height_idx ON power_event_history USING btree(height);"
	indexMissStartHeight         = "CREATE INDEX miss_start_height_idx ON miss USING btree(start_height);"
	indexMissEndHeight           = "CREATE INDEX miss_end_height_idx ON miss USING btree(end_height);"
	indexMissDetailHeight        = "CREATE INDEX miss_detail_height_idx ON miss_detail USING btree(height);"
	indexTransactionHeight       = "CREATE INDEX transaction_height_idx ON transaction USING btree(height);"
	indexTransactionChainID      = "CREATE INDEX transaction_chain_id_idx ON transaction USING btree(chain_id);"
	indexTransactionHash         = "CREATE INDEX transaction_tx_hash_idx ON transaction USING btree(tx_hash);"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*pg.DB
}

// Connect opens a database connections with the given database connection info from config.
func Connect(config *config.Database) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     config.Host + ":" + config.Port,
		User:     config.User,
		Password: config.Password,
		Database: config.Table,
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

// CreateTables creates database tables using ORM (Object Relational Mapper).
func (db *Database) CreateTables() error {
	for _, table := range []interface{}{(*schema.Block)(nil), (*schema.Evidence)(nil), (*schema.Miss)(nil),
		(*schema.MissDetail)(nil), (*schema.Proposal)(nil), (*schema.PowerEventHistory)(nil), (*schema.Validator)(nil),
		(*schema.Transaction)(nil), (*schema.Vote)(nil), (*schema.Deposit)(nil)} {

		// Disable pluralization
		orm.SetTableNameInflector(func(s string) string {
			return s
		})

		err := db.CreateTable(table, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     columnLength, // replaces PostgreSQL data type `text` to `varchar(n)`
		})

		if err != nil {
			return err
		}
	}

	// Create table indexes and roll back if any index creation fails.
	err := db.createIndexes()
	if err != nil {
		return err
	}

	return nil
}

// createIndexes uses RunInTransaction to run a function in a transaction.
// if function returns an error, transaction is rollbacked, otherwise transaction is committed.
// Create B-Tree indexes to reduce the cost of lookup queries
func (db *Database) createIndexes() error {
	db.RunInTransaction(func(tx *pg.Tx) (err error) {
		_, err = db.Model(schema.Block{}).Exec(indexBlockHeight)
		if err != nil {
			return fmt.Errorf("failed to create block height index: %s", err)
		}
		_, err = db.Model(schema.Validator{}).Exec(indexValidatorRank)
		if err != nil {
			return fmt.Errorf("failed to create validator rank index: %s", err)
		}
		_, err = db.Model(schema.Validator{}).Exec(indexValidatorStatus)
		if err != nil {
			return fmt.Errorf("failed to create validator status index: %s", err)
		}
		_, err = db.Model(schema.PowerEventHistory{}).Exec(indexPowerEventHistoryHeight)
		if err != nil {
			return fmt.Errorf("failed to create power event history height index: %s", err)
		}
		_, err = db.Model(schema.Miss{}).Exec(indexMissStartHeight)
		if err != nil {
			return fmt.Errorf("failed to create miss start height index: %s", err)
		}
		_, err = db.Model(schema.Miss{}).Exec(indexMissEndHeight)
		if err != nil {
			return fmt.Errorf("failed to create miss end height index: %s", err)
		}
		_, err = db.Model(schema.MissDetail{}).Exec(indexMissDetailHeight)
		if err != nil {
			return fmt.Errorf("failed to create miss detail height index: %s", err)
		}
		_, err = db.Model(schema.Transaction{}).Exec(indexTransactionHeight)
		if err != nil {
			return fmt.Errorf("failed to create tx height index: %s", err)
		}
		_, err = db.Model(schema.Transaction{}).Exec(indexTransactionChainID)
		if err != nil {
			return fmt.Errorf("failed to create tx chain id index: %s", err)
		}
		_, err = db.Model(schema.Transaction{}).Exec(indexTransactionHash)
		if err != nil {
			return fmt.Errorf("failed to create tx hash index: %s", err)
		}

		return nil
	})

	return nil
}

// --------------------
// Query
// --------------------

// QueryLatestBlockHeight queries latest block height in database
func (db *Database) QueryLatestBlockHeight() (int64, error) {
	var block schema.Block
	err := db.Model(&block).
		Order("height DESC").
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

	return block.Height, nil
}

// QueryValidators returns validators info
func (db *Database) QueryValidators() ([]schema.Validator, error) {
	var vals []schema.Validator
	err := db.Model(&vals).
		Column("id", "identity", "moniker").
		Select()

	if err != nil {
		return vals, err
	}

	return vals, nil
}

// QueryValidator returns validator information
func (db *Database) QueryValidator(address string) (schema.Validator, error) {
	var validator schema.Validator
	switch {
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ConsensusPubPrefix()):
		err := db.Model(&validator).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validator, err
		}
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ValidatorAddrPrefix()):
		err := db.Model(&validator).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validator, err
		}
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32AccountAddrPrefix()):
		err := db.Model(&validator).
			Where("address = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validator, err
		}
	}
	return validator, nil
}

// QueryValidatorID returns the id number of a validator from power_event_history table
func (db *Database) QueryValidatorID(proposer string) (schema.PowerEventHistory, error) {
	var powerEventHistory schema.PowerEventHistory
	err := db.Model(&powerEventHistory).
		Column("id_validator", "voting_power").
		Where("proposer = ?", proposer).
		Order("id DESC"). // Lastly input data
		Limit(1).
		Select()
	if err != nil {
		return powerEventHistory, err
	}
	return powerEventHistory, nil
}

// QueryHighestValidatorID returns highest id number of a validator from power_event_history table
func (db *Database) QueryHighestValidatorID() (int, error) {
	var powerEventHistory schema.PowerEventHistory
	err := db.Model(&powerEventHistory).
		Column("id_validator").
		Order("id_validator DESC").
		Limit(1).
		Select()
	if err != nil {
		return 0, err
	}
	return powerEventHistory.IDValidator, nil
}

// QueryAccountMobile queries account information
func (db *Database) QueryAccountMobile(address string) (schema.AccountMobile, error) {
	var account schema.AccountMobile
	_ = db.Model(&account).
		Where("address = ?", address).
		Select()

	return account, nil
}

// QueryAlarmTokens queries user's alarm tokens
func (db *Database) QueryAlarmTokens(address string) ([]string, error) {
	var accounts []schema.AccountMobile
	_ = db.Model(&accounts).
		Column("alarm_token").
		Where("address = ?", address).
		Select()

	var result []string
	if len(accounts) > 0 {
		for _, account := range accounts {
			result = append(result, account.AlarmToken)
		}
	}

	return result, nil
}

// QueryHighestRankValidatorByStatus queries highest rank of a validator by status
func (db *Database) QueryHighestRankValidatorByStatus(status int) int {
	var val schema.Validator
	_ = db.Model(&val).
		Where("status = ?", status).
		Order("rank DESC").
		Limit(1).
		Select()

	return val.Rank
}

// QueryMissingPreviousBlock queries if a validator has missed previous block.
func (db *Database) QueryMissingPreviousBlock(consAddrHex string, prevHeight int64) schema.Miss {
	var prevMiss schema.Miss
	_ = db.Model(&prevMiss).
		Where("end_height = ? AND address = ?", prevHeight, consAddrHex).
		Order("end_height DESC").
		Select()

	return prevMiss
}

// --------------------
// Exists
// --------------------

// ExistProposal queries to find if the proposal id exists in database.
func (db *Database) ExistProposal(proposalID int64) (exist bool, err error) {
	var proposal schema.Proposal
	exist, err = db.Model(&proposal).
		Where("id = ?", proposalID).
		Exists()

	if err != nil {
		return false, err
	}

	return exist, nil
}

// ExistVote checks to see if a vote exists in database.
func (db *Database) ExistVote(proposalID uint64, voter string) (exist bool, err error) {
	var vote schema.Vote
	exist, err = db.Model(&vote).
		Where("proposal_id = ? AND voter = ?", proposalID, voter).
		Exists()

	if err != nil {
		return false, err
	}

	return exist, nil
}

// ExistValidator checks to see if a validator exists
func (db *Database) ExistValidator(valAddr string) (bool, error) {
	var validator schema.Validator
	exist, err := db.Model(&validator).
		Where("operator_address = ?", valAddr).
		Exists()

	if err != nil {
		return false, err
	}

	return exist, nil
}

// --------------------
// Insert or Update
// --------------------

// InsertOrUpdateProposals inserts if not exist already or updates proposals.
func (db *Database) InsertOrUpdateProposals(proposals []schema.Proposal) error {
	var proposal schema.Proposal
	for _, p := range proposals {
		ok, _ := db.ExistProposal(p.ID)
		if !ok {
			err := db.Insert(&p)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Model(&proposal).
				Set("title = ?", p.Title).
				Set("description = ?", p.Description).
				Set("proposal_type = ?", p.ProposalType).
				Set("proposal_status = ?", p.ProposalStatus).
				Set("yes = ?", p.Yes).
				Set("abstain = ?", p.Abstain).
				Set("no = ?", p.No).
				Set("no_with_veto = ?", p.NoWithVeto).
				Set("deposit_end_time = ?", p.DepositEndtime).
				Set("total_deposit_amount = ?", p.TotalDepositAmount).
				Set("total_deposit_denom = ?", p.TotalDepositDenom).
				Set("submit_time = ?", p.SubmitTime).
				Set("voting_start_time = ?", p.VotingStartTime).
				Set("voting_end_time = ?", p.VotingEndTime).
				Where("id = ?", p.ID).
				Update()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// InsertOrUpdateValidators inserts validators or updates validators information.
func (db *Database) InsertOrUpdateValidators(vals []schema.Validator) error {
	var validator schema.Validator
	for _, val := range vals {
		ok, _ := db.ExistValidator(val.OperatorAddress)
		if !ok {
			err := db.Insert(&val)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Model(&validator).
				Set("rank = ?", val.Rank).
				Set("consensus_pubkey = ?", val.ConsensusPubkey).
				Set("proposer = ?", val.Proposer).
				Set("jailed = ?", val.Jailed).
				Set("status = ?", val.Status).
				Set("tokens = ?", val.Tokens).
				Set("delegator_shares = ?", val.DelegatorShares).
				Set("moniker = ?", val.Moniker).
				Set("identity = ?", val.Identity).
				Set("website = ?", val.Website).
				Set("details = ?", val.Details).
				Set("unbonding_height = ?", val.UnbondingHeight).
				Set("unbonding_time = ?", val.UnbondingTime).
				Set("commission_rate = ?", val.CommissionRate).
				Set("commission_max_rate = ?", val.CommissionMaxRate).
				Set("commission_change_rate = ?", val.CommissionChangeRate).
				Set("update_time = ?", val.UpdateTime).
				Set("min_self_delegation = ?", val.MinSelfDelegation).
				Where("operator_address = ?", val.OperatorAddress).
				Update()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// InsertExportedData saves exported blockchain data
// if function returns an error transaction is rollbacked, otherwise transaction is committed.
func (db *Database) InsertExportedData(e schema.ExportData) error {
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		err := tx.Insert(&e.ResultBlock)
		if err != nil {
			return fmt.Errorf("failed to insert result block: %s", err)
		}

		if len(e.ResultEvidence) > 0 {
			err := tx.Insert(&e.ResultEvidence)
			if err != nil {
				return fmt.Errorf("failed to insert result evidence: %s", err)
			}
		}

		if len(e.ResultGenesisValidatorsSet) > 0 {
			err := tx.Insert(&e.ResultGenesisValidatorsSet)
			if err != nil {
				return fmt.Errorf("failed to insert result genesis validator set: %s", err)
			}
		}

		if len(e.ResultValidatorsPowerEventHistory) > 0 {
			err := tx.Insert(&e.ResultValidatorsPowerEventHistory)
			if err != nil {
				return fmt.Errorf("failed to insert result validator power event history: %s", err)
			}
		}

		if len(e.ResultMissBlocks) > 0 {
			err := tx.Insert(&e.ResultMissBlocks)
			if err != nil {
				return fmt.Errorf("failed to insert result miss blocks: %s", err)
			}
		}

		if len(e.ResultMissDetailBlocks) > 0 {
			err := tx.Insert(&e.ResultMissDetailBlocks)
			if err != nil {
				return fmt.Errorf("failed to insert result miss detail blocks: %s", err)
			}
		}

		if len(e.ResultTxs) > 0 {
			err := tx.Insert(&e.ResultTxs)
			if err != nil {
				return fmt.Errorf("failed to insert result txs: %s", err)
			}
		}

		if len(e.ResultDeposits) > 0 {
			err := tx.Insert(&e.ResultDeposits)
			if err != nil {
				return fmt.Errorf("failed to insert result deposits: %s", err)
			}
		}

		if len(e.ResultAccumulatedMissBlocks) > 0 {
			var miss schema.Miss
			for _, rmb := range e.ResultAccumulatedMissBlocks {
				_, err := tx.Model(&miss).
					Set("address = ?", rmb.Address).
					Set("start_height = ?", rmb.StartHeight).
					Set("end_height = ?", rmb.EndHeight).
					Set("missing_count = ?", rmb.MissingCount).
					Set("start_time = ?", rmb.StartTime).
					Set("end_time = ?", e.ResultBlock.Timestamp).
					Where("end_height = ? AND address = ?", rmb.EndHeight-int64(1), rmb.Address).
					Update()

				if err != nil {
					return fmt.Errorf("failed to update result accumulated miss blocks: %s", err)
				}
			}
		}

		if len(e.ResultProposals) > 0 {
			var proposal schema.Proposal
			for _, rp := range e.ResultProposals {
				ok, _ := db.ExistProposal(rp.ID)
				if !ok {
					err := tx.Insert(&rp)
					if err != nil {
						return fmt.Errorf("failed to insert result proposal: %s", err)
					}
				} else {
					_, err := tx.Model(&proposal).
						Set("tx_hash = ?", rp.TxHash).
						Set("proposer = ?", rp.Proposer).
						Set("initial_deposit_amount = ?", rp.InitialDepositAmount).
						Set("initial_deposit_denom = ?", rp.InitialDepositDenom).
						Where("id = ?", rp.ID).
						Update()

					if err != nil {
						return fmt.Errorf("failed to update result proposal: %s", err)
					}
				}

			}
		}

		if len(e.ReusltVotes) > 0 {
			var vote schema.Vote
			for _, rv := range e.ReusltVotes {
				ok, _ := db.ExistVote(rv.ProposalID, rv.Voter)
				if !ok {
					err := tx.Insert(&rv)
					if err != nil {
						return fmt.Errorf("failed to insert result vote: %s", err)
					}
				} else {
					_, err := tx.Model(&vote).
						Set("height = ?", rv.Height).
						Set("option = ?", rv.Option).
						Set("tx_hash = ?", rv.TxHash).
						Set("gas_wanted = ?", rv.GasWanted).
						Set("gas_used = ?", rv.GasUsed).
						Set("timestamp = ?", rv.Timestamp).
						Where("proposal_id = ? AND voter = ?", rv.ProposalID, rv.Voter).
						Update()

					if err != nil {
						return fmt.Errorf("failed to update result vote: %s", err)
					}
				}
			}
		}

		return nil
	})

	// Roll back if any insertion fails.
	if err != nil {
		return err
	}

	return nil
}

// UpdateValidatorsKeyBaseURL updates the given validators' keybase url information
func (db *Database) UpdateValidatorsKeyBaseURL(vals []schema.Validator) (bool, error) {
	var validator schema.Validator

	for _, val := range vals {
		_, err := db.Model(&validator).
			Set("keybase_url = ?", val.KeybaseURL).
			Where("id = ?", val.ID).
			Update()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
