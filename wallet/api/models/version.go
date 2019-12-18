package models

import "time"

type AppVersion struct {
	IdfVersion uint16    `json:"idf_version,omitempty" sql:",pk"`
	AppName    string    `json:"app_name,omitempty"`
	DeviceType string    `json:"device_type,omitempty"`
	Accpetable uint16    `json:"acceptable,omitempty"`
	Latest     uint16    `json:"latest,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
}
