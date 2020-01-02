package types

var (
	PushNotificationSentTitle   = "Sent "
	PushNotificationSentMessage = "You have just sent "

	PushNotificationReceivedTitle   = "Received "
	PushNotificationReceivedMessage = "You have just received "
)

type PushNotificationPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Txid   string `json:"txid"`
	Amount string `json:"amount"`
}

type PushNotificationServerPayload struct {
	Notifications []PushNotifications `json:"notifications"`
}

type PushNotifications struct {
	Tokens   []string             `json:"tokens"`
	Platform int8                 `json:"platform"`
	Title    string               `json:"title"`
	Message  string               `json:"message"`
	Data     PushNotificationData `json:"data"`
}

type PushNotificationData struct {
	NotifyTo string `json:"notifyto"`
	Txid     string `json:"txid"`
	Type     string `json:"type"`
}
