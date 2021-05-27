package db

import (
	"context"
	"fmt"
	"strings"

	mblconfig "github.com/cosmostation/mintscan-backend-library/config"
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"
	mddb "github.com/cosmostation/mintscan-database/db"
	"github.com/cosmostation/mintscan-database/schema"
	mdschema "github.com/cosmostation/mintscan-database/schema"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	pg "github.com/go-pg/pg/v10"
)

var (
	// columnLength is the column length of varchar type in every table.
	// This needs to be considered again to set it to what specific length is needed, but right now set it to 99999.
	columnLength = 99999
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL.
type Database struct {
	*mddb.Database
}

// Connect opens a database connections with the given database connection info from config.
func Connect(dbcfg *mblconfig.DatabaseConfig) *Database {
	db := mddb.Connect(dbcfg.Host, dbcfg.Port, dbcfg.User, dbcfg.Password, dbcfg.DBName, dbcfg.Schema, dbcfg.Timeout)

	return &Database{db}
}

// CreateTables creates database tables using ORM (Object Relational Mapper).
func (db *Database) CreateTablesAndIndexes() {
	// 생성 오류 시 패닉
	db.CreateTables()
}

func (db *Database) QueryTxForPowerEventHistory(beginHeight, endHeight int64) ([]mdschema.RawTransaction, error) {
	var txs []mdschema.RawTransaction
	_, err := db.Query(&txs, "select t.* from stargate_final.raw_transaction t where exists ( select 1 from transaction_account as ta where height >= ? and height < ? and msg_type in (?, ?, ?, ?) and t.tx_hash = ta.tx_hash) order by t.height asc ", beginHeight, endHeight, mbltypes.StakingMsgCreateValidator, mbltypes.StakingMsgDelegate, mbltypes.StakingMsgBeginRedelegate, mbltypes.StakingMsgUndelegate)
	if err != nil {
		if err == pg.ErrNoRows {
			return txs, nil
		}
		return txs, err
	}

	return txs, nil
}

// QueryAccountMobile queries account information
func (db *Database) QueryAccountMobile(address string) (*mdschema.AccountMobile, error) {
	var account *mdschema.AccountMobile
	_ = db.Model(&account).
		Where("address = ?", address).
		Select()

	return account, nil
}

// QueryAlarmTokens queries user's alarm tokens
func (db *Database) QueryAlarmTokens(address string) ([]string, error) {
	var accounts []mdschema.AccountMobile
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

// QueryValidatorByAnyAddr returns a validator information by any type of address format
func (db *Database) QueryValidatorByAnyAddr(anyAddr string) (mdschema.Validator, error) {
	var val mdschema.Validator
	var err error

	switch {
	// jeonghwan
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32ConsensusPubPrefix()): // Bech32 prefix for validator public key
		err = db.Model(&val).
			Where("consensus_pubkey = ?", anyAddr).
			Limit(1).
			Select()
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32ValidatorAddrPrefix()): // Bech32 prefix for validator address
		err = db.Model(&val).
			Where("operator_address = ?", anyAddr).
			Limit(1).
			Select()
	case strings.HasPrefix(anyAddr, sdktypes.GetConfig().GetBech32AccountAddrPrefix()): // Bech32 prefix for account address
		err = db.Model(&val).
			Where("address = ?", anyAddr).
			Limit(1).
			Select()
	case len(anyAddr) == 40: // Validator consensus address in hex
		anyAddr := strings.ToUpper(anyAddr)
		err = db.Model(&val).
			Where("proposer = ?", anyAddr).
			Limit(1).
			Select()
	default:
		err = db.Model(&val).
			Where("moniker = ?", anyAddr). // Validator moniker
			Limit(1).
			Select()
	}

	if err != nil {
		if err == pg.ErrNoRows {
			return mdschema.Validator{}, nil
		}
		return mdschema.Validator{}, err
	}

	return val, nil
}

