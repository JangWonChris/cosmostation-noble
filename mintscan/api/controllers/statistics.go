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
func StatsController(codec *codec.Codec, config *config.Config, db *pg.DB, router *mux.Router, rpcClient *client.HTTP) {
	router.HandleFunc("/stats/market", func(w http.ResponseWriter, r *http.Request) {
		services.GetMarketStats(config, db, rpcClient, w, r)
	})

	router.HandleFunc("/stats/network", func(w http.ResponseWriter, r *http.Request) {
		services.GetNetworkStats(config, db, rpcClient, w, r)
	})
}
