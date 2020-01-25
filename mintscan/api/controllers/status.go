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

// StatusController passes requests to its respective service
func StatusController(codec *codec.Codec, config *config.Config, db *db.Database, r *mux.Router, rpcClient *client.HTTP) {
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		services.GetStatus(config, db, rpcClient, w, r)
	})
}
