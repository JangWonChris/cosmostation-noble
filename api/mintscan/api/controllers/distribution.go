package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// Passes requests to its respective service
func DistributionController(Config *config.Config, DB *pg.DB, r *mux.Router, RPCClient *client.HTTP) {
	r.HandleFunc("/distribution/delegators/{delegatorAddr}/withdraw_address", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegatorWithdrawAddress(Config, DB, RPCClient, w, r)
	})
	r.HandleFunc("/distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegatorRewards(Config, DB, RPCClient, w, r)
	})
	r.HandleFunc("/distribution/community_pool", func(w http.ResponseWriter, r *http.Request) {
		services.GetCommunityPool(Config, DB, RPCClient, w, r)
	})
}
