package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// Passes requests to its respective service
func TransactionController(codec *codec.Codec, config *config.Config, db *pg.DB, router *mux.Router, rpcClient *client.HTTP) {
	router.HandleFunc("/txs", func(w http.ResponseWriter, r *http.Request) {
		services.GetTxs(codec, db, rpcClient, w, r)
	})
	router.HandleFunc("/tx/{hash}", func(w http.ResponseWriter, r *http.Request) {
		services.GetTx(codec, config, db, rpcClient, w, r)
	})
	router.HandleFunc("/tx/broadcast/{hash}", func(w http.ResponseWriter, r *http.Request) {
		services.BroadcastTx(codec, rpcClient, w, r)
	})
}
