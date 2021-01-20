package exporter

import (
	"encoding/json"
	"fmt"
	"log"

	// internal
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/custom"

	// core
	"github.com/cosmostation/mintscan-backend-library/db/schema"
	"github.com/cosmostation/mintscan-backend-library/types"

	// sdk
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktypestx "github.com/cosmos/cosmos-sdk/types/tx"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getTxs decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getTxs(block *tmctypes.ResultBlock, txResps []*sdktypes.TxResponse) ([]schema.Transaction, error) {
	txs := make([]schema.Transaction, 0)

	if len(txResps) <= 0 {
		return txs, nil
	}

	for _, txResp := range txResps {
		txI := txResp.GetTx()
		tx, ok := txI.(*sdktypestx.Tx)
		if !ok {
			return txs, fmt.Errorf("unsupported type")
		}

		msgs := tx.GetBody().GetMessages()
		jsonRaws := make([]json.RawMessage, len(msgs), len(msgs))
		var err error
		for i, msg := range msgs {
			jsonRaws[i], err = custom.AppCodec.MarshalJSON(msg)
			if err != nil {
				return txs, fmt.Errorf("failed to marshal message of transaction : %s", err)
			}
		}
		msgsBz, err := json.Marshal(jsonRaws)
		if err != nil {
			return txs, fmt.Errorf("failed to marshal set of transactions : %s", err)
		}

		feeBz, err := custom.AppCodec.MarshalJSON(tx.GetAuthInfo().GetFee())
		if err != nil {
			return txs, fmt.Errorf("failed to marshal tx fee: %s", err)
		}

		type SIG struct {
			Signatures []byte
		}

		sigs := make([]SIG, len(tx.GetSignatures()), len(tx.GetSignatures()))
		for i, s := range tx.GetSignatures() {
			sigs[i].Signatures = s
		}

		sigsBz, err := json.Marshal(sigs)
		if err != nil {
			return txs, fmt.Errorf("failed to marshal tx signatures: %s", err)
		}

		logsBz, err := json.Marshal(txResp.Logs.String())
		if err != nil {
			return txs, fmt.Errorf("failed to marshal tx logs: %s", err)
		}

		t := &schema.Transaction{
			ChainID:    block.Block.ChainID,
			Height:     txResp.Height,
			Code:       txResp.Code,
			TxHash:     txResp.TxHash,
			GasWanted:  txResp.GasWanted,
			GasUsed:    txResp.GasUsed,
			Messages:   string(msgsBz),
			Fee:        string(feeBz),
			Signatures: string(sigsBz),
			Logs:       string(logsBz),
			RawLog:     txResp.RawLog,
			Memo:       tx.GetBody().Memo,
			Timestamp:  txResp.Timestamp,
		}

		txs = append(txs, *t)
	}

	return txs, nil
}

// getTxsChunk decodes transactions in a block and return a format of database transaction.
func (ex *Exporter) getTxsJSONChunk(block *tmctypes.ResultBlock, txResps []*sdktypes.TxResponse) ([]schema.RawTransaction, error) {
	txChunk := make([]schema.RawTransaction, len(txResps), len(txResps))
	if len(txResps) <= 0 {
		return txChunk, nil
	}

	for i, txResp := range txResps {
		chunk, err := custom.AppCodec.MarshalJSON(txResp)
		if err != nil {
			log.Println(err)
			return txChunk, fmt.Errorf("failed to marshal tx : %s", err)
		}
		txChunk[i].ChainID = block.Block.ChainID
		txChunk[i].Height = txResp.Height
		txChunk[i].TxHash = txResp.TxHash
		txChunk[i].Chunk = chunk
		// show result
		// fmt.Println(jsonString[i])
	}

	return txChunk, nil
}

func (ex *Exporter) transactionAccount(chainID string, txResps []*sdktypes.TxResponse) (tms []schema.TransactionAccount, err error) {
	if len(txResps) <= 0 {
		return nil, nil
	}

	for _, txResp := range txResps {
		msgs := txResp.GetTx().GetMsgs()

		txHash := txResp.TxHash
		height := txResp.Height

		if len(msgs) > 0 {
			//모든 msgs는 동일한 서명인을 갖기 때문에, index = 0 일때만 추가
			sub := getAccountSlice(chainID, "signer", height, txHash, getSignerAddress(msgs[0].GetSigners())...)
			tms = append(tms, sub...)
		}
		for _, msg := range msgs {

			msgType, accounts := types.AccountExporterFromCosmosTxMsg(&msg)

			// 코스모스 기본 메세지 타입이 아니면, msgType이 비어있다.
			if msgType == "" {
				msgType, accounts = custom.AccountExporterFromCustomTxMsg(&msg, txResp.TxHash)
				if msgType == "" {
					fmt.Printf("Undefined msg Type : %T(hash = %s)\n", msg, txResp.TxHash)
				}
			}
			if len(accounts) > 0 {
				sub := getAccountSlice(chainID, msgType, height, txHash, accounts...)
				// Tx 내에서 발생한 msg별로 어카운트를 수집했으므로, 슬라이스에 저장한다.
				tms = append(tms, sub...)
			}
		} // end msgs for loop

	} // 모든 tx 완료
	return tms, nil
}

func getAccountSlice(chainID, msgType string, height int64, txHash string, accounts ...string) (tms []schema.TransactionAccount) {
	tms = make([]schema.TransactionAccount, len(accounts), len(accounts))
	for i, acc := range accounts {
		tms[i].ChainID = chainID
		tms[i].MsgType = msgType
		tms[i].Height = height
		tms[i].TxHash = txHash
		tms[i].AccountAddress = acc
	}
	return tms
}

func getSignerAddress(accAddrs []sdktypes.AccAddress) (address []string) {
	for _, addr := range accAddrs {
		address = append(address, addr.String())
	}

	return address
}

// func (ex *Exporter) transactionAccount(txResps []*sdktypes.TxResponse) (tms []schema.TransactionAccount, err error) {
// 	// 별도의 스키마 필요
// 	// id, account, hash, timestamp(불필요, 조인하면 되기 때문)
// 	// id, tx_hash(id는 불가능, db에 저장할 때 알 수 없음)
// 	if len(txResps) <= 0 {
// 		return nil, nil
// 	}

// 	for _, txResp := range txResps {
// 		msgs := txResp.GetTx().GetMsgs()

// 		txHash := txResp.TxHash
// 		accounts := make([]string, 0)

// 		for _, msg := range msgs {
// 			switch msg := msg.(type) {
// 			case *authvestingtypes.MsgCreateVestingAccount:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.FromAddress, msg.ToAddress)

// 				//bank
// 			case *banktypes.MsgSend:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.FromAddress, msg.ToAddress)
// 			case *banktypes.MsgMultiSend:
// 				// 추가 필요
// 				fmt.Printf("Type : %T\n", msg)
// 				fmt.Println(msg)

// 				//crisis
// 			case *crisistypes.MsgVerifyInvariant:
// 				fmt.Printf("Type : %T\n", msg)
// 				fmt.Println(msg.Sender)
// 				accounts = append(accounts, msg.Sender)

// 				//distribution
// 			case *distributiontypes.MsgSetWithdrawAddress:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.DelegatorAddress)
// 			case *distributiontypes.MsgWithdrawDelegatorReward:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.DelegatorAddress, msg.ValidatorAddress)
// 			case *distributiontypes.MsgWithdrawValidatorCommission:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.ValidatorAddress)
// 			case *distributiontypes.MsgFundCommunityPool:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.Depositor)

// 				//evidence
// 			case *evidencetypes.MsgSubmitEvidence:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.Submitter)

// 				//gov
// 			case *govtypes.MsgSubmitProposal:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.Proposer)
// 			case *govtypes.MsgVote:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.Voter)
// 			case *govtypes.MsgDeposit:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.Depositor)

// 				//slashing
// 			case *slashingtypes.MsgUnjail:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.ValidatorAddr)

// 				//staking
// 			case *stakingtypes.MsgCreateValidator:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.DelegatorAddress, msg.ValidatorAddress)
// 			case *stakingtypes.MsgEditValidator:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.ValidatorAddress)
// 			case *stakingtypes.MsgDelegate:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.DelegatorAddress, msg.ValidatorAddress)
// 			case *stakingtypes.MsgBeginRedelegate:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.ValidatorDstAddress)
// 			case *stakingtypes.MsgUndelegate:
// 				fmt.Printf("Type : %T\n", msg)
// 				accounts = append(accounts, msg.DelegatorAddress, msg.ValidatorAddress)

// 				//ibc transaction 추가/보완 필요
// 			case *transfertypes.MsgTransfer:
// 				fmt.Printf("Type : %T\n", msg)
// 				log.Println(msg.Sender)
// 				log.Println(msg.Receiver) //체인 밖 주소이므로 필요 없을 것 같다.
// 				accounts = append(accounts, msg.Sender)

// 			//client, connection, channel은 별도로 얘기가 필요하다.

// 			default:
// 				fmt.Printf("Undefined Type : %T\n", msg)
// 			}

// 		} // end msgs for loop

// 		// 하나의 Tx에서 발생한 모든 어카운트를 수집했으므로, 슬라이스로 저장한다.
// 		sub, err := getAccountSlice(txHash, accounts...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tms = append(tms, sub...)

// 	} // 모든 tx 완료
// 	return tms, nil
// }

// func getAccountSlice(txHash string, accounts ...string) (tms []schema.TransactionAccount, err error) {
// 	tms = make([]schema.TransactionAccount, len(accounts), len(accounts))
// 	for i, acc := range accounts {
// 		tms[i].TxHash = txHash
// 		tms[i].AccountAddress = acc
// 	}
// 	return tms, nil
// }
