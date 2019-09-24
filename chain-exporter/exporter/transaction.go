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

// getTransactionInfo provides information about each transaction in every block
func (ces *ChainExporterService) getTransactionInfo(height int64) ([]*dtypes.TransactionInfo, []*dtypes.VoteInfo,
	[]*dtypes.DepositInfo, []*dtypes.ProposalInfo, []*dtypes.ValidatorSetInfo, error) {

	transactionInfo := make([]*dtypes.TransactionInfo, 0)
	voteInfo := make([]*dtypes.VoteInfo, 0)
	depositInfo := make([]*dtypes.DepositInfo, 0)
	proposalInfo := make([]*dtypes.ProposalInfo, 0)
	validatorSetInfo := make([]*dtypes.ValidatorSetInfo, 0)

	// query current block
	block, err := ces.rpcClient.Block(&height)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if len(block.Block.Data.Txs) > 0 {
		for _, tx := range block.Block.Data.Txs {
			// use tx codec to unmarshal binary length prefix
			var sdkTx sdk.Tx
			_ = ces.codec.UnmarshalBinaryLengthPrefixed([]byte(tx), &sdkTx)

			// transaction hash
			txByte := crypto.Sha256(tx)
			txHash := hex.EncodeToString(txByte)
			txHash = strings.ToUpper(txHash)

			// unmarshal general transaction format
			var generalTx dtypes.GeneralTx
			resp, _ := resty.R().Get(ces.config.Node.LCDURL + "/txs/" + txHash)
			err := json.Unmarshal(resp.Body(), &generalTx)
			if err != nil {
				fmt.Printf("unmarshal generalTx error - %v\n", err)
			}

			// save all txs in PostgreSQL if it is success or fail
			if len(generalTx.Tx.Value.Msg) == 1 {
				tempTransactionInfo := &dtypes.TransactionInfo{
					Height:  block.Block.Height,
					TxHash:  txHash,
					MsgType: generalTx.Tx.Value.Msg[0].Type,
					Memo:    generalTx.Tx.Value.Memo,
					Time:    block.BlockMeta.Header.Time,
				}
				transactionInfo = append(transactionInfo, tempTransactionInfo)
			} else {
				tempTransactionInfo := &dtypes.TransactionInfo{
					Height:  block.Block.Height,
					TxHash:  txHash,
					MsgType: dtypes.MultiMsg,
					Memo:    generalTx.Tx.Value.Memo,
					Time:    block.BlockMeta.Header.Time,
				}
				transactionInfo = append(transactionInfo, tempTransactionInfo)
			}

			// check log to see if tx is success
			for j, log := range generalTx.Logs {
				if log.Success {
					switch generalTx.Tx.Value.Msg[j].Type {
					case "cosmos-sdk/MsgCreateValidator":
						var msgCreateValidator dtypes.MsgCreateValidator
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgCreateValidator)

						/*
							[기술적 한계] > 동일한 블록안에 create_validator 메시지가 2개 이상 있을 경우 마지막으로 저장된 id_validator를 가져오면 겹친다.
						*/

						// query the highest height of id_validator
						highestIDValidatorNum, _ := utils.QueryHighestIDValidatorNum(ces.db)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgCreateValidator.Value.Amount.String(), 64) // parseFloat from sdk.Dec.String()
						newVotingPowerAmount = float64(newVotingPowerAmount) / 1000000

						tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
							IDValidator:          highestIDValidatorNum + 1,
							Height:               height,
							Proposer:             utils.ConsAddrFromConsPubkey(msgCreateValidator.Pubkey), // new validator's proposer address needs to be converted
							VotingPower:          newVotingPowerAmount,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgCreateValidator.Value.Denom,
							EventType:            dtypes.EventTypeMsgCreateValidator,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

					case "cosmos-sdk/MsgDelegate":
						var msgDelegate dtypes.MsgDelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgDelegate)

						// query validator information fro validator_infos table
						validatorInfo, _ := utils.QueryValidatorInfo(ces.db, msgDelegate.ValidatorAddress)

						// query to get id_validator of lastly inserted data
						idValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorInfo.Proposer)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgDelegate.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
						newVotingPowerAmount = newVotingPowerAmount / 1000000

						// current voting power of a validator
						var votingPower float64
						validators, _ := ces.rpcClient.Validators(&height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorInfo.Proposer {
								votingPower = float64(validator.VotingPower)
							}
						}

						/*
							[기술적 한계] - Certus One 17번째 블록에 두번 - cosmoshub-1
										동일한 블록에서 서로 다른 주소에서 동일한 검증인에게 위임한 트랜잭션이 있을 경우 현재 VotingPower는 같다.
						*/

						tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
							IDValidator:          idValidatorSetInfo.IDValidator,
							Height:               height,
							Moniker:              validatorInfo.Moniker,
							OperatorAddress:      validatorInfo.OperatorAddress,
							Proposer:             validatorInfo.Proposer,
							VotingPower:          votingPower + newVotingPowerAmount,
							EventType:            dtypes.EventTypeMsgDelegate,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgDelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

					case "cosmos-sdk/MsgUndelegate":
						var msgUndelegate dtypes.MsgUndelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgUndelegate)

						// query validator info
						validatorInfo, _ := utils.QueryValidatorInfo(ces.db, msgUndelegate.ValidatorAddress)

						// query to get id_validator of lastly inserted data
						idValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorInfo.Proposer)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgUndelegate.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
						newVotingPowerAmount = -newVotingPowerAmount / 1000000

						// current voting power of a validator
						var votingPower float64
						validators, _ := ces.rpcClient.Validators(&height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorInfo.Proposer {
								votingPower = float64(validator.VotingPower)
							}
						}

						// substract the undelegated amount from the validator
						tempValidatorSetInfo := &dtypes.ValidatorSetInfo{
							IDValidator:          idValidatorSetInfo.IDValidator,
							Height:               height,
							Moniker:              validatorInfo.Moniker,
							OperatorAddress:      validatorInfo.OperatorAddress,
							Proposer:             validatorInfo.Proposer,
							VotingPower:          votingPower + newVotingPowerAmount,
							EventType:            dtypes.EventTypeMsgUndelegate,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgUndelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

					case "cosmos-sdk/MsgBeginRedelegate":
						/*
							[Note]
								+ for ValidatorDstAddress
								- for ValidatorSrcAddress
						*/

						var msgBeginRedelegate dtypes.MsgBeginRedelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgBeginRedelegate)

						// query validator_dst_address info
						validatorDstInfo, _ := utils.QueryValidatorInfo(ces.db, msgBeginRedelegate.ValidatorDstAddress)
						dstValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorDstInfo.Proposer)

						// query validator_src_address info
						validatorSrcInfo, _ := utils.QueryValidatorInfo(ces.db, msgBeginRedelegate.ValidatorSrcAddress)
						srcValidatorSetInfo, _ := utils.QueryIDValidatorSetInfo(ces.db, validatorSrcInfo.Proposer)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgBeginRedelegate.Amount.Amount.String(), 64)
						newVotingPowerAmount = newVotingPowerAmount / 1000000

						// current destination validator's voting power
						var dstValidatorVotingPower float64
						validators, _ := ces.rpcClient.Validators(&height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorDstInfo.Proposer {
								dstValidatorVotingPower = float64(validator.VotingPower)
							}
						}

						// current source validator's voting power
						var srcValidatorVotingPower float64
						validators, _ = ces.rpcClient.Validators(&height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorSrcInfo.Proposer {
								srcValidatorVotingPower = float64(validator.VotingPower)
							}
						}

						// add the redelegated amount to validator_dst_address
						tempDstValidatorSetInfo := &dtypes.ValidatorSetInfo{
							IDValidator:          dstValidatorSetInfo.IDValidator,
							Height:               height,
							Moniker:              validatorDstInfo.Moniker,
							OperatorAddress:      validatorDstInfo.OperatorAddress,
							Proposer:             validatorDstInfo.Proposer,
							VotingPower:          dstValidatorVotingPower + newVotingPowerAmount,
							EventType:            dtypes.EventTypeMsgBeginRedelegate,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgBeginRedelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						validatorSetInfo = append(validatorSetInfo, tempDstValidatorSetInfo)

						// substract the redelegated amount from validator_src_address
						tempSrcValidatorSetInfo := &dtypes.ValidatorSetInfo{
							IDValidator:          srcValidatorSetInfo.IDValidator,
							Height:               height,
							Moniker:              validatorSrcInfo.Moniker,
							OperatorAddress:      validatorSrcInfo.OperatorAddress,
							Proposer:             validatorSrcInfo.Proposer,
							VotingPower:          srcValidatorVotingPower - newVotingPowerAmount,
							EventType:            dtypes.EventTypeMsgBeginRedelegate,
							NewVotingPowerAmount: -newVotingPowerAmount,
							NewVotingPowerDenom:  msgBeginRedelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						validatorSetInfo = append(validatorSetInfo, tempSrcValidatorSetInfo)

					case "cosmos-sdk/MsgSubmitProposal":
						var msgSubmitProposal dtypes.MsgSubmitProposal
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgSubmitProposal)

						// take care of multi-msg
						var proposalID int64
						for _, event := range generalTx.Events {
							for _, attribute := range event.Attributes {
								if attribute.Key == "proposal_id" {
									proposalID, _ = strconv.ParseInt(attribute.Value, 10, 64)
								}
							}
						}

						var initialDepositAmount string
						var initialDepositDenom string

						if len(msgSubmitProposal.InitialDeposit) > 0 {
							initialDepositAmount = msgSubmitProposal.InitialDeposit[0].Amount
							initialDepositDenom = msgSubmitProposal.InitialDeposit[0].Denom
						}

						tempProposalInfo := &dtypes.ProposalInfo{
							ID:                   proposalID,
							TxHash:               generalTx.TxHash,
							Proposer:             msgSubmitProposal.Proposer,
							InitialDepositAmount: initialDepositAmount,
							InitialDepositDenom:  initialDepositDenom,
						}
						proposalInfo = append(proposalInfo, tempProposalInfo)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempDepositInfo := &dtypes.DepositInfo{
							Height:     height,
							ProposalID: proposalID,
							Depositor:  msgSubmitProposal.Proposer,
							Amount:     initialDepositAmount,
							Denom:      initialDepositDenom,
							TxHash:     generalTx.TxHash,
							GasWanted:  gasWanted,
							GasUsed:    gasUsed,
							Time:       block.BlockMeta.Header.Time,
						}
						depositInfo = append(depositInfo, tempDepositInfo)

					case "cosmos-sdk/MsgVote":
						var msgVote dtypes.MsgVote
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgVote)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						proposalID, _ := strconv.ParseInt(msgVote.ProposalID, 10, 64)
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempVoteInfo := &dtypes.VoteInfo{
							Height:     height,
							ProposalID: proposalID,
							Voter:      msgVote.Voter,
							Option:     msgVote.Option,
							TxHash:     generalTx.TxHash,
							GasWanted:  gasWanted,
							GasUsed:    gasUsed,
							Time:       block.BlockMeta.Header.Time,
						}
						voteInfo = append(voteInfo, tempVoteInfo)

					case "cosmos-sdk/MsgDeposit":
						var msgDeposit dtypes.MsgDeposit
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgDeposit)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						proposalID, _ := strconv.ParseInt(msgDeposit.ProposalID, 10, 64)
						amount := msgDeposit.Amount[0].Amount
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempDepositInfo := &dtypes.DepositInfo{
							Height:     height,
							ProposalID: proposalID,
							Depositor:  msgDeposit.Depositor,
							Amount:     amount,
							Denom:      msgDeposit.Amount[j].Denom,
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
			}
		}
	}

	return transactionInfo, voteInfo, depositInfo, proposalInfo, validatorSetInfo, nil
}
