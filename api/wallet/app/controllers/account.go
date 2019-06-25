package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// AccountController passes requests to its respective service
func AccountController(r *mux.Router, c *client.HTTP, DB *pg.DB) {
	r.HandleFunc("/account/register", func(w http.ResponseWriter, r *http.Request) {
		services.Register(DB, w, r)
	}).Methods("POST")
}
