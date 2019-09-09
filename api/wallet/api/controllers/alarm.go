package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/api/services"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// AlarmController passes requests to its respective service
func AlarmController(r *mux.Router, c *client.HTTP, db *pg.DB) {
	r.HandleFunc("/account/alarm/update", func(w http.ResponseWriter, r *http.Request) {
		services.UpdateAlarmStatus(db, w, r)
	}).Methods("POST")
}
