package exporter

import (
	"context"
	"strconv"

	"go.uber.org/zap"

	// mbl
	"github.com/cosmostation/mintscan-backend-library/db/schema"
	"github.com/cosmostation/mintscan-backend-library/types"

	// cosmos-sdk
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getPowerEventHistory returns voting power event history of validators by decoding transactions in a block.
func (ex *Exporter) getPowerEventHistory(block *tmctypes.ResultBlock, txResp []*sdktypes.TxResponse) ([]schema.PowerEventHistory, error) {
	powerEventHistory := make([]schema.PowerEventHistory, 0)

	if len(txResp) <= 0 {
		return powerEventHistory, nil
	}

	for _, tx := range txResp {
		// Other than code equals to 0, it is failed transaction.
		if tx.Code != 0 {
			continue
		}

		msgs := tx.GetTx().GetMsgs()

		for _, msg := range msgs {

			// zap.S().Infof("MsgType: %s | Hash: %s", rt.Body.Messages[0].Type, tx.TxHash)
			switch m := msg.(type) {
			// case staking.MsgCreateValidator:
			case *stakingtypes.MsgCreateValidator:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgCreateValidator := stdTx.Msgs[0].(staking.MsgCreateValidator)

				// Query the highest height of id_validator
				// TODO: Note that if two `create_validator` mesesages included in the same block then
				// id_validator may overlap. Needs to find other way to handle this.
				highestIDValidatorNum, _ := ex.db.QueryHighestValidatorID()

				newVotingPowerAmount, _ := strconv.ParseFloat(m.Value.Amount.String(), 64) // parseFloat from sdk.Dec.String()
				newVotingPowerAmount = float64(newVotingPowerAmount) / 1000000

				peh := &schema.PowerEventHistory{
					IDValidator:          highestIDValidatorNum + 1,
					Height:               tx.Height,
					Proposer:             m.Pubkey.String(), //jeonghwan : pubkey로부터 address 구하는 인터페이스가 string으로 변경 됨
					VotingPower:          newVotingPowerAmount,
					NewVotingPowerAmount: newVotingPowerAmount,
					NewVotingPowerDenom:  m.Value.Denom,
					MsgType:              types.StakingMsgCreateValidator,
					TxHash:               tx.TxHash,
					Timestamp:            block.Block.Header.Time,
				}

				powerEventHistory = append(powerEventHistory, *peh)

			// case staking.MsgDelegate:
			case *stakingtypes.MsgDelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type, tx.TxHash)

				// msgDelegate := stdTx.Msgs[0].(staking.MsgDelegate)

				// Query the validator's information.
				valInfo, _ := ex.db.QueryValidatorByAnyAddr(m.ValidatorAddress)

				// Query id_validator of lastly inserted data.
				validatorID, _ := ex.db.QueryValidatorID(valInfo.Proposer)

				newVotingPowerAmount, _ := strconv.ParseFloat(m.Amount.String(), 64) // parseFloat from sdk.Dec.String()
				newVotingPowerAmount = newVotingPowerAmount / 1000000

				// Get current voting power of the validator.
				// TODO: Note that if two MsgDelegate messages in one transaction, then
				// the validator's voting power may not be calculated correctly.
				var votingPower float64
				vals, _ := ex.client.RPC.GetValidatorsInHeight(tx.Height, 1, 150)
				for _, val := range vals.Validators {
					if val.Address.String() == valInfo.Proposer {
						votingPower = float64(val.VotingPower)
					}
				}

				peh := &schema.PowerEventHistory{
					IDValidator:          validatorID.IDValidator,
					Height:               tx.Height,
					Moniker:              valInfo.Moniker,
					OperatorAddress:      valInfo.OperatorAddress,
					Proposer:             valInfo.Proposer,
					VotingPower:          votingPower + newVotingPowerAmount,
					MsgType:              types.StakingMsgDelegate,
					NewVotingPowerAmount: newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					Timestamp:            block.Block.Header.Time,
				}

				powerEventHistory = append(powerEventHistory, *peh)

			// case staking.MsgUndelegate:
			case *stakingtypes.MsgUndelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgUndelegate := stdTx.Msgs[0].(staking.MsgUndelegate)

				// Query the validator's information.
				valInfo, _ := ex.db.QueryValidatorByAnyAddr(m.ValidatorAddress)

				// Query d_validator of lastly inserted data.
				validatorID, _ := ex.db.QueryValidatorID(valInfo.Proposer)

				newVotingPowerAmount, _ := strconv.ParseFloat(m.Amount.String(), 64) // parseFloat from sdk.Dec.String()
				newVotingPowerAmount = -newVotingPowerAmount / 1000000               // needs to be negative value

				// Get current voting power of the validator.
				var votingPower float64
				vals, _ := ex.client.RPC.GetValidatorsInHeight(tx.Height, types.DefaultQueryValidatorsPage, types.DefaultQueryValidatorsPerPage)
				for _, val := range vals.Validators {
					if val.Address.String() == valInfo.Proposer {
						votingPower = float64(val.VotingPower)
					}
				}

				peh := &schema.PowerEventHistory{
					IDValidator:          validatorID.IDValidator,
					Height:               tx.Height,
					Moniker:              valInfo.Moniker,
					OperatorAddress:      valInfo.OperatorAddress,
					Proposer:             valInfo.Proposer,
					VotingPower:          votingPower + newVotingPowerAmount,
					MsgType:              types.StakingMsgUndelegate,
					NewVotingPowerAmount: newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					Timestamp:            block.Block.Header.Time,
				}

				powerEventHistory = append(powerEventHistory, *peh)

			// case staking.MsgBeginRedelegate:
			case *stakingtypes.MsgBeginRedelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgBeginRedelegate := m.(staking.MsgBeginRedelegate)

				// Query validator_dst_address information.
				valDstInfo, _ := ex.db.QueryValidatorByAnyAddr(m.ValidatorDstAddress)
				dstpowerEventHistory, _ := ex.db.QueryValidatorID(valDstInfo.Proposer)

				// Query validator_src_address information.
				valSrcInfo, _ := ex.db.QueryValidatorByAnyAddr(m.ValidatorSrcAddress)
				srcpowerEventHistory, _ := ex.db.QueryValidatorID(valSrcInfo.Proposer)

				newVotingPowerAmount, _ := strconv.ParseFloat(m.Amount.String(), 64)
				newVotingPowerAmount = newVotingPowerAmount / 1000000

				// Get current destination validator's voting power.
				var dstValVotingPower float64
				vals, _ := ex.client.RPC.GetValidatorsInHeight(tx.Height, 1, 150)
				for _, val := range vals.Validators {
					if val.Address.String() == valDstInfo.Proposer {
						dstValVotingPower = float64(val.VotingPower)
					}
				}

				// Get current source validator's voting power.
				var srcValVotingPower float64
				vals, _ = ex.client.RPC.GetValidatorsInHeight(tx.Height, 1, 150)
				for _, val := range vals.Validators {
					if val.Address.String() == valSrcInfo.Proposer {
						srcValVotingPower = float64(val.VotingPower)
					}
				}

				dpeh := &schema.PowerEventHistory{
					IDValidator:          dstpowerEventHistory.IDValidator,
					Height:               tx.Height,
					Moniker:              valDstInfo.Moniker,
					OperatorAddress:      valDstInfo.OperatorAddress,
					Proposer:             valDstInfo.Proposer,
					VotingPower:          dstValVotingPower + newVotingPowerAmount, // Add
					MsgType:              types.StakingMsgBeginRedelegate,
					NewVotingPowerAmount: newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					Timestamp:            block.Block.Header.Time,
				}

				powerEventHistory = append(powerEventHistory, *dpeh)

				speh := &schema.PowerEventHistory{
					IDValidator:          srcpowerEventHistory.IDValidator,
					Height:               tx.Height,
					Moniker:              valSrcInfo.Moniker,
					OperatorAddress:      valSrcInfo.OperatorAddress,
					Proposer:             valSrcInfo.Proposer,
					VotingPower:          srcValVotingPower - newVotingPowerAmount, // Substract
					MsgType:              types.StakingMsgBeginRedelegate,
					NewVotingPowerAmount: -newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					Timestamp:            block.Block.Header.Time,
				}

				powerEventHistory = append(powerEventHistory, *speh)

			default:
				continue
			}
		}
	}

	return powerEventHistory, nil
}

