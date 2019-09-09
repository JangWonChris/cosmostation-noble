package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// VersionController passes requests to its respective service
func VersionController(r *mux.Router, c *client.HTTP, db *pg.DB) {
	r.HandleFunc("/app/version/{deviceType}", func(w http.ResponseWriter, r *http.Request) {
		services.GetVersion(db, w, r)
	}).Methods("GET")
}
