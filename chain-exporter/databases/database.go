package databases

import (
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// ConnectDatabase connects to PostgreSQL
func ConnectDatabase(Config *config.Config) *pg.DB {
	database := pg.Connect(&pg.Options{
		Addr:     Config.DB.Host,
		User:     Config.DB.User,
		Password: Config.DB.Password,
		Database: Config.DB.Table,
	})
	return database
}

// CreateSchema creates database tables using ORM
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*schema.BlockInfo)(nil), (*schema.EvidenceInfo)(nil), (*schema.MissInfo)(nil),
		(*schema.MissDetailInfo)(nil), (*schema.ProposalInfo)(nil), (*schema.ValidatorSetInfo)(nil), (*schema.ValidatorInfo)(nil),
		(*schema.TransactionInfo)(nil), (*schema.VoteInfo)(nil), (*schema.DepositInfo)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     5000, // replaces PostgreSQL data type `text` to `varchar(n)`
		})
		if err != nil {
			return err
		}
	}

	// RunInTransaction runs a function in a transaction.
	// If function returns an error transaction is rollbacked, otherwise transaction is committed.
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		// Create indexes to reduce the cost of lookup queries in case of server traffic jams (B-Tree Index)
		_, err := db.Model(schema.BlockInfo{}).Exec(`CREATE INDEX block_info_height_idx ON block_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.ValidatorInfo{}).Exec(`CREATE INDEX validator_info_rank_idx ON validator_infos USING btree(rank);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.MissDetailInfo{}).Exec(`CREATE INDEX miss_detail_info_height_idx ON miss_detail_infos USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.MissInfo{}).Exec(`CREATE INDEX miss_info_start_height_idx ON miss_infos USING btree(start_height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.TransactionInfo{}).Exec(`CREATE INDEX transaction_info_height_idx ON transaction_infos USING btree(height);`)
		if err != nil {
			return err
		}

		return nil
	})

	// Roll back if any index creation fails.
	if err != nil {
		return err
	}

	return nil
}

// QueryValidatorInfo returns validator information
func QueryValidatorInfo(db *pg.DB, address string) (schema.ValidatorInfo, error) {
	var validatorInfo schema.ValidatorInfo
	switch {
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ConsensusPubPrefix()):
		err := db.Model(&validatorInfo).
			Where("address = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validatorInfo, err
		}
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ValidatorAddrPrefix()):
		err := db.Model(&validatorInfo).
			Where("operator_address = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validatorInfo, err
		}
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32AccountAddrPrefix()):
		err := db.Model(&validatorInfo).
			Where("consensus_pubkey = ?", address).
			Limit(1).
			Select()
		if err != nil {
			return validatorInfo, err
		}
	}
	return validatorInfo, nil
}

// QueryIDValidatorSetInfo returns id of a validator from validator_set_infos table
func QueryIDValidatorSetInfo(db *pg.DB, proposer string) (schema.ValidatorSetInfo, error) {
	var validatorSetInfo schema.ValidatorSetInfo
	err := db.Model(&validatorSetInfo).
		Column("id_validator", "voting_power").
		Where("proposer = ?", proposer).
		Order("id DESC"). // Lastly input data
		Limit(1).
		Select()
	if err != nil {
		return validatorSetInfo, err
	}
	return validatorSetInfo, nil
}

// QueryHighestIDValidatorNum returns highest id of a validator from validator_set_infos table
func QueryHighestIDValidatorNum(db *pg.DB) (int, error) {
	var validatorSetInfo schema.ValidatorSetInfo
	err := db.Model(&validatorSetInfo).
		Column("id_validator").
		Order("id_validator DESC").
		Limit(1).
		Select()
	if err != nil {
		return 0, err
	}
	return validatorSetInfo.IDValidator, nil
}