// getValidatorsUptime has three slices
// missDetail gets every block
func (ex *Exporter) getValidatorsUptime(prevBlock *tmctypes.ResultBlock,
	block *tmctypes.ResultBlock, vals *tmctypes.ResultValidators) ([]schema.Miss, []schema.Miss, []schema.MissDetail, error) {

	miss := make([]schema.Miss, 0)
	accumMiss := make([]schema.Miss, 0)
	missDetail := make([]schema.MissDetail, 0)

	// MissDetailInfo saves every missing block of validators
	// while MissInfo saves ranges of missing blocks of validators.
	for i, val := range vals.Validators {
		// First block doesn't have any signatures from last commit
		if len(block.Block.LastCommit.Signatures) == 0 {
			break
		}

		// Note that it used to be block.Block.LastCommit.Precommits[i] == nil
		if block.Block.LastCommit.Signatures[i].Signature == nil {
			m := schema.MissDetail{
				Address:   val.Address.String(),
				Height:    prevBlock.Block.Header.Height,
				Proposer:  prevBlock.Block.Header.ProposerAddress.String(),
				Timestamp: prevBlock.Block.Header.Time,
			}

			missDetail = append(missDetail, m)

			// Set initial variables
			startHeight := prevBlock.Block.Header.Height
			endHeight := prevBlock.Block.Header.Height
			missingCount := int64(1)

			// Query if a validator hash missed previous block.
			prevMiss := ex.db.QueryMissingPreviousBlock(val.Address.String(), endHeight-int64(1))

			// Validator hasn't missed previous block.
			if prevMiss.Address == "" {
				m := schema.Miss{
					Address:      val.Address.String(),
					StartHeight:  startHeight,
					EndHeight:    endHeight,
					MissingCount: missingCount,
					StartTime:    prevBlock.Block.Header.Time,
					EndTime:      prevBlock.Block.Header.Time,
				}

				miss = append(miss, m)
			}

			// Validator has missed previous block.
			if prevMiss.Address != "" {
				m := schema.Miss{
					Address:      prevMiss.Address,
					StartHeight:  prevMiss.StartHeight,
					EndHeight:    prevMiss.EndHeight + int64(1),
					MissingCount: prevMiss.MissingCount + int64(1),
					StartTime:    prevMiss.StartTime,
					EndTime:      prevBlock.Block.Header.Time,
				}

				accumMiss = append(accumMiss, m)
			}
		}
	}

	return miss, accumMiss, missDetail, nil
}

