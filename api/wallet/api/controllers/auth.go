package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// AuthController passes requests to its respective service
func AuthController(r *mux.Router, c *client.HTTP, DB *pg.DB) {
	r.HandleFunc("/auth/account/test", func(w http.ResponseWriter, r *http.Request) {
		services.Test(DB, w, r)
	}).Methods("POST")
}
