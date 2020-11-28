package client

import (
	"fmt"
	"net/http"
)

// RequestWithRestServer is general request API from REST Server and
// return without any modification
func (c *Client) RequestWithRestServer(reqParam string) ([]byte, error) {
	// deprecated function, do not use for querying any with rest-server
	resp, err := c.apiClient.R().Get(reqParam)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get respond : %s", resp.Status())
	}

	return resp.Body(), nil
}