// getEvidence provides evidence of malicious wrong-doing by validators.
// There is only DuplicateVoteEvidence. There is no downtime evidence.
func (ex *Exporter) getEvidence(block *tmctypes.ResultBlock) ([]schema.Evidence, error) {
	evidence := make([]schema.Evidence, 0)

	if block.Block.Evidence.Evidence != nil {
		for _, ev := range block.Block.Evidence.Evidence {
			e := schema.Evidence{
				// jeonghwan : ev.Address() are removed
				// Proposer:  strings.ToUpper(string(hex.EncodeToString(ev.Address()))),
				Proposer:  "",
				Height:    ev.Height(),
				Hash:      block.Block.Header.EvidenceHash.String(),
				Timestamp: block.Block.Header.Time,
			}

			evidence = append(evidence, e)
		}
	}

	return evidence, nil
}

// saveValidators parses all validators which are in three different status
// bonded, unbonding, unbonded and save them in database.
func (ex *Exporter) saveValidators() {
	ctx := context.Background()
	bondedVals, err := ex.client.GetValidatorsByStatus(ctx, stakingtypes.Bonded)
	if err != nil {
		zap.S().Errorf("failed to get bonded validators: %s", err)
		return
	}

	// Handle bonded validators sorted by highest tokens and insert or update them.
	err = ex.db.InsertOrUpdateValidators(bondedVals)
	if err != nil {
		zap.S().Errorf("failed to insert or update bonded validators: %s", err)
		return
	}

	unbondingVals, err := ex.client.GetValidatorsByStatus(ctx, stakingtypes.Unbonding)
	if err != nil {
		zap.S().Errorf("failed to get unbonding validators: %s", err)
		return
	}

	// Handle unbonding validators sorted by highest tokens and insert or update them.
	if len(unbondingVals) > 0 {
		highestBondedRank := ex.db.QueryHighestRankValidatorByStatus(types.BondedValidatorStatus)

		for i := range unbondingVals {
			unbondingVals[i].Rank = (highestBondedRank + 1 + i)
		}

		err := ex.db.InsertOrUpdateValidators(unbondingVals)
		if err != nil {
			zap.S().Errorf("failed to insert or update unbonding validators: %s", err)
			return
		}
	}

	unbondedVals, err := ex.client.GetValidatorsByStatus(ctx, stakingtypes.Unbonded)
	if err != nil {
		zap.S().Errorf("failed to get unbonded validators: %s", err)
		return
	}

	// Handle unbonded validators sorted by highest tokens and insert or update them.
	if len(unbondedVals) > 0 {
		unbondingHighestRank := ex.db.QueryHighestRankValidatorByStatus(types.UnbondingValidatorStatus)

		if unbondingHighestRank == 0 {
			unbondingHighestRank = ex.db.QueryHighestRankValidatorByStatus(types.BondedValidatorStatus)
		}

		for i := range unbondedVals {
			unbondedVals[i].Rank = (unbondingHighestRank + 1 + i)
		}

		err := ex.db.InsertOrUpdateValidators(unbondedVals)
		if err != nil {
			zap.S().Errorf("failed to insert or update unbonded validators: %s", err)
			return
		}
	}
}

// saveValidatorsIdentities saves all KeyBase URLs of validators
func (ex *Exporter) saveValidatorsIdentities() {
	vals, _ := ex.db.QueryValidators()

	result, err := ex.client.GetValidatorsIdentities(vals)
	if err != nil {
		zap.S().Errorf("failed to get validator identities: %s", err)
		return
	}

	if len(result) > 0 {
		ex.db.UpdateValidatorsKeyBaseURL(result)
		return
	}
}
