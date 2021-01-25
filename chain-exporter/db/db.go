package db

import (
	"fmt"
	"log"

	//mbl
	lconfig "github.com/cosmostation/mintscan-backend-library/config"
	ldb "github.com/cosmostation/mintscan-backend-library/db"
	"github.com/cosmostation/mintscan-backend-library/db/schema"

	"github.com/go-pg/pg"
)

var (
	// columnLength is the column length of varchar type in every table.
	// This needs to be considered again to set it to what specific length is needed, but right now set it to 99999.
	columnLength = 99999
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*ldb.Database
}

// Connect opens a database connections with the given database connection info from config.
func Connect(config *lconfig.DatabaseConfig) *Database {
	db := ldb.Connect(config)

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
func (db *Database) CreateTablesAndIndexes() {
	// 생성 오류 시 패닉
	db.CreateTables()
}

// --------------------
// Query
// --------------------

// QueryAccountMobile queries account information
func (db *Database) QueryAccountMobile(address string) (*schema.AccountMobile, error) {
	var account *schema.AccountMobile
	_ = db.Model(&account).
		Where("address = ?", address).
		Select()

	return account, nil
}

// --------------------
// Insert or Update
// --------------------

// InsertOrUpdateAccounts inserts if not exist already or update accounts.
func (db *Database) InsertOrUpdateAccounts(accounts []schema.AccountCoin) error {
	for _, acc := range accounts {
		log.Println(acc.AccountAddress)
		log.Println(acc.Available)
		ok, _ := db.ExistAccountAtAccountCoin(acc.AccountAddress)
		if !ok {
			err := db.Insert(&acc)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Model(&schema.AccountCoin{}).
				Set("denom = ?", acc.Denom).
				Set("total = ?", acc.Total).
				Set("available = ?", acc.Available).
				Set("delegated = ?", acc.Delegated).
				Set("undelegated = ?", acc.Undelegated).
				Set("rewards = ?", acc.Rewards).
				Set("commission = ?", acc.Commission).
				Set("vesting = ?", acc.Vesting).
				Set("vested = ?", acc.Vested).
				Set("failed_vested = ?", acc.FailedVested).
				Set("last_tx = ?", acc.LastTx).
				Set("last_tx_time = ?", acc.LastTxTime).
				Where("account_address = ?", acc.AccountAddress).
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
func (db *Database) InsertExportedData(e *schema.ExportData) error {
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if e.ResultBlock != nil {
			err := tx.Insert(e.ResultBlock)
			if err != nil {
				return fmt.Errorf("failed to insert result block: %s", err)
			}
		}

		if len(e.ResultAccounts) > 0 {
			// err := db.InsertOrUpdateAccounts(e.ResultAccounts)
			// if err != nil {
			// 	return fmt.Errorf("failed to insert result accounts: %s", err)
			// }
		}

		if len(e.ResultAccountCoin) > 0 {
			log.Println("insert in ", e.ResultAccountCoin)
			err := db.InsertOrUpdateAccounts(e.ResultAccountCoin)
			if err != nil {
				return fmt.Errorf("failed to insert result account coin: %s", err)
			}
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

		if len(e.ResultTxsAccount) > 0 {
			err := tx.Insert(&e.ResultTxsAccount)
			if err != nil {
				return fmt.Errorf("failed to insert result txs message: %s", err)
			}
		}

		if len(e.ResultDeposits) > 0 {
			err := tx.Insert(&e.ResultDeposits)
			if err != nil {
				return fmt.Errorf("failed to insert result deposits: %s", err)
			}
		}

		if len(e.ResultAccumulatedMissBlocks) > 0 {
			for _, rmb := range e.ResultAccumulatedMissBlocks {
				_, err := tx.Model(&schema.Miss{}).
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
			for _, rp := range e.ResultProposals {
				ok, _ := db.ExistProposal(rp.ID)
				if !ok {
					err := tx.Insert(&rp)
					if err != nil {
						return err
					}
				} else {
					_, err := tx.Model(&schema.Proposal{}).
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

		if len(e.ResultVotes) > 0 {
			for _, rv := range e.ResultVotes {
				ok, _ := db.ExistVote(rv.ProposalID, rv.Voter)
				if !ok {
					err := tx.Insert(&rv)
					if err != nil {
						return fmt.Errorf("failed to insert result votes: %s", err)
					}
				} else {
					_, err := tx.Model(&schema.Vote{}).
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

// InsertGenesisAccount insert the given genesis accounts using Copy command, it will faster than insert
// func (db *Database) InsertGenesisAccount(acc []schema.AccountCoin) error {
// 	err := db.RunInTransaction(func(tx *pg.Tx) error {
// 		if len(acc) > 0 {
// 			err := tx.Insert(&acc)
// 			if err != nil {
// 				return fmt.Errorf("failed to insert result genesis accounts: %s", err)
// 			}
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
