package exporter

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/notification"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"go.uber.org/zap"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// handlePushNotification handles our mobile wallet applications' push notification.
func (ex *Exporter) handlePushNotification(block *tmctypes.ResultBlock, txs []*sdk.TxResponse) error {
	if len(txs) <= 0 {
		return nil
	}

	for _, tx := range txs {
		// Other than code equals to 0, it is failed transaction.
		if tx.Code != 0 {
			return nil
		}

		stdTx, ok := tx.Tx.(auth.StdTx)
		if !ok {
			return fmt.Errorf("unsupported tx type: %s", tx.Tx)
		}

		switch stdTx.Msgs[0].(type) {
		case bank.MsgSend:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgSend := stdTx.Msgs[0].(bank.MsgSend)

			var amount string
			var denom string

			// TODO: need to test for multiple coins in one message.
			if len(msgSend.Amount) > 0 {
				amount = msgSend.Amount[0].Amount.String()
				denom = msgSend.Amount[0].Denom
			}

			payload := types.NewNotificationPayload(types.NotificationPayload{
				From:   msgSend.FromAddress.String(),
				To:     msgSend.ToAddress.String(),
				Txid:   tx.TxHash,
				Amount: amount,
				Denom:  denom,
			})

			// Push notification to both sender and recipient.
			notification := notification.NewNotification()

			fromAccountStatus := notification.VerifyAccountStatus(msgSend.FromAddress.String())
			if fromAccountStatus {
				tokens, _ := ex.db.QueryAlarmTokens(msgSend.FromAddress.String())
				if len(tokens) > 0 {
					notification.Push(*payload, tokens, types.From)
				}
			}

			toAccountStatus := notification.VerifyAccountStatus(msgSend.ToAddress.String())
			if toAccountStatus {
				tokens, _ := ex.db.QueryAlarmTokens(msgSend.ToAddress.String())
				if len(tokens) > 0 {
					notification.Push(*payload, tokens, types.To)
				}
			}

		case bank.MsgMultiSend:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgMultiSend := stdTx.Msgs[0].(bank.MsgMultiSend)

			notification := notification.NewNotification()

			// Push notifications to all accounts in inputs
			for _, input := range msgMultiSend.Inputs {
				var amount string
				var denom string

				if len(input.Coins) > 0 {
					amount = input.Coins[0].Amount.String()
					denom = input.Coins[0].Denom
				}

				payload := &types.NotificationPayload{
					From:   input.Address.String(),
					Txid:   tx.TxHash,
					Amount: amount,
					Denom:  denom,
				}

				fromAccountStatus := notification.VerifyAccountStatus(input.Address.String())
				if fromAccountStatus {
					tokens, _ := ex.db.QueryAlarmTokens(input.Address.String())
					if len(tokens) > 0 {
						notification.Push(*payload, tokens, types.From)
					}
				}
			}

			// Push notifications to all accounts in outputs
			for _, output := range msgMultiSend.Outputs {
				var amount string
				var denom string

				if len(output.Coins) > 0 {
					amount = output.Coins[0].Amount.String()
					denom = output.Coins[0].Denom
				}

				payload := &types.NotificationPayload{
					To:     output.Address.String(),
					Txid:   tx.TxHash,
					Amount: amount,
					Denom:  denom,
				}

				toAcctStatus := notification.VerifyAccountStatus(output.Address.String())
				if toAcctStatus {
					tokens, _ := ex.db.QueryAlarmTokens(output.Address.String())
					if len(tokens) > 0 {
						notification.Push(*payload, tokens, types.To)
					}
				}
			}

		default:
			continue
		}
	}

	return nil
}
