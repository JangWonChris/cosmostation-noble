package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// Passes requests to its respective service
func MintingController(r *mux.Router, RPCClient *client.HTTP, DB *pg.DB, Config *config.Config) {
	r.HandleFunc("/minting/inflation", func(w http.ResponseWriter, r *http.Request) {
		services.GetMintingInflation(RPCClient, DB, Config, w, r)
	}).Methods("GET")
}
