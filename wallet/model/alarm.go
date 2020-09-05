package model

var (
	AlarmTitle   = "Received "
	AlarmMessage = "You have just received "
)

type NotificationReceivePayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Txid   string `json:"txid"`
	Amount string `json:"amount"`
}

type NotificationPayload struct {
	Notifications []Notifications `json:"notifications"`
}

type Notifications struct {
	Tokens   []string `json:"tokens"`
	Platform int8     `json:"platform"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`
	Data     Data     `json:"data"`
}

type Data struct {
	NotifyTo string `json:"notifyto"`
	Txid     string `json:"txid"`
	Type     string `json:"type"`
}
