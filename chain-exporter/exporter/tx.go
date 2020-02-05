package exporter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/notification"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	resty "gopkg.in/resty.v1"
)

// getTransactions returns transactions in a block
func (ex *Exporter) getTransactions(block *tmctypes.ResultBlock) ([]*schema.Vote, []*schema.Deposit, []*schema.Proposal, []*schema.PowerEventHistory, error) {
	vote := make([]*schema.Vote, 0)
	deposit := make([]*schema.Deposit, 0)
	proposal := make([]*schema.Proposal, 0)
	powerEventHistory := make([]*schema.PowerEventHistory, 0)

	if len(block.Block.Data.Txs) > 0 {
		for _, tmTx := range block.Block.Data.Txs {
			var sdkTx sdk.Tx
			_ = ex.cdc.UnmarshalBinaryLengthPrefixed([]byte(tmTx), &sdkTx) // use tx codec to unmarshal binary length prefix

			txHash := fmt.Sprintf("%X", tmTx.Hash())

			var generalTx types.GeneralTx
			resp, _ := resty.R().Get(ex.cfg.Node.LCDEndpoint + "/txs/" + txHash)
			err := json.Unmarshal(resp.Body(), &generalTx)
			if err != nil {
				fmt.Printf("failed to unmarshal generalTx error: %t\n", err)
			}

			// check log to see if tx is success
			for j, log := range generalTx.Logs {
				if log.Success {
					switch generalTx.Tx.Value.Msg[j].Type {
					case "cosmos-sdk/MsgSend":
						var msgSend types.MsgSend
						err = ex.cdc.UnmarshalJSON(generalTx.Tx.Value.Msg[j].Value, &msgSend)
						if err != nil {
							fmt.Printf("failed to unmarshal msgSend: %t\n", err)
						}

						// switch param in config.yaml
						// this param is to start or stop sending push notifications in case of syncing from the scratch
						if ex.cfg.Alarm.Switch {
							var amount string
							var denom string

							if len(msgSend.Amount) > 0 {
								amount = msgSend.Amount[0].Amount.String()
								denom = msgSend.Amount[0].Denom
							}

							pnp := &types.PushNotificationPayload{
								From:   msgSend.FromAddress,
								To:     msgSend.ToAddress,
								Txid:   txHash,
								Amount: amount,
								Denom:  denom,
							}

							// push notification to both from and to accounts
							nof := notification.New()

							fromAcctStatus := nof.VerifyAccount(msgSend.FromAddress)
							if fromAcctStatus {
								tokens, _ := ex.db.QueryAlarmTokens(msgSend.FromAddress)
								if len(tokens) > 0 {
									nof.PushNotification(pnp, tokens, types.FROM)
								}
							}

							toAcctStatus := nof.VerifyAccount(msgSend.ToAddress)
							if toAcctStatus {
								tokens, _ := ex.db.QueryAlarmTokens(msgSend.ToAddress)
								if len(tokens) > 0 {
									nof.PushNotification(pnp, tokens, types.TO)
								}
							}
						}

					case "cosmos-sdk/MsgMultiSend":
						var multiSendTx types.MsgMultiSend
						err = ex.cdc.UnmarshalJSON(generalTx.Tx.Value.Msg[j].Value, &multiSendTx)
						if err != nil {
							fmt.Printf("failed to unmarshal multiSendTx: %t\n", err)
						}

						nof := notification.New()

						// switch param in config.yaml
						// this param is to start or stop sending push notifications in case of syncing from the scratch
						if ex.cfg.Alarm.Switch {
							for _, input := range multiSendTx.Inputs {
								var amount string
								var denom string

								if len(input.Coins) > 0 {
									amount = input.Coins[0].Amount.String()
									denom = input.Coins[0].Denom
								}

								pnp := &types.PushNotificationPayload{
									From:   input.Address.String(),
									Txid:   txHash,
									Amount: amount,
									Denom:  denom,
								}

								fromAcctStatus := nof.VerifyAccount(input.Address.String())
								if fromAcctStatus {
									tokens, _ := ex.db.QueryAlarmTokens(input.Address.String())
									if len(tokens) > 0 {
										nof.PushNotification(pnp, tokens, types.FROM)
									}
								}
							}

							// push notifications to all outputs
							for _, output := range multiSendTx.Outputs {
								var amount string
								var denom string

								if len(output.Coins) > 0 {
									amount = output.Coins[0].Amount.String()
									denom = output.Coins[0].Denom
								}

								pnp := &types.PushNotificationPayload{
									To:     output.Address.String(),
									Txid:   txHash,
									Amount: amount,
									Denom:  denom,
								}

								toAcctStatus := nof.VerifyAccount(output.Address.String())
								if toAcctStatus {
									tokens, _ := ex.db.QueryAlarmTokens(output.Address.String())
									if len(tokens) > 0 {
										nof.PushNotification(pnp, tokens, types.TO)
									}
								}
							}
						}

					case "cosmos-sdk/MsgCreateValidator":
						var msgCreateValidator types.MsgCreateValidator
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgCreateValidator)

						/*
							[기술적 한계] > 동일한 블록안에 create_validator 메시지가 2개 이상 있을 경우 마지막으로 저장된 id_validator를 가져오면 겹친다.
						*/

						// query the highest height of id_validator
						highestIDValidatorNum, _ := ex.db.QueryHighestValidatorID()

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgCreateValidator.Value.Amount.String(), 64) // parseFloat from sdk.Dec.String()
						newVotingPowerAmount = float64(newVotingPowerAmount) / 1000000

						tempPowerEventHistory := &schema.PowerEventHistory{
							IDValidator:          highestIDValidatorNum + 1,
							Height:               height,
							Proposer:             utils.ConsAddrFromConsPubkey(msgCreateValidator.Pubkey), // new validator's proposer address needs to be converted
							VotingPower:          newVotingPowerAmount,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgCreateValidator.Value.Denom,
							EventType:            types.EventTypeMsgCreateValidator,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						powerEventHistory = append(powerEventHistory, tempPowerEventHistory)

					case "cosmos-sdk/MsgDelegate":
						var msgDelegate types.MsgDelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgDelegate)

						// query validator information fro validator_infos table
						validatorInfo, _ := ex.db.QueryValidator(msgDelegate.ValidatorAddress)

						// query to get id_validator of lastly inserted data
						validatorID, _ := ex.db.QueryValidatorID(validatorInfo.Proposer)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgDelegate.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
						newVotingPowerAmount = newVotingPowerAmount / 1000000

						// current voting power of a validator
						var votingPower float64
						validators, _ := ex.client.Validators(height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorInfo.Proposer {
								votingPower = float64(validator.VotingPower)
							}
						}

						/*
							[기술적 한계] - Certus One 17번째 블록에 두번 - cosmoshub-1
										동일한 블록에서 서로 다른 주소에서 동일한 검증인에게 위임한 트랜잭션이 있을 경우 현재 VotingPower는 같다.
						*/

						tempPowerEventHistory := &schema.PowerEventHistory{
							IDValidator:          validatorID.IDValidator,
							Height:               height,
							Moniker:              validatorInfo.Moniker,
							OperatorAddress:      validatorInfo.OperatorAddress,
							Proposer:             validatorInfo.Proposer,
							VotingPower:          votingPower + newVotingPowerAmount,
							EventType:            types.EventTypeMsgDelegate,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgDelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						powerEventHistory = append(powerEventHistory, tempPowerEventHistory)

					case "cosmos-sdk/MsgUndelegate":
						var msgUndelegate types.MsgUndelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgUndelegate)

						// query validator info
						validatorInfo, _ := ex.db.QueryValidator(msgUndelegate.ValidatorAddress)

						// query to get id_validator of lastly inserted data
						validatorID, _ := ex.db.QueryValidatorID(validatorInfo.Proposer)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgUndelegate.Amount.Amount.String(), 64) // parseFloat from sdk.Dec.String()
						newVotingPowerAmount = -newVotingPowerAmount / 1000000

						// current voting power of a validator
						var votingPower float64
						validators, _ := ex.client.Validators(height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorInfo.Proposer {
								votingPower = float64(validator.VotingPower)
							}
						}

						// substract the undelegated amount from the validator
						tempPowerEventHistory := &schema.PowerEventHistory{
							IDValidator:          validatorID.IDValidator,
							Height:               height,
							Moniker:              validatorInfo.Moniker,
							OperatorAddress:      validatorInfo.OperatorAddress,
							Proposer:             validatorInfo.Proposer,
							VotingPower:          votingPower + newVotingPowerAmount,
							EventType:            types.EventTypeMsgUndelegate,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgUndelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						powerEventHistory = append(powerEventHistory, tempPowerEventHistory)

					case "cosmos-sdk/MsgBeginRedelegate":
						/*
							[Note]
								+ for ValidatorDstAddress
								- for ValidatorSrcAddress
						*/

						var msgBeginRedelegate types.MsgBeginRedelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgBeginRedelegate)

						// query validator_dst_address info
						validatorDstInfo, _ := ex.db.QueryValidator(msgBeginRedelegate.ValidatorDstAddress)
						dstpowerEventHistory, _ := ex.db.QueryValidatorID(validatorDstInfo.Proposer)

						// query validator_src_address info
						validatorSrcInfo, _ := ex.db.QueryValidator(msgBeginRedelegate.ValidatorSrcAddress)
						srcpowerEventHistory, _ := ex.db.QueryValidatorID(validatorSrcInfo.Proposer)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgBeginRedelegate.Amount.Amount.String(), 64)
						newVotingPowerAmount = newVotingPowerAmount / 1000000

						// current destination validator's voting power
						var dstValidatorVotingPower float64
						validators, _ := ex.client.Validators(height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorDstInfo.Proposer {
								dstValidatorVotingPower = float64(validator.VotingPower)
							}
						}

						// current source validator's voting power
						var srcValidatorVotingPower float64
						validators, _ = ex.client.Validators(height)
						for _, validator := range validators.Validators {
							if validator.Address.String() == validatorSrcInfo.Proposer {
								srcValidatorVotingPower = float64(validator.VotingPower)
							}
						}

						// add the redelegated amount to validator_dst_address
						tempDstpowerEventHistory := &schema.PowerEventHistory{
							IDValidator:          dstpowerEventHistory.IDValidator,
							Height:               height,
							Moniker:              validatorDstInfo.Moniker,
							OperatorAddress:      validatorDstInfo.OperatorAddress,
							Proposer:             validatorDstInfo.Proposer,
							VotingPower:          dstValidatorVotingPower + newVotingPowerAmount,
							EventType:            types.EventTypeMsgBeginRedelegate,
							NewVotingPowerAmount: newVotingPowerAmount,
							NewVotingPowerDenom:  msgBeginRedelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						powerEventHistory = append(powerEventHistory, tempDstpowerEventHistory)

						// substract the redelegated amount from validator_src_address
						tempSrcpowerEventHistory := &schema.PowerEventHistory{
							IDValidator:          srcpowerEventHistory.IDValidator,
							Height:               height,
							Moniker:              validatorSrcInfo.Moniker,
							OperatorAddress:      validatorSrcInfo.OperatorAddress,
							Proposer:             validatorSrcInfo.Proposer,
							VotingPower:          srcValidatorVotingPower - newVotingPowerAmount,
							EventType:            types.EventTypeMsgBeginRedelegate,
							NewVotingPowerAmount: -newVotingPowerAmount,
							NewVotingPowerDenom:  msgBeginRedelegate.Amount.Denom,
							TxHash:               generalTx.TxHash,
							Time:                 block.BlockMeta.Header.Time,
						}
						powerEventHistory = append(powerEventHistory, tempSrcpowerEventHistory)

					case "cosmos-sdk/MsgSubmitProposal":
						var msgSubmitProposal types.MsgSubmitProposal
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

						tempProposal := &schema.Proposal{
							ID:                   proposalID,
							TxHash:               generalTx.TxHash,
							Proposer:             msgSubmitProposal.Proposer,
							InitialDepositAmount: initialDepositAmount,
							InitialDepositDenom:  initialDepositDenom,
						}
						proposal = append(proposal, tempProposal)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempDeposit := &schema.Deposit{
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
						deposit = append(deposit, tempDeposit)

					case "cosmos-sdk/MsgVote":
						var msgVote types.MsgVote
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgVote)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						proposalID, _ := strconv.ParseInt(msgVote.ProposalID, 10, 64)
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempVote := &schema.Vote{
							Height:     height,
							ProposalID: proposalID,
							Voter:      msgVote.Voter,
							Option:     msgVote.Option,
							TxHash:     generalTx.TxHash,
							GasWanted:  gasWanted,
							GasUsed:    gasUsed,
							Time:       block.BlockMeta.Header.Time,
						}
						vote = append(vote, tempVote)

					case "cosmos-sdk/MsgDeposit":
						var msgDeposit types.MsgDeposit
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgDeposit)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						proposalID, _ := strconv.ParseInt(msgDeposit.ProposalID, 10, 64)
						amount := msgDeposit.Amount[0].Amount
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempDeposit := &schema.Deposit{
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
						deposit = append(deposit, tempDeposit)

					default:
						continue
					}
				}
			}
		}
	}

	return vote, deposit, proposal, powerEventHistory, nil
}

