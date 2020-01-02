package alarm

import (
	"strings"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/databases"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"

	"github.com/go-pg/pg"

	resty "gopkg.in/resty.v1"
)

// PushNotification sends push notification to its respective device
func PushNotification(cf *config.Config, db *pg.DB, pnp types.PushNotificationPayload) {

	// query account information
	var fromAccount types.Account
	var toAccount types.Account

	// return when data is empty
	if account.AlarmToken == "" {
		return
	}

	// check user's alarm status
	if !account.AlarmStatus {
		return
	}

	// push notification payload
	var pns []types.PushNotifications
	tempNotification := types.PushNotifications{
		Tokens:   []string{account.AlarmToken},
		Platform: 2,
		Title:    types.PushNotificationReceivedTitle + pnp.Amount,
		Message:  types.PushNotificationReceivedMessage + pnp.Amount,
		Data: types.PushNotificationData{
			NotifyTo: pnp.To,
			Txid:     pnp.Txid,
		},
	}
	pns = append(pns, tempNotification)

	pnsp := types.PushNotificationServerPayload{
		Notifications: pns,
	}

	// send push notification
	_, err = resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(pnsp).
		Post(cf.Alarm.PushServerURL)
	if err != nil {
		return
	}

}

// VerifyAccount verifes account before sending push notification
func VerifyAccount(pnp types.PushNotificationPayload) {
	cf := config.NewConfig()
	db := databases.Connect(cf)

	defer db.Close()

	fromAccount, _ := db.QueryAccount(db, pnp.From)
	toAccount, _ := db.QueryAccount(db, pnp.To)

	pnp.From = strings.ToLower(pnp.From)
	pnp.To = strings.ToLower(pnp.To)
	pnp.Txid = strings.ToLower(pnp.Txid)

}
