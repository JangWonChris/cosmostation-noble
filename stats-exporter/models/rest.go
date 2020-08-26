package models

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

// ResponseWithHeight is a wrapper for returned values from REST API calls.
type ResponseWithHeight struct {
	Height string          `json:"height"`
	Result json.RawMessage `json:"result"`
}

// ReadRespWithHeight reads response with height that has been changed in REST APIs since v0.36.+.
func ReadRespWithHeight(resp *resty.Response) ResponseWithHeight {
	var responseWithHeight ResponseWithHeight
	err := json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("failed to unmarshal response with height: %s", err)
		return ResponseWithHeight{}
	}

	return responseWithHeight
}
