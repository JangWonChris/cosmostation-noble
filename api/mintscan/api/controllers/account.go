package controllers

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// Passes requests to its respective service
func AccountController(codec *codec.Codec, config *config.Config, db *pg.DB, router *mux.Router, rpcClient *client.HTTP) {
	// router.HandleFunc("/account/{address}", func(w http.ResponseWriter, r *http.Request) {
	// 	// TEST
	// 	clientIP := realip.FromRequest(r) // FromRequest return client's real public IP address from http request headers.
	// 	log.Println("GET /account/{address}: ", clientIP)

	// 	services.GetAccountInfo(codec, config, db, rpcClient, w, r)
	// }).Methods("GET")
	router.HandleFunc("/account/balance/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetBalance(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	router.HandleFunc("/account/delegations/rewards/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegationsRewards(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	router.HandleFunc("/account/delegations/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegations(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	router.HandleFunc("/account/commission/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetCommission(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
	router.HandleFunc("/account/unbonding-delegations/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetUnbondingDelegations(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
}
