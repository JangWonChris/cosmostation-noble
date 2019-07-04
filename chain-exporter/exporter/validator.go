package exporter

import (
	"encoding/hex"
	"strings"

	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
)

func (ces *ChainExporterService) getValidatorSetInfo(height int64) ([]*dtypes.ValidatorSetInfo, []*dtypes.MissInfo, []*dtypes.MissInfo, []*dtypes.MissDetailInfo, error) {
	// This is for syncing at the right height
	nextHeight := height + 1

	// Query the current block
	block, err := ces.RPCClient.Block(&height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Query the next block to access the commits
	nextBlock, err := ces.RPCClient.Block(&nextHeight)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Get validator set for the block
	validators, err := ces.RPCClient.Validators(&height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	validatorSetInfo := make([]*dtypes.ValidatorSetInfo, 0)
	missInfo := make([]*dtypes.MissInfo, 0)
	accumMissInfo := make([]*dtypes.MissInfo, 0)
	missDetailInfo := make([]*dtypes.MissDetailInfo, 0)

	for i, validator := range validators.Validators {
		// Insert genesis validators as an event_type of create_validator at height 1
		if validators.BlockHeight == 1 {
			tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
				IDValidator:          i + 1,
				Height:               block.Block.Height,
				Proposer:             validator.Address.String(),
				VotingPower:          float64(validator.VotingPower),
				NewVotingPowerAmount: float64(validator.VotingPower),
				NewVotingPowerDenom:  "kva",
				EventType:            "create_validator",
				Time:                 block.BlockMeta.Header.Time,
			}
			validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)
		}

		// Missing information
		if nextBlock.Block.LastCommit.Precommits[i] == nil {
			// Missing Detail information (save every single height)
			tempMissDetailInfo := &dtypes.MissDetailInfo{
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
			var prevMissInfo dtypes.MissInfo
			_ = ces.DB.Model(&prevMissInfo).
				Where("end_height = ? AND address = ?", endHeight-int64(1), validator.Address.String()).
				Order("end_height DESC").
				Select()

			if prevMissInfo.Address == "" {
				tempMissInfo := &dtypes.MissInfo{
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
				tempMissInfo := &dtypes.MissInfo{
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
	return validatorSetInfo, missInfo, accumMissInfo, missDetailInfo, nil
}

func (ces *ChainExporterService) getEvidenceInfo(height int64) ([]*dtypes.EvidenceInfo, error) {
	// This is for syncing at the right height
	nextHeight := height + 1

	// Query the current block
	block, err := ces.RPCClient.Block(&height)
	if err != nil {
		return nil, err
	}

	// Query the next block to access the commits
	nextBlock, err := ces.RPCClient.Block(&nextHeight)
	if err != nil {
		return nil, err
	}

	evidenceInfo := make([]*dtypes.EvidenceInfo, 0)

	// evidenceInfo
	for _, evidence := range nextBlock.Block.Evidence.Evidence {
		tempEvidenceInfo := &dtypes.EvidenceInfo{
			Address: strings.ToUpper(string(hex.EncodeToString(evidence.Address()))),
			Height:  evidence.Height(),
			Hash:    strings.ToUpper(string(hex.EncodeToString(evidence.Hash()))),
			Time:    block.BlockMeta.Header.Time,
		}
		evidenceInfo = append(evidenceInfo, tempEvidenceInfo)
	}

	return evidenceInfo, nil
}
