package exporter

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"go.uber.org/zap"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// handlePushNotification handles our mobile wallet applications' push notification.
func (ex *Exporter) handlePushNotification(block *tmctypes.ResultBlock, txs []*sdkTypes.TxResponse) error {
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

			// Push notification to both sender and recipient
			fromAccount, err := ex.db.QueryAppAccount(msgSend.FromAddress.String())
			if err != nil {
				return fmt.Errorf("unexpected database error: %s", err)
			}

			for _, acct := range fromAccount {
				// Send push notification when alarm token is not empty and status is true.
				if acct.AlarmToken != "" || acct.AlarmStatus {
					ex.notiClient.Push(*payload, acct.AlarmToken, types.From)
				}
			}

			toAccount, err := ex.db.QueryAppAccount(msgSend.ToAddress.String())
			if err != nil {
				return fmt.Errorf("unexpected database error: %s", err)
			}

			for _, acct := range toAccount {
				// Send push notification when alarm token is not empty and status is true.
				if acct.AlarmToken != "" || acct.AlarmStatus {
					ex.notiClient.Push(*payload, acct.AlarmToken, types.To)
				}
			}

		case bank.MsgMultiSend:
			zap.S().Infof("MsgType: %s | Hash: %s", stdTx.Msgs[0].Type(), tx.TxHash)

			msgMultiSend := stdTx.Msgs[0].(bank.MsgMultiSend)

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

				// Handle from address
				inputAccounts, err := ex.db.QueryAppAccount(input.Address.String())
				if err != nil {
					return fmt.Errorf("unexpected database error: %s", err)
				}

				for _, acct := range inputAccounts {
					// Send push notification when alarm token is not empty and status is true.
					if acct.AlarmToken != "" || acct.AlarmStatus {
						ex.notiClient.Push(*payload, acct.AlarmToken, types.From)
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

				// Handle to address
				outputAccounts, err := ex.db.QueryAppAccount(output.Address.String())
				if err != nil {
					return fmt.Errorf("unexpected database error: %s", err)
				}

				for _, acct := range outputAccounts {
					// Send push notification when alarm token is not empty and status is true.
					if acct.AlarmToken != "" || acct.AlarmStatus {
						ex.notiClient.Push(*payload, acct.AlarmToken, types.To)
					}
				}
			}

		default:
			continue
		}
	}

	return nil
}
