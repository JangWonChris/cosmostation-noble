package controllers

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/services"

	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// StatsController passes requests to its respective service
func StatsController(codec *codec.Codec, config *config.Config, db *db.Database, r *mux.Router, rpcClient *client.HTTP) {
	r.HandleFunc("/stats/market", func(w http.ResponseWriter, r *http.Request) {
		services.GetMarketStats(config, db, rpcClient, w, r)
	})

	r.HandleFunc("/stats/network", func(w http.ResponseWriter, r *http.Request) {
		services.GetNetworkStats(config, db, rpcClient, w, r)
	})
}
