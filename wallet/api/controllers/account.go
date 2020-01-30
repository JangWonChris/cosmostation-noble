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
	r.HandleFunc("/account/update", func(w http.ResponseWriter, r *http.Request) {
		services.RegisterOrUpdate(db, w, r)
	}).Methods("POST")
	r.HandleFunc("/account/delete", func(w http.ResponseWriter, r *http.Request) {
		services.Delete(db, w, r)
	}).Methods("DELETE")
}
