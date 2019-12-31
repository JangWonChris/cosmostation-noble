package controllers

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/services"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// Passes requests to its respective service
func AccountController(codec *codec.Codec, config *config.Config, db *pg.DB, r *mux.Router, rpcClient *client.HTTP) {
	r.HandleFunc("/account/balance/{accAddress}", func(w http.ResponseWriter, r *http.Request) {
		services.GetBalance(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	r.HandleFunc("/account/delegations/rewards/{accAddress}", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegationsRewards(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	r.HandleFunc("/account/delegations/{accAddress}", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegations(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	r.HandleFunc("/account/commission/{accAddress}", func(w http.ResponseWriter, r *http.Request) {
		services.GetCommission(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	r.HandleFunc("/account/unbonding-delegations/{accAddress}", func(w http.ResponseWriter, r *http.Request) {
		services.GetUnbondingDelegations(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
}
