package exporter

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/utils"

	resty "gopkg.in/resty.v1"
)

// Query the current delegation between a delegator and a validator > https://lcd.cosmostation.io/staking/delegators/cosmos1lh76pvy2ew873hxxnvl5gqpsgwzm53qpf6x78y/delegations/cosmosvaloper16m93gjfqvnjajzrfyszml8qm92a0w67nwxrca7
// Get all delegations from a delegator > https://lcd.cosmostation.io/staking/validators/cosmosvaloper16m93gjfqvnjajzrfyszml8qm92a0w67nwxrca7/delegations
func (ses *ChainExporterService) SaveValidatorStats() {
	log.Println("Save Validator Stats")

	// Query all validators order by their tokens
	var validatorInfo []types.ValidatorInfo
	err := ses.db.Model(&validatorInfo).
		Order("rank ASC").
		Select()
	if err != nil {
		fmt.Printf("ValidatorInfo DB error - %v\n", err)
	}

	// Validator rank from 1 to 100
	validatorStats := make([]*types.ValidatorStats, 0)
	for _, validator := range validatorInfo {
		// Convert to validator's cosmos address
		cosmosAddress := utils.ConvertOperatorAddressToCosmosAddress(validator.OperatorAddress)

		// Get self-bonded amount
		selfBondedResp, err := resty.R().Get(config.Node.LCDUrl + "/staking/delegators/" + cosmosAddress + "/delegations/" + validator.OperatorAddress)
		if err != nil {
			fmt.Printf("Staking Pool LCD resty - %v\n", err)
		}
		var delegatorsDelegation types.DelegatorsDelegation
		err = json.Unmarshal(selfBondedResp.Body(), &delegatorsDelegation)
		if err != nil {
			fmt.Printf("DelegatorsDelegations unmarshal error - %v\n", err)
		}

		// Get validator's delegation Amount
		valiDelegationResp, err := resty.R().Get(config.Node.LCDUrl + "/staking/validators/" + validator.OperatorAddress + "/delegations")
		if err != nil {
			fmt.Printf("Staking Pool LCD resty - %v\n", err)
		}
		var validatorDelegations []types.DelegatorsDelegation
		err = json.Unmarshal(valiDelegationResp.Body(), &validatorDelegations)
		if err != nil {
			fmt.Printf("DelegatorsDelegations unmarshal error - %v\n", err)
		}

		// Initialize variables, otherwise throws an error if there is no delegations
		var selfBondedAmount sdk.Dec
		selfBondedAmount = sdk.NewDec(0)

		// LCD 요청 값이 없을 경우 아래와 같이 에러가 발생하므로 아래 if 문으로 에러 처리. 추후에 Cosmos SDK 에서 에러 처리를 변경 할 걸로 보인다.
		// {  %!v(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)}
		if delegatorsDelegation.DelegatorAddress != "" {
			selfBondedAmount = delegatorsDelegation.Shares
		}

		var totalDelegations sdk.Dec
		totalDelegations = sdk.NewDec(0)
		if len(validatorDelegations) > 0 {
			for _, delegation := range validatorDelegations {
				if delegation.DelegatorAddress != cosmosAddress { // Go Decimal 관련 삽질
					tempDelegation := delegation.Shares
					totalDelegations = totalDelegations.Add(tempDelegation)
				}
			}
		}

		// DelegatorNum
		delegatorNum := len(validatorDelegations)

		tempValidatorStats := &types.ValidatorStats{
			Moniker:             validator.Moniker,
			OperatorAddress:     validator.OperatorAddress,
			CosmosAddress:       cosmosAddress,
			ProposerAddress:     validator.Proposer,
			SelfBonded1H:        selfBondedAmount.String(),
			DelegatorShares1H:   totalDelegations.String(),
			DelegatorNum1H:      delegatorNum,
			VotingPowerChange1H: selfBondedAmount.Add(totalDelegations).String(),
			Time:                time.Now(),
		}
		validatorStats = append(validatorStats, tempValidatorStats)
	}

	// Save
	_, err = ses.db.Model(&validatorStats).Insert()
	if err != nil {
		fmt.Printf("Save ValidatorStats error - %v\n", err)
	}
}
