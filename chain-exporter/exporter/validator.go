package exporter

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	resty "gopkg.in/resty.v1"
)

// getValidatorSetInfo provides validator set information in every block
func (ces *ChainExporterService) getValidatorSetInfo(height int64) ([]*types.ValidatorSetInfo, []*types.MissInfo, []*types.MissInfo, []*types.MissDetailInfo, error) {
	nextHeight := height + 1

	// query current block
	block, err := ces.rpcClient.Block(&height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// query the next block to access precommits
	nextBlock, err := ces.rpcClient.Block(&nextHeight)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// query validator set for the block height
	validators, err := ces.rpcClient.Validators(&height)
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

	genesisValidatorsInfo := make([]*types.ValidatorSetInfo, 0)
	missInfo := make([]*types.MissInfo, 0)
	accumMissInfo := make([]*types.MissInfo, 0)
	missDetailInfo := make([]*types.MissDetailInfo, 0)

	// validator set for the height
	for i, validator := range validators.Validators {
		// insert genesis validators as an event_type of create_validator at height 1
		if validators.BlockHeight == 1 {
			tempValidatorSetInfo := &types.ValidatorSetInfo{
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
			tempMissDetailInfo := &types.MissDetailInfo{
				Height:   block.BlockMeta.Header.Height,
				Address:  validator.Address.String(),
				Proposer: block.BlockMeta.Header.ProposerAddress.String(),
				Alerted:  false,
				Time:     block.BlockMeta.Header.Time,
			}
			missDetailInfo = append(missDetailInfo, tempMissDetailInfo)

			// initial variables
			startHeight := block.BlockMeta.Header.Height
			endHeight := block.BlockMeta.Header.Height
			missingCount := int64(1)

			// query to check if a validator missed previous block
			var prevMissInfo types.MissInfo
			_ = ces.db.Model(&prevMissInfo).
				Where("end_height = ? AND address = ?", endHeight-int64(1), validator.Address.String()).
				Order("end_height DESC").
				Select()

			if prevMissInfo.Address == "" {
				tempMissInfo := &types.MissInfo{
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
				tempMissInfo := &types.MissInfo{
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
func (ces *ChainExporterService) getEvidenceInfo(height int64) ([]*types.EvidenceInfo, error) {
	nextHeight := height + 1

	// query current block
	block, err := ces.rpcClient.Block(&height)
	if err != nil {
		return nil, err
	}

	// query the next block to access precommits
	nextBlock, err := ces.rpcClient.Block(&nextHeight)
	if err != nil {
		return nil, err
	}

	// cosmoshub-2
	// 848187 = 1C4DB67E79B5BB30663B04245E064E6180EC6EA304EE83A7A879B04544A2EAD0
	// in evidence, there is only DuplicateVoteEvidence. There is no downtime evidence.
	evidenceInfo := make([]*types.EvidenceInfo, 0)
	if nextBlock.Block.Evidence.Evidence != nil {
		for _, evidence := range nextBlock.Block.Evidence.Evidence {
			tempEvidenceInfo := &types.EvidenceInfo{
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
func (ces *ChainExporterService) SaveValidatorKeyBase() error {
	var validatorInfo []types.ValidatorInfo
	err := ces.db.Model(&validatorInfo).
		Column("id", "identity", "moniker").
		Select()
	if err != nil {
		return err
	}

	validatorInfoUpdate := make([]*types.ValidatorInfo, 0)
	for _, validator := range validatorInfo {
		if validator.Identity != "" {
			resp, err := resty.R().Get(ces.config.KeybaseURL + validator.Identity)
			if err != nil {
				fmt.Printf("KeyBase request error - %v\n", err)
			}

			var keyBases types.KeyBase
			err = json.Unmarshal(resp.Body(), &keyBases)
			if err != nil {
				fmt.Printf("KeyBase unmarshal error - %v\n", err)
			}

			// get keybase urls
			var keybaseURL string
			if len(keyBases.Them) > 0 {
				for _, keybase := range keyBases.Them {
					keybaseURL = keybase.Pictures.Primary.URL
				}
			}

			tempValidatorInfo := &types.ValidatorInfo{
				ID:         validator.ID,
				KeybaseURL: keybaseURL,
			}
			validatorInfoUpdate = append(validatorInfoUpdate, tempValidatorInfo)
		}
	}

	if len(validatorInfoUpdate) > 0 {
		var tempValidatorInfo types.ValidatorInfo
		for i := 0; i < len(validatorInfoUpdate); i++ {
			_, err = ces.db.Model(&tempValidatorInfo).
				Set("keybase_url = ?", validatorInfoUpdate[i].KeybaseURL).
				Where("id = ?", validatorInfoUpdate[i].ID).
				Update()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
