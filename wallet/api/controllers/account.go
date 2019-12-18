package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// AccountController passes requests to its respective service
func AccountController(r *mux.Router, c *client.HTTP, db *pg.DB) {
	r.HandleFunc("/account/register", func(w http.ResponseWriter, r *http.Request) {
		services.Register(db, w, r)
	}).Methods("POST")
}
