package models

type PushNotification struct {
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	Tokens   []string `json:"tokens"`
	Platform int8     `json:"platform"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`
	Data     Payload  `json:"data"`
}

type Payload struct {
	NotifyTo string `json:"notifyto"`
	Txid     string `json:"txid"`
}
