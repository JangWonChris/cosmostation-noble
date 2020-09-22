package notification

import (
	"time"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/types"
	"go.uber.org/zap"

	resty "github.com/go-resty/resty/v2"
)

// Notification wraps around configuration for push notification for mobile apps.
type Notification struct {
	notiClient *resty.Client
}

// NewNotification returns new notification instance.
func NewNotification() *Notification {
	config := config.ParseConfig()

	client := resty.New().
		SetHostURL(config.Alarm.PushServerEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	return &Notification{client}
}

// Push sends push notification to local notification server and it delivers the message to
// its respective device. Uses a push notification micro server called gorush.
// More information can be found here in this link. https://github.com/appleboy/gorush
func (nof *Notification) Push(np types.NotificationPayload, token string, target string) {
	var notifications []types.Notification

	// Create new notification payload for a user sending tokens
	if target == types.From {
		platform := int8(2)
		title := types.NotificationSentTitle + np.Amount + np.Denom
		message := types.NotificationSentMessage + np.Amount + np.Denom

		data := types.NewNotificationData(np.From, np.Txid, types.Sent)
		payload := types.NewNotification([]string{token}, platform, title, message, data)

		notifications = append(notifications, payload)

		zap.S().Info("send - push notification")
		zap.S().Infof("hash: %s | from: %s", np.Txid, np.From)
	}

	// Create new notification payload for a user receiving tokens
	if target == types.To {
		platform := int8(2)
		title := types.NotificationReceivedTitle + np.Amount + np.Denom
		message := types.NotificationReceivedMessage + np.Amount + np.Denom

		data := types.NewNotificationData(np.To, np.Txid, types.Received)
		payload := types.NewNotification([]string{token}, platform, title, message, data)

		notifications = append(notifications, payload)

		zap.S().Info("send - push notification")
		zap.S().Infof("hash: %s | to: %s", np.Txid, np.To)
	}

	if len(notifications) <= 0 {
		return
	}

	nsp := types.NotificationServerPayload{
		Notifications: notifications,
	}

	// Send push notification
	resp, err := nof.notiClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(nsp).Post("")

	if resp.IsError() {
		zap.S().Debugf("failed to send push notification: %s", err)
	}
}
