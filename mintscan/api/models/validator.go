package models

type Uptime struct {
	Address      string `json:"address"`
	MissedBlocks int64  `json:"missed_blocks"`
	OverBlocks   int64  `json:"over_blocks"`
}
