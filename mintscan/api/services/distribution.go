package services

import (
	"encoding/json"
	"net/http"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/errors"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/models"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/utils"

	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/rpc/client"
	resty "gopkg.in/resty.v1"
)

// GetDelegatorWithdrawAddress returns delegator's reward withdraw address
func GetDelegatorWithdrawAddress(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	delegatorAddr := vars["delegatorAddr"]

	if !strings.Contains(delegatorAddr, sdk.Bech32PrefixAccAddr) || len(delegatorAddr) != 45 {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// delegator's withdraw_address
	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + delegatorAddr + "/withdraw_address")

	var address string
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &address)
	if err != nil {
		log.Info().Str(models.Service, models.LogDistribution).Str(models.Method, "GetDelegatorWithdrawAddress").Err(err).Msg("unmarshal address error")
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

	if !strings.Contains(delegatorAddr, sdk.Bech32PrefixAccAddr) || !strings.Contains(validatorAddr, sdk.Bech32PrefixValAddr) {
		errors.ErrNotExist(w, http.StatusNotFound)
		return nil
	}

	// query a delegation reward
	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/delegators/" + delegatorAddr + "/rewards/" + validatorAddr)

	coin := make([]models.Coin, 0)

	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &coin)
	if err != nil {
		log.Info().Str(models.Service, models.LogDistribution).Str(models.Method, "GetDelegatorRewards").Err(err).Msg("unmarshal coin error")
	}

	utils.Respond(w, coin)
	return nil
}

// GetCommunityPool returns current community pool
func GetCommunityPool(config *config.Config, db *db.Database, rpcClient *client.HTTP, w http.ResponseWriter, r *http.Request) error {
	resp, _ := resty.R().Get(config.Node.LCDURL + "/distribution/community_pool")

	var coin []models.Coin
	err := json.Unmarshal(models.ReadRespWithHeight(resp).Result, &coin)
	if err != nil {
		log.Info().Str(models.Service, models.LogDistribution).Str(models.Method, "GetCommunityPool").Err(err).Msg("unmarshal coin error")
	}

	utils.Respond(w, coin)
	return nil
}
