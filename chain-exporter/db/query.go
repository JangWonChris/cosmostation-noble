package db

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

// QueryLatestBlockHeight queries latest block height in database
func (db *Database) QueryLatestBlockHeight() (int64, error) {
	var block schema.BlockCosmoshub3
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
	var validators []schema.Validator
	err := db.Model(&validators).
		Column("id", "identity", "moniker").
		Select()
	if err != nil {
		return validators, err
	}
	return validators, nil
}

// QueryValidator returns validator information
func (db *Database) QueryValidator(address string) (schema.Validator, error) {
	var validator schema.Validator
	switch {
	case strings.HasPrefix(address, sdk.GetConfig().GetBech32ConsensusPubPrefix()):
		err := db.Model(&validator).
			Where("address = ?", address).
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
			Where("consensus_pubkey = ?", address).
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

// QueryAccount queries account information
func (db *Database) QueryAccount(address string) (types.Account, error) {
	var account types.Account
	_ = db.Model(&account).
		Where("address = ?", address).
		Select()

	return account, nil
}

// QueryAlarmTokens queries user's alarm tokens
func (db *Database) QueryAlarmTokens(address string) ([]string, error) {
	var accounts []types.Account
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

// QueryFirstRankValidatorByStatus queries highest rank of a validator by status
func (db *Database) QueryFirstRankValidatorByStatus(status int) (schema.Validator, error) {
	var rank schema.Validator
	_ = db.Model(&rank).
		Where("status = ?", status).
		Order("rank DESC").
		Limit(1).
		Select()

	return rank, nil
}

// ExistProposal queries to find out if the same proposal is already saved
func (db *Database) ExistProposal(proposalID int64) (bool, error) {
	var proposal schema.Proposal
	exist, _ := db.Model(&proposal).
		Where("id = ?", proposalID).
		Exists()
	if !exist {
		return exist, nil
	}
	return exist, nil
}

// ExistValidator checks to see if a validator exists
func (db *Database) ExistValidator(valAddr string) (bool, error) {
	var validator schema.Proposal
	ok, err := db.Model(&validator).
		Where("validator_address = ?", valAddr).
		Exists()

	if err != nil {
		return ok, err
	}

	return ok, nil
}
