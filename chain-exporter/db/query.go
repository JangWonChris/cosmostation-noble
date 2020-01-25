package db

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

// QueryValidators returns validators info
func (db *Database) QueryValidators() ([]schema.ValidatorInfo, error) {
	var validators []schema.ValidatorInfo
	err := db.Model(&validators).
		Column("id", "identity", "moniker").
		Select()
	if err != nil {
		return validators, err
	}
	return validators, nil
}

// QueryValidator returns validator information
func (db *Database) QueryValidator(address string) (schema.ValidatorInfo, error) {
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

// QueryValidatorID returns the id number of a validator from validator_set_infos table
func (db *Database) QueryValidatorID(proposer string) (schema.ValidatorSetInfo, error) {
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

// QueryHighestValidatorID returns highest id number of a validator from validator_set_infos table
func (db *Database) QueryHighestValidatorID() (int, error) {
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
func (db *Database) QueryFirstRankValidatorByStatus(status int) (schema.ValidatorInfo, error) {
	var rankInfo schema.ValidatorInfo
	_ = db.Model(&rankInfo).
		Where("status = ?", status).
		Order("rank DESC").
		Limit(1).
		Select()

	return rankInfo, nil
}

// QueryExistProposal queries to find out if the same proposal is already saved
func (db *Database) QueryExistProposal(proposalID int64) (bool, error) {
	var proposalInfo schema.ProposalInfo
	exist, _ := db.Model(&proposalInfo).
		Where("id = ?", proposalID).
		Exists()
	if !exist {
		return exist, nil
	}
	return exist, nil
}
