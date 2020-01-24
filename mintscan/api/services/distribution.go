package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetDelegatorWithdrawAddress returns delegator's reward withdraw address
func GetDelegatorWithdrawAddress(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	delegatorAddr := vars["delegatorAddr"]

	if !utils.VerifyAddress(delegatorAddr) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query delegator's withdraw address
	resp, _ := resty.R().Get(config.Node.LCDEndpoint + "/distribution/delegators/" + delegatorAddr + "/withdraw_address")

	var address string
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &address)
	if err != nil {
		fmt.Printf("failed to unmarshal address: %t\n", err)
	}

	result := make(map[string]string)
	result["withdraw_address"] = address

	utils.Respond(w, result)
	return nil
}

// GetDelegatorRewards returns a withdrawn delegation rewards
func GetDelegatorRewards(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	delegatorAddr := vars["delegatorAddr"]
	validatorAddr := vars["validatorAddr"]

	if !utils.VerifyAddress(delegatorAddr) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	if !utils.VerifyValAddress(delegatorAddr) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// Query a delegation reward
	resp, _ := resty.R().Get(config.Node.LCDEndpoint + "/distribution/delegators/" + delegatorAddr + "/rewards/" + validatorAddr)

	coin := make([]models.Coin, 0)
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &coin)
	if err != nil {
		fmt.Printf("failed to unmarshal coin: %t\n", err)
	}

	utils.Respond(w, coin)
	return nil
}

// GetCommunityPool returns current community pool
func GetCommunityPool(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	resp, _ := resty.R().Get(config.Node.LCDEndpoint + "/distribution/community_pool")

	var coin []models.Coin
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &coin)
	if err != nil {
		fmt.Printf("failed to unmarshal coin: %t\n", err)
	}

	utils.Respond(w, coin)
	return nil
}
