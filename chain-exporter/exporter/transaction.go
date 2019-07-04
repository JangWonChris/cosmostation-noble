package exporter

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"
	"github.com/tendermint/tendermint/crypto"
	resty "gopkg.in/resty.v1"
)

// Handling transaction data
func (ces *ChainExporterService) getTransactionInfo(height int64) ([]*dtypes.TransactionInfo, []*dtypes.VoteInfo, []*dtypes.DepositInfo, []*dtypes.ProposalInfo, error) {
	transactionInfo := make([]*dtypes.TransactionInfo, 0)
	voteInfo := make([]*dtypes.VoteInfo, 0)
	depositInfo := make([]*dtypes.DepositInfo, 0)
	proposalInfo := make([]*dtypes.ProposalInfo, 0)
	validatorSetInfo := make([]*dtypes.ValidatorSetInfo, 0)

	// Query the current block
	block, err := ces.RPCClient.Block(&height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for _, tx := range block.Block.Data.Txs {
		// Use tx codec to unmarshal binary length prefix
		var sdkTx sdk.Tx
		err := ces.Codec.UnmarshalBinaryLengthPrefixed([]byte(tx), &sdkTx)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		// TxHash
		txByte := crypto.Sha256(tx)
		txHash := hex.EncodeToString(txByte)
		txHash = strings.ToUpper(txHash)

		resp, _ := resty.R().Get(ces.Config.Node.LCDURL + "/txs/" + txHash)

		// Unmarshal general transaction format
		var generalTx dtypes.GeneralTx
		_ = json.Unmarshal(resp.Body(), &generalTx)

		// Check log to see if tx is success
		for j, log := range generalTx.Logs {
			if log.Success {
				switch generalTx.Tx.Value.Msg[j].Type {
				case "cosmos-sdk/MsgCreateValidator":
					var createValidatorTx dtypes.CreateValidatorMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &createValidatorTx)

					// 이렇게 넣는것도 문제가 발생!
					// 동일한 블록안에 create_validator 가 있을 경우 id_validator 를 체크하기가 힘들다
					// Query the highest height of id_validator
					var lastValidatorSetInfo dtypes.ValidatorSetInfo
					_ = ces.DB.Model(&lastValidatorSetInfo).
						Column("id_validator").
						Order("id_validator DESC").
						Limit(1).
						Select()

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(createValidatorTx.Value.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = float64(newVotingPowerAmount) / 1000000

					// Insert data
					tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          lastValidatorSetInfo.IDValidator + 1,
						Height:               height,
						Proposer:             utils.ConsensusPubkeyToProposer(createValidatorTx.Pubkey),
						VotingPower:          newVotingPowerAmount,
						NewVotingPowerAmount: newVotingPowerAmount,
						NewVotingPowerDenom:  createValidatorTx.Value.Denom,
						EventType:            "create_validator",
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

				case "cosmos-sdk/MsgDelegate":
					var delegateTx dtypes.DelegateMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &delegateTx)

					// Transaction Messeage
					var tempValidatorInfo dtypes.ValidatorInfo
					_ = ces.DB.Model(&tempValidatorInfo).
						Column("proposer").
						Where("operator_address = ?", delegateTx.ValidatorAddress).
						Limit(1).
						Select()

					// Query last id_validator
					var lastValidatorSetInfo dtypes.ValidatorSetInfo
					_ = ces.DB.Model(&lastValidatorSetInfo).
						Column("id_validator").
						Where("proposer = ?", tempValidatorInfo.Proposer).
						Order("id DESC").
						Limit(1).
						Select()

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(delegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = newVotingPowerAmount / 1000000

					// Current Voting Power
					var votingPower float64
					validators, _ := ces.RPCClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == tempValidatorInfo.Proposer {
							votingPower = float64(validator.VotingPower)
						}
					}

					// 동일한 블록에서 서로 다른 주소에서 동일한 검증인에게 위임한 트랜잭션이 있을 경우 현재 VotingPower는 같다. 기술적 한계 (Certus One 17번째 블록에 두번)
					// Insert data
					tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          lastValidatorSetInfo.IDValidator,
						Height:               height,
						Proposer:             tempValidatorInfo.Proposer,
						VotingPower:          votingPower + newVotingPowerAmount,
						EventType:            "delegate",
						NewVotingPowerAmount: newVotingPowerAmount,
						NewVotingPowerDenom:  delegateTx.Amount.Denom,
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

				case "cosmos-sdk/MsgUndelegate":
					var undelegateTx dtypes.UndelegateMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &undelegateTx)

					// Transaction Messeage
					var tempValidatorInfo dtypes.ValidatorInfo
					_ = ces.DB.Model(&tempValidatorInfo).
						Column("proposer").
						Where("operator_address = ?", undelegateTx.ValidatorAddress).
						Limit(1).
						Select()

					// Query last id_validator
					var lastValidatorSetInfo dtypes.ValidatorSetInfo
					_ = ces.DB.Model(&lastValidatorSetInfo).
						Column("id_validator", "voting_power").
						Where("proposer = ?", tempValidatorInfo.Proposer).
						Order("id DESC").
						Limit(1).
						Select()

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(undelegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = -newVotingPowerAmount / 1000000

					// Current Voting Power
					var votingPower float64
					validators, _ := ces.RPCClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == tempValidatorInfo.Proposer {
							votingPower = float64(validator.VotingPower)
						}
					}

					// Insert data
					tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          lastValidatorSetInfo.IDValidator,
						Height:               height,
						Proposer:             tempValidatorInfo.Proposer,
						VotingPower:          votingPower + newVotingPowerAmount,
						EventType:            "begin_unbonding",
						NewVotingPowerAmount: newVotingPowerAmount,
						NewVotingPowerDenom:  undelegateTx.Amount.Denom,
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

				case "cosmos-sdk/MsgBeginRedelegate":
					var redelegateTx dtypes.RedelegateMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &redelegateTx)

					/*
						Redelegate 당한 검증인은 -
						Redelegate 한 검증인은 +
					*/

					// Query validator_dst_address's proposer address
					var tempDstValidatorInfo dtypes.ValidatorInfo
					_ = ces.DB.Model(&tempDstValidatorInfo).
						Column("proposer").
						Where("operator_address = ?", redelegateTx.ValidatorDstAddress).
						Limit(1).
						Select()

					// Query validator_src_address's proposer address
					var tempSrcValidatorInfo dtypes.ValidatorInfo
					_ = ces.DB.Model(&tempSrcValidatorInfo).
						Column("proposer").
						Where("operator_address = ?", redelegateTx.ValidatorSrcAddress).
						Limit(1).
						Select()

					// Query last id_validator
					var lastDstValidatorSetInfo dtypes.ValidatorSetInfo
					_ = ces.DB.Model(&lastDstValidatorSetInfo).
						Column("id_validator", "voting_power").
						Where("proposer = ?", tempDstValidatorInfo.Proposer).
						Order("id DESC").
						Limit(1).
						Select()

					// Query last id_validator
					var lastSrcValidatorSetInfo dtypes.ValidatorSetInfo
					_ = ces.DB.Model(&lastSrcValidatorSetInfo).
						Column("id_validator", "voting_power").
						Where("proposer = ?", tempSrcValidatorInfo.Proposer).
						Order("id DESC").
						Limit(1).
						Select()

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(redelegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = newVotingPowerAmount / 1000000

					// Current Destination Validator's Voting Power
					var dstValidatorVotingPower float64
					validators, _ := ces.RPCClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == tempDstValidatorInfo.Proposer {
							dstValidatorVotingPower = float64(validator.VotingPower)
						}
					}

					// Insert ValidatorDstAddress data
					tempDstValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          lastDstValidatorSetInfo.IDValidator,
						Height:               height,
						Proposer:             tempDstValidatorInfo.Proposer,
						VotingPower:          dstValidatorVotingPower + newVotingPowerAmount,
						EventType:            "begin_redelegate",
						NewVotingPowerAmount: newVotingPowerAmount,
						NewVotingPowerDenom:  redelegateTx.Amount.Denom,
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempDstValidatorSetInfo)

					// Current Source Validator's Voting Power
					var srcValidatorVotingPower float64
					validators, _ = ces.RPCClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == tempSrcValidatorInfo.Proposer {
							srcValidatorVotingPower = float64(validator.VotingPower)
						}
					}

					// Insert ValidatorSrcAddress data
					tempSrcValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          lastSrcValidatorSetInfo.IDValidator,
						Height:               height,
						Proposer:             tempSrcValidatorInfo.Proposer,
						VotingPower:          srcValidatorVotingPower - newVotingPowerAmount,
						EventType:            "begin_redelegate",
						NewVotingPowerAmount: -newVotingPowerAmount,
						NewVotingPowerDenom:  redelegateTx.Amount.Denom,
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempSrcValidatorSetInfo)

				case "cosmos-sdk/MsgSubmitProposal":
					var submitTx dtypes.SubmitProposalMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &submitTx)

					// 멀티메시지 일 경우 Tags 번호가 달라져서 아래와 같이 key 값을 찾고 value를 넣어줘야 된다
					// 141050 블록높이: 7de25c478cf26eb6843c6a1b7a1cb550c8ab77ba9563a252677c059572bea6c3
					var proposalID int64
					for _, tag := range generalTx.Tags {
						if tag.Key == "proposal-id" {
							proposalID, _ = strconv.ParseInt(tag.Value, 10, 64)
						}
					}

					initialDepositAmount, _ := strconv.ParseInt(submitTx.InitialDeposit[0].Amount, 10, 64)
					initialDepositDenom := submitTx.InitialDeposit[0].Denom

					// Insert data
					tempProposalInfo := &dtypes.ProposalInfo{
						ID:                   proposalID,
						TxHash:               generalTx.TxHash,
						Proposer:             submitTx.Proposer,
						InitialDepositAmount: string(initialDepositAmount),
						InitialDepositDenom:  initialDepositDenom,
					}
					proposalInfo = append(proposalInfo, tempProposalInfo)

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
					gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

					// Insert data
					tempDepositInfo := &dtypes.DepositInfo{
						Height:     height,
						ProposalID: proposalID,
						Depositor:  submitTx.Proposer,
						Amount:     initialDepositAmount,
						Denom:      initialDepositDenom,
						TxHash:     generalTx.TxHash,
						GasWanted:  gasWanted,
						GasUsed:    gasUsed,
						Time:       block.BlockMeta.Header.Time,
					}
					depositInfo = append(depositInfo, tempDepositInfo)

				case "cosmos-sdk/MsgVote":
					var voteTx dtypes.VoteMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &voteTx)

					// Transaction Messeage
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					proposalID, _ := strconv.ParseInt(voteTx.ProposalID, 10, 64)
					gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
					gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)
					tempVoteInfo := &dtypes.VoteInfo{
						Height:     height,
						ProposalID: proposalID,
						Voter:      voteTx.Voter,
						Option:     voteTx.Option,
						TxHash:     generalTx.TxHash,
						GasWanted:  gasWanted,
						GasUsed:    gasUsed,
						Time:       block.BlockMeta.Header.Time,
					}
					voteInfo = append(voteInfo, tempVoteInfo)

				case "cosmos-sdk/MsgDeposit":
					var depositTx dtypes.DepositMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &depositTx)

					// Transaction Messeage
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					proposalID, _ := strconv.ParseInt(depositTx.ProposalID, 10, 64)
					amount, _ := strconv.ParseInt(depositTx.Amount[0].Amount, 10, 64)
					gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
					gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)
					tempDepositInfo := &dtypes.DepositInfo{
						Height:     height,
						ProposalID: proposalID,
						Depositor:  depositTx.Depositor,
						Amount:     amount,
						Denom:      depositTx.Amount[j].Denom,
						TxHash:     generalTx.TxHash,
						GasWanted:  gasWanted,
						GasUsed:    gasUsed,
						Time:       block.BlockMeta.Header.Time,
					}
					depositInfo = append(depositInfo, tempDepositInfo)
				default:
					continue
				}

				// PostgreSQL
				tempTransactionInfo := &dtypes.TransactionInfo{
					Height:  block.Block.Height,
					TxHash:  txHash,
					MsgType: generalTx.Tx.Value.Msg[j].Type,
					Time:    block.BlockMeta.Header.Time,
				}
				transactionInfo = append(transactionInfo, tempTransactionInfo)
			}
		}
	}

	return transactionInfo, voteInfo, depositInfo, proposalInfo, nil
}
