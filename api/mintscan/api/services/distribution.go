package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"

	resty "gopkg.in/resty.v1"
)

// GetDelegatorWithdrawAddress returns delegator's reward withdraw address
func GetDelegatorWithdrawAddress(Config *config.Config, DB *pg.DB, RPCClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	delegatorAddr := vars["delegatorAddr"]

	// Check the validity of cosmos address
	if !strings.Contains(delegatorAddr, sdk.Bech32PrefixAccAddr) || len(delegatorAddr) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query LCD - delegator's withdraw_address
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(Config.Node.LCDURL + "/distribution/delegators/" + delegatorAddr + "/withdraw_address")

	// Unmarshal struct
	var address string
	err := json.Unmarshal(resp.Body(), &address)
	if err != nil {
		fmt.Printf("distribution/delegators/{delegatorAddr}/withdraw_address unmarshal pool error - %v\n", err)
	}

	// Result
	result := make(map[string]string)
	result["withdraw_address"] = address

	return json.NewEncoder(w).Encode(result)
}

// GetDelegatorRewards returns a withdrawn delegation rewards
func GetDelegatorRewards(Config *config.Config, DB *pg.DB, RPCClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	// Receive address
	vars := mux.Vars(r)
	delegatorAddr := vars["delegatorAddr"]
	validatorAddr := vars["validatorAddr"]

	// Check the validity of cosmos address & validator address
	if !strings.Contains(delegatorAddr, sdk.Bech32PrefixAccAddr) ||
		!strings.Contains(validatorAddr, sdk.Bech32PrefixValAddr) ||
		len(delegatorAddr) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query LCD - Query a delegation reward
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(Config.Node.LCDURL + "/distribution/delegators/" + delegatorAddr + "/rewards/" + validatorAddr)

	// Unmarshal struct
	var coin []models.Coin
	err := json.Unmarshal(resp.Body(), &coin)
	if err != nil {
		fmt.Printf("distribution/community_pool unmarshal pool error - %v\n", err)
	}

	return json.NewEncoder(w).Encode(coin)
}

// GetCommunityPool returns current community pool
func GetCommunityPool(Config *config.Config, DB *pg.DB, RPCClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	// Query LCD - stake pool to get bonded and unbonded tokens
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, _ := resty.R().Get(Config.Node.LCDURL + "/distribution/community_pool")

	// Unmarshal struct
	var coin []models.Coin
	err := json.Unmarshal(resp.Body(), &coin)
	if err != nil {
		fmt.Printf("distribution/community_pool unmarshal pool error - %v\n", err)
	}

	return json.NewEncoder(w).Encode(coin)
}
