package exporter

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	resty "gopkg.in/resty.v1"
)

// getValidatorSetInfo provides validator set information in every block
func (ex Exporter) getValidatorSetInfo(height int64) ([]*schema.ValidatorSetInfo, []*schema.MissInfo, []*schema.MissInfo, []*schema.MissDetailInfo, error) {
	nextHeight := height + 1

	// Query current block
	block, err := ex.client.Block(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Query the next block to access precommits
	nextBlock, err := ex.client.Block(nextHeight)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Query validator set for the block height
	validators, err := ex.client.Validators(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	/*
		DO NOT SORT validator set. This causes miss_infos incorrect data
	*/

	// Sort bondedValidators by highest tokens
	// sort.Slice(validators.Validators[:], func(i, j int) bool {
	// 	tempToken1 := validators.Validators[i].VotingPower
	// 	tempToken2 := validators.Validators[j].VotingPower
	// 	return tempToken1 > tempToken2
	// })

	genesisValidatorsInfo := make([]*schema.ValidatorSetInfo, 0)
	missInfo := make([]*schema.MissInfo, 0)
	accumMissInfo := make([]*schema.MissInfo, 0)
	missDetailInfo := make([]*schema.MissDetailInfo, 0)

	// Get validator set for the height
	for i, validator := range validators.Validators {
		// insert genesis validators as an event_type of create_validator at height 1
		if validators.BlockHeight == 1 {
			tempValidatorSetInfo := &schema.ValidatorSetInfo{
				IDValidator:          i + 1,
				Height:               block.Block.Height,
				Proposer:             validator.Address.String(),
				VotingPower:          float64(validator.VotingPower),
				NewVotingPowerAmount: float64(validator.VotingPower),
				NewVotingPowerDenom:  types.Denom,
				EventType:            types.EventTypeMsgCreateValidator,
				Time:                 block.BlockMeta.Header.Time,
			}
			genesisValidatorsInfo = append(genesisValidatorsInfo, tempValidatorSetInfo)
		}

		// MissDetailInfo saves every missing information of validators
		// MissInfo saves ranges of missing information of validators
		// check if a validator misses previous block
		if nextBlock.Block.LastCommit.Precommits[i] == nil {
			tempMissDetailInfo := &schema.MissDetailInfo{
				Height:   block.BlockMeta.Header.Height,
				Address:  validator.Address.String(),
				Proposer: block.BlockMeta.Header.ProposerAddress.String(),
				Alerted:  false,
				Time:     block.BlockMeta.Header.Time,
			}
			missDetailInfo = append(missDetailInfo, tempMissDetailInfo)

			// Initial variables
			startHeight := block.BlockMeta.Header.Height
			endHeight := block.BlockMeta.Header.Height
			missingCount := int64(1)

			// Query to check if a validator missed previous block
			var prevMissInfo schema.MissInfo
			_ = ex.db.Model(&prevMissInfo).
				Where("end_height = ? AND address = ?", endHeight-int64(1), validator.Address.String()).
				Order("end_height DESC").
				Select()

			if prevMissInfo.Address == "" {
				tempMissInfo := &schema.MissInfo{
					Address:      validator.Address.String(),
					StartHeight:  startHeight,
					EndHeight:    endHeight,
					MissingCount: missingCount,
					StartTime:    block.BlockMeta.Header.Time,
					EndTime:      block.BlockMeta.Header.Time,
					Alerted:      false,
				}
				missInfo = append(missInfo, tempMissInfo)
			} else {
				tempMissInfo := &schema.MissInfo{
					Address:      prevMissInfo.Address,
					StartHeight:  prevMissInfo.StartHeight,
					EndHeight:    prevMissInfo.EndHeight + int64(1),
					MissingCount: prevMissInfo.MissingCount + int64(1),
					StartTime:    prevMissInfo.StartTime,
					EndTime:      block.BlockMeta.Header.Time,
					Alerted:      false,
				}
				accumMissInfo = append(accumMissInfo, tempMissInfo)
			}
			continue
		}
	}
	return genesisValidatorsInfo, missInfo, accumMissInfo, missDetailInfo, nil
}

// getEvidenceInfo provides evidence (slashing) information
func (ex Exporter) getEvidenceInfo(height int64) ([]*schema.EvidenceInfo, error) {
	nextHeight := height + 1

	// Query current block
	block, err := ex.client.Block(height)
	if err != nil {
		return nil, err
	}

	// Query the next block to access precommits
	nextBlock, err := ex.client.Block(nextHeight)
	if err != nil {
		return nil, err
	}

	// cosmoshub-2
	// 848187 = 1C4DB67E79B5BB30663B04245E064E6180EC6EA304EE83A7A879B04544A2EAD0
	// in evidence, there is only DuplicateVoteEvidenex. There is no downtime evidenex.
	evidenceInfo := make([]*schema.EvidenceInfo, 0)
	if nextBlock.Block.Evidence.Evidence != nil {
		for _, evidence := range nextBlock.Block.Evidence.Evidence {
			tempEvidenceInfo := &schema.EvidenceInfo{
				Proposer: strings.ToUpper(string(hex.EncodeToString(evidence.Address()))),
				Height:   evidence.Height(),
				Hash:     nextBlock.BlockMeta.Header.EvidenceHash.String(),
				Time:     block.BlockMeta.Header.Time,
			}
			evidenceInfo = append(evidenceInfo, tempEvidenceInfo)
		}
	}

	return evidenceInfo, nil
}

// SaveValidatorKeyBase saves keybase urls for every validator
func (ex Exporter) SaveValidatorKeyBase() error {
	validatorsInfo := make([]*schema.ValidatorInfo, 0)

	// query validators info
	validators, _ := ex.db.QueryValidators()

	for _, validator := range validators {
		if validator.Identity != "" {
			resp, err := resty.R().Get(ex.cfg.KeybaseURL + validator.Identity)
			if err != nil {
				fmt.Printf("failed to request KeyBase: %v \n", err)
			}

			var keyBases types.KeyBase
			err = json.Unmarshal(resp.Body(), &keyBases)
			if err != nil {
				fmt.Printf("failed to unmarshal KeyBase: %v \n", err)
			}

			var keybaseURL string
			if len(keyBases.Them) > 0 {
				for _, keybase := range keyBases.Them {
					keybaseURL = keybase.Pictures.Primary.URL
				}
			}

			tempValidatorInfo := &schema.ValidatorInfo{
				ID:         validator.ID,
				KeybaseURL: keybaseURL,
			}
			validatorsInfo = append(validatorsInfo, tempValidatorInfo)
		}
	}

	if len(validatorsInfo) > 0 {
		for _, validator := range validatorsInfo {
			result, err := ex.db.UpdateKeyBase(validator.ID, validator.KeybaseURL)
			if !result {
				log.Printf("failed to update KeyBase URL: %v \n", err)
			}
		}
	}

	return nil
}
