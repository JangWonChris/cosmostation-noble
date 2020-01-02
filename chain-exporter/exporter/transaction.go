package exporter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/chain-exporter/codec"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/notification"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/schema"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/utils"

	resty "gopkg.in/resty.v1"
)

// getTransactionInfo provides information about each transaction in every block
func (ces ChainExporterService) getTransactionInfo(height int64) ([]*schema.TransactionInfo, []*schema.VoteInfo,
	[]*schema.DepositInfo, []*schema.ProposalInfo, []*schema.ValidatorSetInfo, error) {

	transactionInfo := make([]*schema.TransactionInfo, 0)
	voteInfo := make([]*schema.VoteInfo, 0)
	depositInfo := make([]*schema.DepositInfo, 0)
	proposalInfo := make([]*schema.ProposalInfo, 0)
	validatorSetInfo := make([]*schema.ValidatorSetInfo, 0)

	// query current block
	block, err := ces.rpcClient.Block(&height)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if len(block.Block.Data.Txs) > 0 {
		for _, tmTx := range block.Block.Data.Txs {
			// use tx codec to unmarshal binary length prefix
			var sdkTx sdk.Tx
			_ = ces.codec.UnmarshalBinaryLengthPrefixed([]byte(tmTx), &sdkTx)

			txHash := fmt.Sprintf("%X", tmTx.Hash())

			tx, err := ces.Tx(txHash)
			if err != nil {
				fmt.Printf("failed to get tx %s: %s", txHash, err)
				continue
			}

			// Save txInfo in PostgreSQL database
			var tempTxInfo schema.TransactionInfo
			tempTxInfo, err = ces.SetTx(tx, txHash)
			if err != nil {
				fmt.Printf("failed to persist transaction %s: %s", txHash, err)
			}
			transactionInfo = append(transactionInfo, &tempTxInfo)

			var generalTx types.GeneralTx
			resp, _ := resty.R().Get(ces.config.Node.LCDURL + "/txs/" + txHash)
			err = json.Unmarshal(resp.Body(), &generalTx)
			if err != nil {
				fmt.Printf("unmarshal generalTx error - %v\n", err)
			}

			// check log to see if tx is success
			for j, log := range generalTx.Logs {
				if log.Success {
					switch generalTx.Tx.Value.Msg[j].Type {
					case "cosmos-sdk/MsgSend":
						var msgSend types.MsgSend
						err = ces.codec.UnmarshalJSON(generalTx.Tx.Value.Msg[j].Value, &msgSend)
						if err != nil {
							fmt.Printf("failed to JSON encode msgSend: %s", err)
						}

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

						fromAccount := nof.VerifyAccount(msgSend.FromAddress)
						if fromAccount != nil {
							nof.PushNotification(pnp, fromAccount.AlarmToken, types.FROM)
						}

						toAccount := nof.VerifyAccount(msgSend.ToAddress)
						if toAccount != nil {
							nof.PushNotification(pnp, toAccount.AlarmToken, types.TO)
						}

					case "cosmos-sdk/MultiSend":
						var multiSendTx types.MsgMultiSend
						err = ces.codec.UnmarshalJSON(generalTx.Tx.Value.Msg[j].Value, &multiSendTx)
						if err != nil {
							fmt.Println("Unmarshal MsgMultiSend JSON Error: ", err)
						}

						fmt.Println("=======================================[multisend]")
						fmt.Println(multiSendTx.Inputs)
						fmt.Println(multiSendTx.Outputs)
						fmt.Println("=======================================")

					case "cosmos-sdk/MsgCreateValidator":
						var msgCreateValidator types.MsgCreateValidator
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgCreateValidator)

						/*
							[기술적 한계] > 동일한 블록안에 create_validator 메시지가 2개 이상 있을 경우 마지막으로 저장된 id_validator를 가져오면 겹친다.
						*/

						// query the highest height of id_validator
						highestIDValidatorNum, _ := ces.db.QueryHighestIDValidatorNum()

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						newVotingPowerAmount, _ := strconv.ParseFloat(msgCreateValidator.Value.Amount.String(), 64) // parseFloat from sdk.Dec.String()
						newVotingPowerAmount = float64(newVotingPowerAmount) / 1000000

						tempValidatorSetInfo := &schema.ValidatorSetInfo{
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
						validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

					case "cosmos-sdk/MsgDelegate":
						var msgDelegate types.MsgDelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgDelegate)

						// query validator information fro validator_infos table
						validatorInfo, _ := ces.db.QueryValidatorInfo(msgDelegate.ValidatorAddress)

						// query to get id_validator of lastly inserted data
						idValidatorSetInfo, _ := ces.db.QueryIDValidatorSetInfo(validatorInfo.Proposer)

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

						tempValidatorSetInfo := &schema.ValidatorSetInfo{
							IDValidator:          idValidatorSetInfo.IDValidator,
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
						validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

					case "cosmos-sdk/MsgUndelegate":
						var msgUndelegate types.MsgUndelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgUndelegate)

						// query validator info
						validatorInfo, _ := ces.db.QueryValidatorInfo(msgUndelegate.ValidatorAddress)

						// query to get id_validator of lastly inserted data
						idValidatorSetInfo, _ := ces.db.QueryIDValidatorSetInfo(validatorInfo.Proposer)

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
						tempValidatorSetInfo := &schema.ValidatorSetInfo{
							IDValidator:          idValidatorSetInfo.IDValidator,
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
						validatorSetInfo = append(validatorSetInfo, tempValidatorSetInfo)

					case "cosmos-sdk/MsgBeginRedelegate":
						/*
							[Note]
								+ for ValidatorDstAddress
								- for ValidatorSrcAddress
						*/

						var msgBeginRedelegate types.MsgBeginRedelegate
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgBeginRedelegate)

						// query validator_dst_address info
						validatorDstInfo, _ := ces.db.QueryValidatorInfo(msgBeginRedelegate.ValidatorDstAddress)
						dstValidatorSetInfo, _ := ces.db.QueryIDValidatorSetInfo(validatorDstInfo.Proposer)

						// query validator_src_address info
						validatorSrcInfo, _ := ces.db.QueryValidatorInfo(msgBeginRedelegate.ValidatorSrcAddress)
						srcValidatorSetInfo, _ := ces.db.QueryIDValidatorSetInfo(validatorSrcInfo.Proposer)

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
						tempDstValidatorSetInfo := &schema.ValidatorSetInfo{
							IDValidator:          dstValidatorSetInfo.IDValidator,
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
						validatorSetInfo = append(validatorSetInfo, tempDstValidatorSetInfo)

						// substract the redelegated amount from validator_src_address
						tempSrcValidatorSetInfo := &schema.ValidatorSetInfo{
							IDValidator:          srcValidatorSetInfo.IDValidator,
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
						validatorSetInfo = append(validatorSetInfo, tempSrcValidatorSetInfo)

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

						tempProposalInfo := &schema.ProposalInfo{
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

						tempDepositInfo := &schema.DepositInfo{
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
						var msgVote types.MsgVote
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgVote)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						proposalID, _ := strconv.ParseInt(msgVote.ProposalID, 10, 64)
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempVoteInfo := &schema.VoteInfo{
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
						var msgDeposit types.MsgDeposit
						_ = json.Unmarshal(generalTx.Tx.Value.Msg[j].Value, &msgDeposit)

						height, _ := strconv.ParseInt(generalTx.Height, 10, 64)
						proposalID, _ := strconv.ParseInt(msgDeposit.ProposalID, 10, 64)
						amount := msgDeposit.Amount[0].Amount
						gasWanted, _ := strconv.ParseInt(generalTx.GasWanted, 10, 64)
						gasUsed, _ := strconv.ParseInt(generalTx.GasUsed, 10, 64)

						tempDepositInfo := &schema.DepositInfo{
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

// Tx queries for a transaction from the REST client and decodes it into a sdk.Tx
// if the transaction exists. An error is returned if the tx doesn't exist or
// decoding fails.
func (ces ChainExporterService) Tx(hash string) (sdk.TxResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/txs/%s", ces.config.Node.LCDURL, hash))
	if err != nil {
		return sdk.TxResponse{}, err
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var tx sdk.TxResponse

	if err := ces.codec.UnmarshalJSON(bz, &tx); err != nil {
		return sdk.TxResponse{}, err
	}

	return tx, nil
}

// SetTx stores a transaction and returns the resulting record ID. An error is
// returned if the operation fails.
func (ces ChainExporterService) SetTx(tx sdk.TxResponse, txHash string) (schema.TransactionInfo, error) {
	stdTx, ok := tx.Tx.(auth.StdTx)
	if !ok {
		return schema.TransactionInfo{}, fmt.Errorf("unsupported tx type: %T", tx.Tx)
	}

	msgsBz, err := ceCodec.Codec.MarshalJSON(stdTx.GetMsgs())
	if err != nil {
		return schema.TransactionInfo{}, fmt.Errorf("failed to JSON encode tx messages: %s", err)
	}

	feeBz, err := ceCodec.Codec.MarshalJSON(stdTx.Fee)
	if err != nil {
		return schema.TransactionInfo{}, fmt.Errorf("failed to JSON encode tx fee: %s", err)
	}

	// convert Tendermint signatures into a more human-readable format
	sigs := make([]types.Signature, len(stdTx.GetSignatures()), len(stdTx.GetSignatures()))
	for i, sig := range stdTx.GetSignatures() {
		consPubKey, err := sdk.Bech32ifyConsPub(sig.PubKey) // nolint: typecheck
		if err != nil {
			return schema.TransactionInfo{}, fmt.Errorf("failed to convert validator public key %s: %s\n", sig.PubKey, err)
		}

		sigs[i] = types.Signature{
			Address:   sig.Address().String(),
			Signature: base64.StdEncoding.EncodeToString(sig.Signature),
			Pubkey:    consPubKey,
		}
	}

	sigsBz, err := ceCodec.Codec.MarshalJSON(sigs)
	if err != nil {
		return schema.TransactionInfo{}, fmt.Errorf("failed to JSON encode tx signatures: %s", err)
	}

	eventsBz, err := ceCodec.Codec.MarshalJSON(tx.Events)
	if err != nil {
		return schema.TransactionInfo{}, fmt.Errorf("failed to JSON encode tx events: %s", err)
	}

	logsBz, err := ceCodec.Codec.MarshalJSON(tx.Logs)
	if err != nil {
		return schema.TransactionInfo{}, fmt.Errorf("failed to JSON encode tx logs: %s", err)
	}

	return schema.TransactionInfo{
		Height:     tx.Height,
		TxHash:     txHash,
		GasWanted:  tx.GasWanted,
		GasUsed:    tx.GasUsed,
		Messages:   string(msgsBz),
		Fee:        string(feeBz),
		Signatures: string(sigsBz),
		Logs:       string(logsBz),
		Events:     string(eventsBz),
		Memo:       stdTx.GetMemo(),
		Time:       tx.Timestamp,
	}, err
}