// InsertExportedData saves exported blockchain data
// if function returns an error transaction is rollbacked, otherwise transaction is committed
func (db *Database) InsertExportedData(e *schema.BasicData) error {
	err := db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		if e.Block != nil {
			_, err := tx.Model(e.Block).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result block: %s", err)
			}
		}

		if len(e.GenesisAccounts) > 0 {
			_, err := tx.Model(&e.GenesisAccounts).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result genesis accounts: %s", err)
			}
		}

		if len(e.Evidence) > 0 {
			_, err := tx.Model(&e.Evidence).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result evidence: %s", err)
			}
		}

		if len(e.GenesisValidatorsSet) > 0 {
			_, err := tx.Model(&e.GenesisValidatorsSet).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result genesis validator set: %s", err)
			}
		}

		if len(e.ValidatorsPowerEventHistory) > 0 {
			_, err := tx.Model(&e.ValidatorsPowerEventHistory).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result validator power event history: %s", err)
			}
		}

		if len(e.MissBlocks) > 0 {
			_, err := tx.Model(&e.MissBlocks).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result miss blocks: %s", err)
			}
		}

		if len(e.MissDetailBlocks) > 0 {
			_, err := tx.Model(&e.MissDetailBlocks).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result miss detail blocks: %s", err)
			}
		}

		if len(e.Transactions) > 0 {
			for i := range e.Transactions {
				if e.Block.ID != 0 {
					e.Transactions[i].BlockID = e.Block.ID
				} else {
					return fmt.Errorf("failed to insert result txs, can not get block.id")
				}
			}
			_, err := tx.Model(&e.Transactions).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result txs: %s", err)
			}
		}

		lenTMA := len(e.SourceTransactionMessageAccounts)
		if lenTMA > 0 {
			limit := 100
			args := make([]string, 0)
			for i := 0; i < lenTMA; i += limit {
				if i+limit > lenTMA {
					limit = lenTMA - i
				}
				args = append(args, parseTMAToArg(e.SourceTransactionMessageAccounts[i:i+limit]))
			}
			for i := range args {
				query := "select public.f_insert_tx_msg_acc" + args[i] // refine. 과 같은 스키마 명을 놓치지 않도록 조심해야 함.
				_, err := tx.Exec(query)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		if len(e.Deposits) > 0 {
			_, err := tx.Model(&e.Deposits).Insert()
			if err != nil {
				return fmt.Errorf("failed to insert result deposits: %s", err)
			}
		}

		if len(e.AccumulatedMissBlocks) > 0 {
			for _, rmb := range e.AccumulatedMissBlocks {
				_, err := tx.Model(&schema.Miss{}).
					Set("address = ?", rmb.Address).
					Set("start_height = ?", rmb.StartHeight).
					Set("end_height = ?", rmb.EndHeight).
					Set("missing_count = ?", rmb.MissingCount).
					Set("start_time = ?", rmb.StartTime).
					Set("end_time = ?", e.Block.Timestamp).
					Where("end_height = ? AND address = ?", rmb.EndHeight-int64(1), rmb.Address).
					Update()
				if err != nil {
					return fmt.Errorf("failed to update result accumulated miss blocks: %s", err)
				}
			}
		}

		if len(e.Proposals) > 0 {
			for _, rp := range e.Proposals {
				ok, _ := db.ExistProposal(rp.ID)
				if !ok {
					_, err := tx.Model(&rp).Insert()
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

		if len(e.Votes) > 0 {
			for _, rv := range e.Votes {
				ok, _ := db.ExistVote(rv.ProposalID, rv.Voter)
				if !ok {
					_, err := tx.Model(&rv).Insert()
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

func parseTMAToArg(tma []schema.TMA) string {
	arg := "("
	fmt.Println("len of partial tma", len(tma))
	for i := range tma {
		if i != 0 {
			arg += ","
		}
		arg += fmt.Sprintf("('%s', '%s', '%s')", tma[i].TxHash, tma[i].MsgType, tma[i].AccountAddress)
	}
	arg += ")"
	return arg
}
