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
func StatsController(r *mux.Router, RPCClient *client.HTTP, DB *pg.DB, Config *config.Config) {
	r.HandleFunc("/stats/market", func(w http.ResponseWriter, r *http.Request) {
		services.GetMarketInfo(RPCClient, DB, Config, w, r)
	})

	r.HandleFunc("/stats/network", func(w http.ResponseWriter, r *http.Request) {
		services.GetNetworkStats(RPCClient, DB, Config, w, r)
	})
}
