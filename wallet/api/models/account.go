package models

import "time"

// Account is an account for our cosmostation mobile wallet app users
type Account struct {
	IdfAccount  uint16    `json:"idf_account,omitempty" sql:",pk"`
	ChainID     uint16    `json:"chain_id,omitempty"`
	Address     string    `json:"address,omitempty"`
	AlarmToken  string    `json:"alarm_token,omitempty"`
	DeviceType  string    `json:"device_type,omitempty"`
	AlarmStatus bool      `json:"alarm_status,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
}
