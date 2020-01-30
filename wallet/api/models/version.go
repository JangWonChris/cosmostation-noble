package models

import "time"

const (
	Android     = "android"
	IOS         = "ios"
	ForceUpdate = false
)

type AppVersion struct {
	IdfVersion uint16    `json:"idf_version,omitempty" sql:",pk"`
	AppName    string    `json:"app_name"`
	DeviceType string    `json:"device_type"`
	Version    uint16    `json:"version"`
	Enable     bool      `json:"enable" sql:"default:false"`
	Timestamp  time.Time `json:"timestamp" sql:"default:now()"`
}
