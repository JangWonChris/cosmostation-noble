package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// Passes requests to its respective service
func TransactionController(r *mux.Router, RPCClient *client.HTTP, DB *pg.DB, Codec *codec.Codec, Config *config.Config) {
	r.HandleFunc("/txs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.GetTxs(Codec, RPCClient, DB, w, r)
	})
	r.HandleFunc("/tx/{hash}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.GetTx(Codec, RPCClient, DB, Config, w, r)
	})
	r.HandleFunc("/tx/broadcast/{hash}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.BroadcastTx(Codec, RPCClient, w, r)
	})
}
