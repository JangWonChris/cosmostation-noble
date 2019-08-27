package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Passes requests to its respective service
func ValidatorController(codec *codec.Codec, config *config.Config, db *pg.DB, router *mux.Router, rpcClient *client.HTTP) {
	router.HandleFunc("/staking/validators", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidators(db, rpcClient, w, r)
	})
	router.HandleFunc("/staking/validator/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidator(db, rpcClient, w, r)
	})
	router.HandleFunc("/staking/validator/misses/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorBlockMisses(db, rpcClient, w, r)
	})
	router.HandleFunc("/staking/validator/misses/detail/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorBlockMissesDetail(db, rpcClient, w, r)
	})
	router.HandleFunc("/staking/validator/events/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorEvents(db, w, r)
	})
	router.HandleFunc("/staking/redelegations", func(w http.ResponseWriter, r *http.Request) {
		services.GetRedelegations(config, db, w, r)
	})
	router.HandleFunc("/staking/validator/delegations/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorDelegations(codec, config, db, rpcClient, w, r)
	})
}
