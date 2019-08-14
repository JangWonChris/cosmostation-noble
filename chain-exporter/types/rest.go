package types

import "encoding/json"

type ResponseWithHeight struct {
	Height int64           `json:"height"`
	Result json.RawMessage `json:"result"`
}
