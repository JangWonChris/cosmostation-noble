package models

import "time"

// Account is an account for our cosmostation mobile wallet app users
type Account struct {
	IdfAccount  uint16    `json:"idf_account" sql:",pk"`
	ChainID     uint16    `json:"chain_id" sql:",notnull"`
	DeviceType  string    `json:"device_type" sql:",notnull"`
	Address     string    `json:"address" sql:",unique, notnull"`
	AlarmToken  string    `json:"alarm_token" sql:",notnull"`
	AlarmStatus bool      `json:"alarm_status" sql:",notnull"`
	Timestamp   time.Time `json:"timestamp,omitempty" sql:"default:now()"`
}
