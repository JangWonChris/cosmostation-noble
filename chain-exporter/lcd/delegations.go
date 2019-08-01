package lcd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	"github.com/go-pg/pg"
	resty "gopkg.in/resty.v1"
)

/*
	CURRENTLY THIS METHOD IS NOT USED
*/
// SaveValidatorDelegations queries each validator's delegations and save them in the database
func SaveValidatorDelegations(db *pg.DB, config *config.Config) {
	// Query all validators' operating addresses
	var validatorInfo []dtypes.ValidatorInfo
	_ = db.Model(&validatorInfo).
		Column("cosmos_address", "operator_address").
		Order("id ASC").
		Select()

	// Query each validator's delegations
	validatorDelegationsInfo := make([]*dtypes.ValidatorDelegationsInfo, 0)
	for _, validator := range validatorInfo {
		resp, _ := resty.R().Get("https://lcd-do-not-abuse.cosmostation.io/staking/validators/" + validator.OperatorAddress + "/delegations")

		// Parse ValidatorDelegations struct
		var validatorDelegations []dtypes.ValidatorDelegations
		_ = json.Unmarshal(resp.Body(), &validatorDelegations)

		// var totalDelegatorShares float64
		var selfDelegatedShares float64
		var othersShares float64

		for _, validatorDelegation := range validatorDelegations {
			// Calculate self-delegated and others shares
			if validatorDelegation.DelegatorAddress == validator.Address {
				selfDelegatedShares, _ = strconv.ParseFloat(validatorDelegation.Shares.String(), 64)
			} else {
				tempOthersShares, _ := strconv.ParseFloat(validatorDelegation.Shares.String(), 64)
				othersShares += tempOthersShares
			}
		}

		// Insert validator delegations data
		tempValidatorDelegationsInfo := &dtypes.ValidatorDelegationsInfo{
			Address:             validator.Address,
			OperatorAddress:     validator.OperatorAddress,
			TotalShares:         selfDelegatedShares + othersShares,
			SelfDelegatedShares: selfDelegatedShares,
			OthersShares:        othersShares,
			DelegatorNum:        len(validatorDelegations),
			Time:                time.Now(),
		}
		validatorDelegationsInfo = append(validatorDelegationsInfo, tempValidatorDelegationsInfo)
	}

	fmt.Println(len(validatorDelegationsInfo))

	// Save & Update validatorDelegationsInfo
	// _, err := db.Model(&validatorDelegationsInfo).
	// 	OnConflict("(operator_address) DO UPDATE").
	// 	Set("total_shares = EXCLUDED.total_shares").
	// 	Set("self_delegated_shares = EXCLUDED.self_delegated_shares").
	// 	Set("others_shares = EXCLUDED.others_shares").
	// 	Set("delegator_num = EXCLUDED.delegator_num").
	// 	Insert()
	// if err != nil {
	// 	fmt.Printf("error - save and update validatorinfo: %v\n", err)
	// }
}
