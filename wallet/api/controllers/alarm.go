package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/services"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// AlarmController passes requests to its respective service
func AlarmController(r *mux.Router, c *client.HTTP, db *pg.DB) {
	r.HandleFunc("/alarm/push", func(w http.ResponseWriter, r *http.Request) {
		services.PushNotification(db, w, r)
	}).Methods("POST")
}
