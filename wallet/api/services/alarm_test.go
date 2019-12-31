package services

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"
)

func TestNotificationPayload(t *testing.T) {
	var notification []models.Notifications
	tempNotification := models.Notifications{
		Tokens:   []string{"APA91bEEVWV7AbszQutnuAlFZfn9aXucZUCo_sbTltmKB7F1_l3n2TtlR31HmPx04xSw6kl0V0Fafjn4koAqydPMR8heKv8n_9Zr0bLIjsjVzfOFXC2jjWbfTxhURSDTnW0_Zvh1s6J5"},
		Platform: 2,
		Message:  "Hello World",
		Title:    "Title",
		Data: models.Data{
			NotifyTo: "cosmos1ma02nlc7lchu7caufyrrqt4r6v2mpsj92s3mw7",
			Txid:     "06FA072B36E4D9D0E99C9BAA826794DE11109F697916F3B0A93FCA8919754827",
		},
	}
	notification = append(notification, tempNotification)

	notificationPayload := models.NotificationPayload{
		Notifications: notification,
	}

	result, err := json.Marshal(notificationPayload)
	if err != nil {
		t.Errorf("error whhen marshal notificationPayload %s", err)
	}

	fmt.Println(string(result))
}
