package controllers

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// Passes requests to its respective service
func DistributionController(codec *codec.Codec, config *config.Config, db *pg.DB, router *mux.Router, rpcClient *client.HTTP) {
	router.HandleFunc("/distribution/delegators/{delegatorAddr}/withdraw_address", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegatorWithdrawAddress(config, db, rpcClient, w, r)
	})
	router.HandleFunc("/distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}", func(w http.ResponseWriter, r *http.Request) {
		services.GetDelegatorRewards(config, db, rpcClient, w, r)
	})
	router.HandleFunc("/distribution/community_pool", func(w http.ResponseWriter, r *http.Request) {
		services.GetCommunityPool(config, db, rpcClient, w, r)
	})
}
