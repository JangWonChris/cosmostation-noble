package exporter

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "gopkg.in/resty.v1"
)

// getGenesisValidatorSet returns validator set in genesis
func (ex *Exporter) getGenesisValidatorSet(block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) ([]*schema.PowerEventHistory, error) {
	genesisValsSet := make([]*schema.PowerEventHistory, 0)

	// Get validator set for the height
	for i, validator := range vals.Validators {
		if vals.BlockHeight == 1 {
			tempValsSet := &schema.PowerEventHistory{
				IDValidator:          i + 1,
				Height:               block.Block.Height,
				Proposer:             validator.Address.String(),
				VotingPower:          float64(validator.VotingPower),
				NewVotingPowerAmount: float64(validator.VotingPower),
				NewVotingPowerDenom:  types.Denom,
				EventType:            types.EventTypeMsgCreateValidator,
				Time:                 block.BlockMeta.Header.Time,
			}
			genesisValsSet = append(genesisValsSet, tempValsSet)
		}
	}

	return genesisValsSet, nil
}

// getPowerEventHistory provides validator set information in every block
func (ex *Exporter) getPowerEventHistory(prevBlock *tmctypes.ResultBlock, block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) ([]*schema.Miss, []*schema.Miss, []*schema.MissDetail, error) {

	miss := make([]*schema.Miss, 0)
	accumMiss := make([]*schema.Miss, 0)
	missDetail := make([]*schema.MissDetail, 0)

	// Get validator set for the height
	for i, validator := range vals.Validators {
		// MissDetailInfo saves every missing information of validators
		// MissInfo saves ranges of missing information of validators
		// check if a validator misses previous block
		if block.Block.LastCommit.Precommits[i] == nil {
			tempMissDetail := &schema.MissDetail{
				Height:   prevBlock.BlockMeta.Header.Height,
				Address:  validator.Address.String(),
				Proposer: prevBlock.BlockMeta.Header.ProposerAddress.String(),
				Alerted:  false,
				Time:     prevBlock.BlockMeta.Header.Time,
			}
			missDetail = append(missDetail, tempMissDetail)

			// Initial variables
			startHeight := prevBlock.BlockMeta.Header.Height
			endHeight := prevBlock.BlockMeta.Header.Height
			missingCount := int64(1)

			// Query to check if a validator missed previous block
			var prevMiss schema.Miss
			_ = ex.db.Model(&prevMiss).
				Where("end_height = ? AND address = ?", endHeight-int64(1), validator.Address.String()).
				Order("end_height DESC").
				Select()

			if prevMiss.Address == "" {
				tempMiss := &schema.Miss{
					Address:      validator.Address.String(),
					StartHeight:  startHeight,
					EndHeight:    endHeight,
					MissingCount: missingCount,
					StartTime:    prevBlock.BlockMeta.Header.Time,
					EndTime:      prevBlock.BlockMeta.Header.Time,
					Alerted:      false,
				}
				miss = append(miss, tempMiss)
			} else {
				tempMiss := &schema.Miss{
					Address:      prevMiss.Address,
					StartHeight:  prevMiss.StartHeight,
					EndHeight:    prevMiss.EndHeight + int64(1),
					MissingCount: prevMiss.MissingCount + int64(1),
					StartTime:    prevMiss.StartTime,
					EndTime:      prevBlock.BlockMeta.Header.Time,
					Alerted:      false,
				}
				accumMiss = append(accumMiss, tempMiss)
			}
			continue
		}
	}
	return miss, accumMiss, missDetail, nil
}

// getEvidence provides evidence (slashing) information
func (ex *Exporter) getEvidence(block *tmctypes.ResultBlock, nextBlock *tmctypes.ResultBlock) ([]*schema.Evidence, error) {
	// In cosmoshub2: 848187 = 1C4DB67E79B5BB30663B04245E064E6180EC6EA304EE83A7A879B04544A2EAD0
	// in evidence, there is only DuplicateVoteEvidence. There is no downtime evidence.
	evidence := make([]*schema.Evidence, 0)
	if nextBlock.Block.Evidence.Evidence != nil {
		for _, evi := range nextBlock.Block.Evidence.Evidence {
			tempEvidence := &schema.Evidence{
				Proposer: strings.ToUpper(string(hex.EncodeToString(evi.Address()))),
				Height:   evi.Height(),
				Hash:     nextBlock.BlockMeta.Header.EvidenceHash.String(),
				Time:     block.BlockMeta.Header.Time,
			}
			evidence = append(evidence, tempEvidence)
		}
	}

	return evidence, nil
}

// SaveValidatorKeyBase saves keybase urls for every validator
func (ex *Exporter) SaveValidatorKeyBase() error {
	resultValidators := make([]*schema.Validator, 0)

	// query validators
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

			tempValidator := &schema.Validator{
				ID:         validator.ID,
				KeybaseURL: keybaseURL,
			}
			resultValidators = append(resultValidators, tempValidator)
		}
	}

	if len(resultValidators) > 0 {
		for _, validator := range resultValidators {
			result, err := ex.db.UpdateKeyBase(validator.ID, validator.KeybaseURL)
			if !result {
				log.Printf("failed to update KeyBase URL: %v \n", err)
			}
		}
	}

	return nil
}
