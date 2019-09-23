package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/models/types"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/utils"

	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetDelegatorWithdrawAddress returns delegator's reward withdraw address
func GetDelegatorWithdrawAddress(config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	delegatorAddr := vars["delegatorAddr"]

	// check the validity of account
	if !strings.Contains(delegatorAddr, sdk.Bech32PrefixAccAddr) || len(delegatorAddr) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// delegator's withdraw_address
	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + delegatorAddr + "/withdraw_address")

	var responseWithHeight types.ResponseWithHeight
	err := json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal responseWithHeight error - %v\n", err)
	}

	// Unmarshal struct
	var address string
	err = json.Unmarshal(responseWithHeight.Result, &address)
	if err != nil {
		fmt.Printf("unmarshal distribution/delegators/{delegatorAddr}/withdraw_address pool error - %v\n", err)
	}

	result := make(map[string]string)
	result["withdraw_address"] = address

	utils.Respond(w, result)
	return nil
}

// GetDelegatorRewards returns a withdrawn delegation rewards
func GetDelegatorRewards(config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	delegatorAddr := vars["delegatorAddr"]
	validatorAddr := vars["validatorAddr"]

	// check the validity of account & validator address
	if !strings.Contains(delegatorAddr, sdk.Bech32PrefixAccAddr) || !strings.Contains(validatorAddr, sdk.Bech32PrefixValAddr) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query a delegation reward
	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + delegatorAddr + "/rewards/" + validatorAddr)

	var responseWithHeight types.ResponseWithHeight
	err := json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal responseWithHeight error - %v\n", err)
	}

	// Unmarshal struct
	var coin []models.Coin
	err = json.Unmarshal(responseWithHeight.Result, &coin)
	if err != nil {
		fmt.Printf("unmarshal distribution/community_pool pool error - %v\n", err)
	}

	utils.Respond(w, coin)
	return nil
}

// GetCommunityPool returns current community pool
func GetCommunityPool(config *config.Config, db *pg.DB, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	// query stake pool - bonded and not bonded tokens
	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/community_pool")

	var responseWithHeight types.ResponseWithHeight
	err := json.Unmarshal(resp.Body(), &responseWithHeight)
	if err != nil {
		fmt.Printf("unmarshal responseWithHeight error - %v\n", err)
	}

	var coin []models.Coin
	err = json.Unmarshal(responseWithHeight.Result, &coin)
	if err != nil {
		fmt.Printf("unmarshal distribution/community_pool pool error - %v\n", err)
	}

	utils.Respond(w, coin)
	return nil
}
