package types

// These are parameters that are used for sending push notification to push notification server.
var (
	From = "from"
	To   = "to"

	Sent     = "sent"
	Received = "received"

	NotificationSentTitle   = "Sent "
	NotificationSentMessage = "You have just sent "

	NotificationReceivedTitle   = "Received "
	NotificationReceivedMessage = "You have just received "
)

// NotificationServerPayload defines the structure for a list of payloads from push notification server.
type NotificationServerPayload struct {
	Notifications []Notification `json:"notifications"`
}

// Notification defines the structure for payload from push notification server.
type Notification struct {
	Tokens   []string         `json:"tokens"`
	Platform int8             `json:"platform"`
	Title    string           `json:"title"`
	Message  string           `json:"message"`
	Data     NotificationData `json:"data"`
}

// NotificationData defines the structure for notification data.
type NotificationData struct {
	NotifyTo string `json:"notifyto"`
	Txid     string `json:"txid"`
	Type     string `json:"type"`
}

// NotificationPayload defines the structure for notification payload.
type NotificationPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Txid   string `json:"txid"`
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

// NewNotification returns new Notification.
func NewNotification(tokens []string, platform int8, title string, message string, data NotificationData) Notification {
	return Notification{tokens, platform, title, message, data}
}

// NewNotificationPayload returns new NotificationPayload.
func NewNotificationPayload(payload NotificationPayload) *NotificationPayload {
	return &NotificationPayload{
		From:   payload.From,
		To:     payload.To,
		Txid:   payload.Txid,
		Amount: payload.Amount,
		Denom:  payload.Denom,
	}
}

// NewNotificationData returns new NotificationData.
func NewNotificationData(notifyto string, txid string, notifyType string) NotificationData {
	return NotificationData{notifyto, txid, notifyType}
}
