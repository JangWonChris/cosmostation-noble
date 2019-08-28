package models

import "encoding/json"

// ResponseWithHeight is a wrapper for returned values from REST API calls
type ResponseWithHeight struct {
	Height string          `json:"height"`
	Result json.RawMessage `json:"result"`
}
