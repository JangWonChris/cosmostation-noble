package exporter

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dtypes "github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"
	"github.com/tendermint/tendermint/crypto"
	resty "gopkg.in/resty.v1"
)

// Handling transaction data
func (ces *ChainExporterService) getTransactionInfo(height int64) ([]*dtypes.TransactionInfo, []*dtypes.VoteInfo,
	[]*dtypes.DepositInfo, []*dtypes.ProposalInfo, []*dtypes.ValidatorSetInfo, error) {
	transactionInfo := make([]*dtypes.TransactionInfo, 0)
	voteInfo := make([]*dtypes.VoteInfo, 0)
	depositInfo := make([]*dtypes.DepositInfo, 0)
	proposalInfo := make([]*dtypes.ProposalInfo, 0)
	validatorSetInfo := make([]*dtypes.ValidatorSetInfo, 0)

	// Query the current block
	block, err := ces.rpcClient.Block(&height)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	for _, tx := range block.Block.Data.Txs {
		// Use tx codec to unmarshal binary length prefix
		var sdkTx sdk.Tx
		_ = ces.codec.UnmarshalBinaryLengthPrefixed([]byte(tx), &sdkTx)

		// Tx hash
		txByte := crypto.Sha256(tx)
		txHash := hex.EncodeToString(txByte)
		txHash = strings.ToUpper(txHash)

		// Unmarshal general transaction format
		var generalTx dtypes.GeneralTx
		resp, _ := resty.R().Get(ces.config.Node.LCDURL + "/txs/" + txHash)
		_ = json.Unmarshal(resp.Body(), &generalTx)

		// Check log to see if tx is success
		for j, log := range generalTx.Logs {
			if log.Success {
				switch generalTx.Tx.Value.Msg[j].Type {
				case "cosmos-sdk/MsgCreateValidator":
					var createValidatorTx dtypes.CreateValidatorMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &createValidatorTx)

					// [기술적 한계] > 동일한 블록안에 create_validator 메시지가 2개 이상 있을 경우 마지막으로 저장된 id_validator를 가져오면 겹친다.

					// Query the highest height of id_validator
					highestIDValidatorNum, _ := utils.QueryHighestIDValidatorNum(ces.db)

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(createValidatorTx.Value.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = float64(newVotingPowerAmount) / 1000000

					// Insert data
					tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          highestIDValidatorNum + 1,
						Height:               height,
						Proposer:             utils.ConsensusPubkeyToProposer(createValidatorTx.Pubkey), // New validator's proposer address needs to be converted
						VotingPower:          newVotingPowerAmount,
						NewVotingPowerAmount: newVotingPowerAmount,
						NewVotingPowerDenom:  createValidatorTx.Value.Denom,
						EventType:            dtypes.EventTypeMsgCreateValidator,
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

				case "cosmos-sdk/MsgDelegate":
					var delegateTx dtypes.DelegateMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &delegateTx)

					// Query validator info
					validatorInfo, _ := utils.QueryValidatorInfo(ces.db, delegateTx.ValidatorAddress)

					// Query to get id_validator of lastly inserted data
					idValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorInfo.Proposer)

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(delegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = newVotingPowerAmount / 1000000

					// Current Voting Power
					var votingPower float64
					validators, _ := ces.rpcClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == validatorInfo.Proposer {
							votingPower = float64(validator.VotingPower)
						}
					}

					// 기술적 한계 (Certus One 17번째 블록에 두번 - cosmoshub-1)
					// 동일한 블록에서 서로 다른 주소에서 동일한 검증인에게 위임한 트랜잭션이 있을 경우 현재 VotingPower는 같다.
					// Insert data
					tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          idValidatorSetInfo.IDValidator,
						Height:               height,
						Moniker:              validatorInfo.Moniker,
						OperatorAddress:      validatorInfo.OperatorAddress,
						Proposer:             validatorInfo.Proposer,
						VotingPower:          votingPower + newVotingPowerAmount,
						EventType:            dtypes.EventTypeMsgDelegate,
						NewVotingPowerAmount: newVotingPowerAmount,
						NewVotingPowerDenom:  delegateTx.Amount.Denom,
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

				case "cosmos-sdk/MsgUndelegate":
					var undelegateTx dtypes.UndelegateMsgValueTx
					_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &undelegateTx)

					// Query validator info
					validatorInfo, _ := utils.QueryValidatorInfo(ces.db, undelegateTx.ValidatorAddress)

					// Query to get id_validator of lastly inserted data
					idValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorInfo.Proposer)

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(undelegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = -newVotingPowerAmount / 1000000

					// Current Voting Power
					var votingPower float64
					validators, _ := ces.rpcClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == validatorInfo.Proposer {
							votingPower = float64(validator.VotingPower)
						}
					}

					// Substract the undelegated amount from the validator
					tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          idValidatorSetInfo.IDValidator,
						Height:               height,
						Moniker:              validatorInfo.Moniker,
						OperatorAddress:      validatorInfo.OperatorAddress,
						Proposer:             block.BlockMeta.Header.ProposerAddress.String(),
						VotingPower:          votingPower + newVotingPowerAmount,
						EventType:            dtypes.EventTypeMsgUndelegate,
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
						Note : + for ValidatorDstAddress | - for ValidatorSrcAddress
					*/

					// Query validator_dst_address info
					validatorDstInfo, _ := utils.QueryValidatorInfo(ces.db, redelegateTx.ValidatorDstAddress)
					dstValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorDstInfo.Proposer)

					// Query validator_src_address info
					validatorSrcInfo, _ := utils.QueryValidatorInfo(ces.db, redelegateTx.ValidatorSrcAddress)
					srcValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorSrcInfo.Proposer)

					// Conversion
					height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
					newVotingPowerAmount, _ := strconv.ParseFloat(redelegateTx.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
					newVotingPowerAmount = newVotingPowerAmount / 1000000

					// Current destination validator's voting power
					var dstValidatorVotingPower float64
					validators, _ := ces.rpcClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == validatorDstInfo.Proposer {
							dstValidatorVotingPower = float64(validator.VotingPower)
						}
					}

					// Current source validator's voting power
					var srcValidatorVotingPower float64
					validators, _ = ces.rpcClient.Validators(&height)
					for _, validator := range validators.Validators {
						if validator.Address.String() == validatorSrcInfo.Proposer {
							srcValidatorVotingPower = float64(validator.VotingPower)
						}
					}

					// Add the redelegated amount to validator_dst_address
					tempDstValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          dstValidatorSetInfo.IDValidator,
						Height:               height,
						Moniker:              validatorDstInfo.Moniker,
						OperatorAddress:      validatorDstInfo.OperatorAddress,
						Proposer:             validatorDstInfo.Proposer,
						VotingPower:          dstValidatorVotingPower + newVotingPowerAmount,
						EventType:            dtypes.EventTypeMsgBeginRedelegate,
						NewVotingPowerAmount: newVotingPowerAmount,
						NewVotingPowerDenom:  redelegateTx.Amount.Denom,
						TxHash:               generalTx.TxHash,
						Time:                 block.BlockMeta.Header.Time,
					}
					validatorSetInfo = append(validatorSetInfo, tempDstValidatorSetInfo)

					// Substract the redelegated amount from validator_src_address
					tempSrcValidatorSetInfo := &dtypes.ValidatorSetInfo{
						IDValidator:          srcValidatorSetInfo.IDValidator,
						Height:               height,
						Moniker:              validatorSrcInfo.Moniker,
						OperatorAddress:      validatorSrcInfo.OperatorAddress,
						Proposer:             validatorSrcInfo.Proposer,
						VotingPower:          srcValidatorVotingPower - newVotingPowerAmount,
						EventType:            dtypes.EventTypeMsgBeginRedelegate,
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

					initialDepositAmount, _ := strconv.ParseFloat(submitTx.InitialDeposit[0].Amount, 64)
					depositAmount := fmt.Sprintf("%f", initialDepositAmount)
					initialDepositDenom := submitTx.InitialDeposit[0].Denom

					// Insert data
					tempProposalInfo := &dtypes.ProposalInfo{
						ID:                   proposalID,
						TxHash:               generalTx.TxHash,
						Proposer:             submitTx.Proposer,
						InitialDepositAmount: depositAmount,
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
						Amount:     depositAmount,
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
						Amount:     string(amount),
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
			}

			// PostgreSQL : save all txs whether it is success or fail
			tempTransactionInfo := &dtypes.TransactionInfo{
				Height:  block.Block.Height,
				TxHash:  txHash,
				MsgType: generalTx.Tx.Value.Msg[j].Type,
				Time:    block.BlockMeta.Header.Time,
			}
			transactionInfo = append(transactionInfo, tempTransactionInfo)
		}
	}

	return transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo, nil
}
