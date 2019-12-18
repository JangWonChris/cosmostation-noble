package models

import "time"

const (
	Android = "android"
	IOS     = "ios"
)

type AppVersion struct {
	IdfVersion uint16    `json:"idf_version,omitempty" sql:",pk"`
	AppName    string    `json:"app_name,omitempty"`
	DeviceType string    `json:"device_type"`
	Acceptable uint16    `json:"acceptable,omitempty"`
	Latest     uint16    `json:"latest"`
	Timestamp  time.Time `json:"timestamp" sql:"default:now()"`
}
