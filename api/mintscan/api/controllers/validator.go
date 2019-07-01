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
func ValidatorController(r *mux.Router, RPCClient *client.HTTP, DB *pg.DB, Codec *codec.Codec, Config *config.Config) {
	r.HandleFunc("/staking/validators", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidators(RPCClient, DB, w, r)
	})
	r.HandleFunc("/staking/validator/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidator(RPCClient, DB, w, r)
	})
	r.HandleFunc("/staking/validator/misses/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorBlockMisses(RPCClient, DB, w, r)
	})
	r.HandleFunc("/staking/validator/misses/detail/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorBlockMissesDetail(RPCClient, DB, w, r)
	})
	r.HandleFunc("/staking/validator/events/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorEvents(DB, w, r)
	})
	r.HandleFunc("/staking/redelegations", func(w http.ResponseWriter, r *http.Request) {
		services.GetRedelegations(DB, Config, w, r)
	})
	// Currently not used due to Full Node requests performance issue
	r.HandleFunc("/staking/validator/delegations/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetValidatorDelegations(Codec, RPCClient, DB, Config, w, r)
	})
}
