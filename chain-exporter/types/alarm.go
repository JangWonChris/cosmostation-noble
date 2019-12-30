package types

type NotificationPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Txid   string `json:"txid"`
	Amount string `json:"amount"`
}
