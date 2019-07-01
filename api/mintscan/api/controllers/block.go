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
func BlockController(r *mux.Router, RPCClient *client.HTTP, DB *pg.DB, Config *config.Config) {
	r.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		services.GetBlocks(DB, w, r)
	})
	r.HandleFunc("/blocks/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetProposedBlocks(DB, w, r)
	})
}