// getTxs returns transactions information in a block
func (ex *Exporter) getTxs(txResp []sdk.TxResponse) ([]*schema.TxCosmoshub3, error) {
	txs := make([]*schema.TxCosmoshub3, 0)

	if len(txResp) > 0 {
		for _, tx := range txResp {
			stdTx, ok := tx.Tx.(auth.StdTx)
			if !ok {
				return txs, fmt.Errorf("unsupported tx type: %T", tx.Tx)
			}

			msgsBz, err := ceCodec.Codec.MarshalJSON(stdTx.GetMsgs())
			if err != nil {
				return txs, fmt.Errorf("failed to unmarshal tx messages: %t", err)
			}

			feeBz, err := ceCodec.Codec.MarshalJSON(stdTx.Fee)
			if err != nil {
				return txs, fmt.Errorf("failed to unmarshal tx fee: %t", err)
			}

			// convert Tendermint signatures into a more human-readable format
			sigs := make([]types.Signature, len(stdTx.GetSignatures()), len(stdTx.GetSignatures()))
			for i, sig := range stdTx.GetSignatures() {
				consPubKey, err := sdk.Bech32ifyConsPub(sig.PubKey) // nolint: typecheck
				if err != nil {
					return txs, fmt.Errorf("failed to convert validator public key %t\n: %t", sig.PubKey, err)
				}

				sigs[i] = types.Signature{
					Address:   sig.Address().String(),
					Signature: base64.StdEncoding.EncodeToString(sig.Signature),
					Pubkey:    consPubKey,
				}
			}

			sigsBz, err := ceCodec.Codec.MarshalJSON(sigs)
			if err != nil {
				return txs, fmt.Errorf("failed to unmarshal tx signatures: %t", err)
			}

			eventsBz, err := ceCodec.Codec.MarshalJSON(tx.Events)
			if err != nil {
				return txs, fmt.Errorf("failed to unmarshal tx events: %t", err)
			}

			logsBz, err := ceCodec.Codec.MarshalJSON(tx.Logs)
			if err != nil {
				return txs, fmt.Errorf("failed to unmarshal tx logs: %t", err)
			}

			tempTx := &schema.TxCosmoshub3{
				Height:     tx.Height,
				TxHash:     tx.TxHash,
				GasWanted:  tx.GasWanted,
				GasUsed:    tx.GasUsed,
				Messages:   string(msgsBz),
				Fee:        string(feeBz),
				Signatures: string(sigsBz),
				Logs:       string(logsBz),
				Events:     string(eventsBz),
				Memo:       stdTx.GetMemo(),
				Time:       tx.Timestamp,
			}

			txs = append(txs, tempTx)
		}
	}

	return txs, nil
}
