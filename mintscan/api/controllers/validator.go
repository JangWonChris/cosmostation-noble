package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/codec"
)

// ValidatorController passes requests to its respective service
func ValidatorController(codec *codec.Codec, config *config.Config, db *db.Database, r *mux.Router, rpcClient *client.HTTP) {
	r.HandleFunc("/staking/validators", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidators(db, rpcClient, w, r)
	})
	r.HandleFunc("/staking/validator/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidator(db, rpcClient, w, r)
	})
	r.HandleFunc("/staking/validator/misses/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorBlockMisses(db, rpcClient, w, r)
	})
	r.HandleFunc("/staking/validator/misses/detail/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorBlockMissesDetail(db, rpcClient, w, r)
	})
	r.HandleFunc("/staking/validator/events/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorEvents(db, w, r)
	})
	r.HandleFunc("/staking/validator/events/{address}/count", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorEventsTotalCount(db, w, r)
	})
	r.HandleFunc("/staking/redelegations", func(w http.ResponseWriter, r *http.Request) {
		services.GetRedelegations(config, db, w, r)
	})
	r.HandleFunc("/staking/validator/delegations/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorDelegations(codec, config, db, rpcClient, w, r)
	})
}
