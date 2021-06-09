package exporter

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmostation/cosmostation-cosmos/notification"
	"github.com/cosmostation/cosmostation-cosmos/types"

	"go.uber.org/zap"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// handlePushNotification handles our mobile wallet applications' push notification.
func (ex *Exporter) handlePushNotification(block *tmctypes.ResultBlock, txResp []*sdktypes.TxResponse) error {
	if len(txResp) <= 0 {
		return nil
	}

	for _, tx := range txResp {
		// Other than code equals to 0, it is failed transaction.
		if tx.Code != 0 {
			continue
		}

		msgs := tx.GetTx().GetMsgs()

		for _, msg := range msgs {

			// stdTx, ok := tx.Tx.(auth.StdTx)
			// if !ok {
			// 	return fmt.Errorf("unsupported tx type: %s", tx.Tx)
			// }

			switch m := msg.(type) {
			// case bank.MsgSend:
			case *banktypes.MsgSend:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgSend := m.(banktypestypes.MsgSend)

				var amount string
				var denom string

				// TODO: need to test for multiple coins in one message.
				if len(m.Amount) > 0 {
					amount = m.Amount[0].Amount.String()
					denom = m.Amount[0].Denom
				}

				payload := types.NewNotificationPayload(types.NotificationPayload{
					From:   m.FromAddress,
					To:     m.ToAddress,
					Txid:   tx.TxHash,
					Amount: amount,
					Denom:  denom,
				})

				// Push notification to both sender and recipient.
				notification := notification.NewNotification()

				fromAccountStatus := notification.VerifyAccountStatus(m.FromAddress)
				if fromAccountStatus {
					tokens, _ := ex.db.QueryAlarmTokens(m.FromAddress)
					if len(tokens) > 0 {
						notification.Push(*payload, tokens, types.From)
					}
				}

				toAccountStatus := notification.VerifyAccountStatus(m.ToAddress)
				if toAccountStatus {
					tokens, _ := ex.db.QueryAlarmTokens(m.ToAddress)
					if len(tokens) > 0 {
						notification.Push(*payload, tokens, types.To)
					}
				}

			case *banktypes.MsgMultiSend:
				zap.S().Infof("MsgType: %s | Hash: %s", m.Type(), tx.TxHash)

				// msgMultiSend := m.(banktypes.MsgMultiSend)

				notification := notification.NewNotification()

				// Push notifications to all accounts in inputs
				for _, input := range m.Inputs {
					var amount string
					var denom string

					if len(input.Coins) > 0 {
						amount = input.Coins[0].Amount.String()
						denom = input.Coins[0].Denom
					}

					payload := &types.NotificationPayload{
						From:   input.Address,
						Txid:   tx.TxHash,
						Amount: amount,
						Denom:  denom,
					}

					fromAccountStatus := notification.VerifyAccountStatus(input.Address)
					if fromAccountStatus {
						tokens, _ := ex.db.QueryAlarmTokens(input.Address)
						if len(tokens) > 0 {
							notification.Push(*payload, tokens, types.From)
						}
					}
				}

				// Push notifications to all accounts in outputs
				for _, output := range m.Outputs {
					var amount string
					var denom string

					if len(output.Coins) > 0 {
						amount = output.Coins[0].Amount.String()
						denom = output.Coins[0].Denom
					}

					payload := &types.NotificationPayload{
						To:     output.Address,
						Txid:   tx.TxHash,
						Amount: amount,
						Denom:  denom,
					}

					toAcctStatus := notification.VerifyAccountStatus(output.Address)
					if toAcctStatus {
						tokens, _ := ex.db.QueryAlarmTokens(output.Address)
						if len(tokens) > 0 {
							notification.Push(*payload, tokens, types.To)
						}
					}
				}

			default:
				continue
			}
		}
	}

	return nil
}
