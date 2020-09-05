package schema

import "time"

// AppAccount defines the structure for Cosmostation's wallet app users information.
type AppAccount struct {
	IdfAccount  uint16    `json:"idf_account" sql:",pk"`
	ChainID     uint16    `json:"chain_id,omitempty" sql:",notnull"`
	DeviceType  string    `json:"device_type,omitempty" sql:",notnull"`
	Address     string    `json:"address" sql:",notnull"`
	AlarmToken  string    `json:"alarm_token" sql:",notnull"`
	AlarmStatus bool      `json:"alarm_status" sql:",notnull"`
	Timestamp   time.Time `json:"timestamp,omitempty" sql:"default:now()"`
}
