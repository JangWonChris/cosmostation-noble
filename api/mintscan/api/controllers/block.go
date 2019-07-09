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
func BlockController(codec *codec.Codec, config *config.Config, db *pg.DB, router *mux.Router, rpcClient *client.HTTP) {
	router.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		services.GetBlocks(db, w, r)
	})
	router.HandleFunc("/blocks/{address}", func(w http.ResponseWriter, r *http.Request) {
		services.GetProposedBlocks(db, w, r)
	})
}
