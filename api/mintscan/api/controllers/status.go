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
func StatusController(r *mux.Router, RPCClient *client.HTTP, DB *pg.DB, Config *config.Config) {
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.GetStatus(RPCClient, DB, Config, w, r)
	})
}
