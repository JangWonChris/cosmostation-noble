package exporter

import (
	"context"
	"time"

	"go.uber.org/zap"

	// mbl
	"github.com/cosmostation/cosmostation-cosmos/chain-config/custom"
	"github.com/cosmostation/mintscan-backend-library/db/schema"
	"github.com/cosmostation/mintscan-backend-library/types"

	// cosmos-sdk
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getPowerEventHistory returns voting power event history of validators by decoding transactions in a block.
func (ex *Exporter) getPowerEventHistoryNew( /*block *tmctypes.ResultBlock,*/ txResp []*sdktypes.TxResponse) ([]schema.PowerEventHistory, error) {
	/*
		구현 방향 정리 :
		1. validator 테이블에 존재하는 데이터(검증인 정보)를 중복으로 저장 할 필요가 없다.
			필요하면 조인 연산을 통해 결과를 만들어 내도록 하고, exporter에서는 조인에 필요한 키(validator operator address)만 취한다.
			이렇게 되면, chain-exporter에서 power event 저장 시 사용되는 로직 일부를 제거할 수 있다.

		2. validator 테이블의 ID를 관계를 엮어서 가져오는게 아닌, 조회를 통해 넣고 있다.
			이 역시 validator operator address를 이용하면, 외래 키로 이용이 가능하기 때문에 이 컬럼 역시 제거한다.

		3. 특정 높이의 consensus power를 power_event_history 테이블에 저장하지 않는다.
			따라서, 변화량만 저장한다.
			따라서, validator로부터 전체 리스트를 가져올 필요가 없다.
			따라서, transaction-account 테이블로부터 이 데이터를 만들어 낼 수 있다.(이렇게 하면, 프론트 공수가 추가적으로 들어간다.)

		4. (3)의 결정에 따라, 블록 별 검증인의 consensus 변화 추이를 계산 할 전체 변화량을 저장 할 필요가 있다.
			4-1. 노드에 특정 높이의 검증인 집합을 요청하고, 그 값을 사용한다. (이전 체인 데이터의 재생산이 어려움)
			4-2. staking tx로부터 검증인의 보팅 파워 변화를 계산하고 그 결과를 테이블에 저장한다. (이전 체인 히스토리도 모두 지원 할 수 있음, 데이터의 정합성 검증 추가 필요)
	*/
	powerEventHistory := make([]schema.PowerEventHistory, 0)

	if len(txResp) <= 0 {
		return powerEventHistory, nil
	}

	for _, tx := range txResp {
		if tx.Code != 0 {
			// Code != 0 이면, 성공한 tx가 아니므로 무시한다.
			continue
		}

		timestamp, _ := time.Parse(time.RFC3339, tx.Timestamp) // 임시
		msgs := tx.GetTx().GetMsgs()

		for _, msg := range msgs {

			switch m := msg.(type) {
			case *stakingtypes.MsgCreateValidator:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				newVotingPowerAmount := float64(m.Value.Amount.Quo(custom.PowerReduction).Int64())

				peh := &schema.PowerEventHistory{
					Height:               tx.Height,
					OperatorAddress:      m.ValidatorAddress,
					MsgType:              types.StakingMsgCreateValidator,
					NewVotingPowerAmount: newVotingPowerAmount,
					NewVotingPowerDenom:  m.Value.Denom,
					TxHash:               tx.TxHash,
					// Timestamp:            block.Block.Header.Time,
					Timestamp: timestamp,
				}

				powerEventHistory = append(powerEventHistory, *peh)

			case *stakingtypes.MsgDelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				newVotingPowerAmount := float64(m.Amount.Amount.Quo(custom.PowerReduction).Int64())

				peh := &schema.PowerEventHistory{
					Height:               tx.Height,
					OperatorAddress:      m.ValidatorAddress,
					MsgType:              types.StakingMsgDelegate,
					NewVotingPowerAmount: newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					// Timestamp:            block.Block.Header.Time,
					Timestamp: timestamp,
				}

				powerEventHistory = append(powerEventHistory, *peh)

			case *stakingtypes.MsgUndelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				newVotingPowerAmount := float64(m.Amount.Amount.Quo(custom.PowerReduction).Int64())

				peh := &schema.PowerEventHistory{
					Height:               tx.Height,
					OperatorAddress:      m.ValidatorAddress,
					MsgType:              types.StakingMsgUndelegate,
					NewVotingPowerAmount: -newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					// Timestamp:            block.Block.Header.Time,
					Timestamp: timestamp,
				}

				powerEventHistory = append(powerEventHistory, *peh)

			case *stakingtypes.MsgBeginRedelegate:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				newVotingPowerAmount := float64(m.Amount.Amount.Quo(custom.PowerReduction).Int64())

				// destination (add power)
				dpeh := &schema.PowerEventHistory{
					Height:               tx.Height,
					OperatorAddress:      m.ValidatorDstAddress,
					MsgType:              types.StakingMsgBeginRedelegate,
					NewVotingPowerAmount: newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					// Timestamp:            block.Block.Header.Time,
					Timestamp: timestamp,
				}

				powerEventHistory = append(powerEventHistory, *dpeh)

				//source (subtract power)
				speh := &schema.PowerEventHistory{
					Height:               tx.Height,
					OperatorAddress:      m.ValidatorSrcAddress,
					MsgType:              types.StakingMsgBeginRedelegate,
					NewVotingPowerAmount: -newVotingPowerAmount,
					NewVotingPowerDenom:  m.Amount.Denom,
					TxHash:               tx.TxHash,
					// Timestamp:            block.Block.Header.Time,
					Timestamp: timestamp,
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
