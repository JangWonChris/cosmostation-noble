package models

import "time"

type Account struct {
	IdfAccount uint8     `json:"idf_account,omitempty" sql:",pk"`
	AlarmToken string    `json:"alarm_token,omitempty"`
	DeviceType string    `json:"device_type,omitempty"`
	Address    string    `json:"address,omitempty"`
	Nickname   string    `json:"nickname,omitempty"`
	CoinType   string    `json:"coin_type,omitempty"`
	Status     bool      `json:"status,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}
