package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	resty "github.com/go-resty/resty/v2"
)

const (
	// DefaultLimit is the default limit of response items for APIs.
	DefaultLimit = 100

	// DefaultPowerEventHistoryLimit is the default limit for Power Event History API.
	DefaultPowerEventHistoryLimit = 50

	// DefaultBefore is the default parameter that will be used to query database previous items from the default param.
	DefaultBefore = 0

	// DefaultAfter is he default parameter that will be used to query database next items from the default param.
	DefaultAfter = -1
)

// ResponseWithHeight is a wrapper for returned values from REST API calls.
type ResponseWithHeight struct {
	Height string          `json:"height"`
	Result json.RawMessage `json:"result"`
}

// ReadRespWithHeight reads response with height that has been changed in REST APIs from v0.36.0.
func ReadRespWithHeight(resp *resty.Response) ResponseWithHeight {
	var responseWithHeight ResponseWithHeight
	err := json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal responseWithHeight error - %s\n", err)
	}
	return responseWithHeight
}

// ParseHTTPArgsWithBeforeAfterLimit parses the request's URL and returns all arguments pairs.
// It separates page and limit used for pagination where a default limit can be provided.
func ParseHTTPArgsWithBeforeAfterLimit(r *http.Request, defaultBefore, defaultAfter, defaultLimit int) (before, after, limit int, err error) {
	beforeStr := r.FormValue("before")
	if beforeStr == "" {
		before = defaultBefore
	} else {
		before, err = strconv.Atoi(beforeStr)
		if err != nil {
			return before, after, limit, errors.New("failed to convert to integer type")
		}

		if before < 0 {
			return before, after, limit, errors.New("before param must be equal or greater than 0")
		}
	}

	afterStr := r.FormValue("after")
	if afterStr == "" {
		after = defaultAfter
	} else {
		after, err = strconv.Atoi(afterStr)
		if err != nil {
			return before, after, limit, errors.New("failed to convert to integer type")
		}

		if after < 0 {
			return before, after, limit, errors.New("after param must be equal or greater than 0")
		}
	}

	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return before, after, limit, errors.New("failed to convert to integer type")
		}

		if limit <= 0 {
			return before, after, limit, errors.New("limit param must be greater than 0")
		}
	}

	return before, after, limit, nil
}
