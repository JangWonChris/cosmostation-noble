package notification

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/db"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	resty "gopkg.in/resty.v1"
)

// Notification implemnts a wrapper around configuration for this project
type Notification struct {
	cfg *config.Config
	db  *db.Database
}

func New() Notification {
	cfg := config.ParseConfig()

	return Notification{
		cfg: cfg,
		db:  db.Connect(&cfg.DB),
	}
}

// PushNotification sends push notification to its respective device
func (nof *Notification) PushNotification(pnp *types.PushNotificationPayload, tokens []string, target string) {
	var pns []types.PushNotifications

	switch target {
	case "from":
		tempNotification := types.PushNotifications{
			Tokens:   tokens,
			Platform: 2,
			Title:    types.PushNotificationSentTitle + pnp.Amount + pnp.Denom,
			Message:  types.PushNotificationSentMessage + pnp.Amount + pnp.Denom,
			Data: types.PushNotificationData{
				NotifyTo: pnp.From,
				Txid:     pnp.Txid,
				Type:     types.SENT,
			},
		}
		pns = append(pns, tempNotification)
		fmt.Printf("sent push notification - Hash: %s, From: %s \n", pnp.Txid, pnp.From)
	case "to":
		tempNotification := types.PushNotifications{
			Tokens:   tokens,
			Platform: 2,
			Title:    types.PushNotificationReceivedTitle + pnp.Amount + pnp.Denom,
			Message:  types.PushNotificationReceivedMessage + pnp.Amount + pnp.Denom,
			Data: types.PushNotificationData{
				NotifyTo: pnp.To,
				Txid:     pnp.Txid,
				Type:     types.RECEIVED,
			},
		}
		pns = append(pns, tempNotification)
		fmt.Printf("sent push notification - Hash: %s, To: %s \n", pnp.Txid, pnp.To)
	default:
		fmt.Printf("invalid target: %s ", target)
	}

	pnsp := types.PushNotificationServerPayload{
		Notifications: pns,
	}

	// send push notification
	_, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(pnsp).
		Post(nof.cfg.Alarm.PushServerEndpoint)
	if err != nil {
		fmt.Printf("failed to push notification %s: ", err)
	}
}

// VerifyAccount verifes account before sending push notification
func (nof *Notification) VerifyAccount(address string) bool {
	var account types.Account
	account, _ = nof.db.QueryAccount(address)

	// return when data is empty
	if account.AlarmToken == "" {
		return false
	}

	// check user's alarm status
	if !account.AlarmStatus {
		return false
	}

	return true
}
