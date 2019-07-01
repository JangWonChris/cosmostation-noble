package controllers

import (
	"log"
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tomasen/realip"
)

// Passes requests to its respective service
func AccountController(r *mux.Router, c *client.HTTP, DB *pg.DB, Config *config.Config) {
	r.HandleFunc("/account/{address}", func(w http.ResponseWriter, r *http.Request) {
		// TEST
		clientIP := realip.FromRequest(r) // FromRequest return client's real public IP address from http request headers.
		log.Println("GET /account/{address}: ", clientIP)

		services.GetAccountInfo(DB, Config, w, r)
	}).Methods("GET")
}
