package controllers

import (
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tomasen/realip"
)

// Passes requests to its respective service
func AccountController(codec *codec.Codec, config *config.Config, db *pg.DB, router *mux.Router, rpcClient *client.HTTP) {
	router.HandleFunc("/account/{address}", func(w http.ResponseWriter, r *http.Request) {
		// TEST
		clientIP := realip.FromRequest(r) // FromRequest return client's real public IP address from http request headers.
		log.Println("GET /account/{address}: ", clientIP)

		services.GetAccountInfo(codec, config, db, rpcClient, w, r)
	}).Methods("GET")
}
