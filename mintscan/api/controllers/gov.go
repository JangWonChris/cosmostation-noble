package controllers

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/gorilla/mux"
)

// GovController passes requests to its respective service
func GovController(codec *codec.Codec, config *config.Config, db *db.Database, r *mux.Router, rpcClient *client.HTTP) {
	r.HandleFunc("/gov/proposals", func(w http.ResponseWriter, r *http.Request) {
		services.GetProposals(db, config, w, r)
	})
	r.HandleFunc("/gov/proposal/{proposalId}", func(w http.ResponseWriter, r *http.Request) {
		services.GetProposal(db, config, w, r)
	})
	r.HandleFunc("/gov/proposal/votes/{proposalId}", func(w http.ResponseWriter, r *http.Request) {
		services.GetVotes(db, config, w, r)
	})
	r.HandleFunc("/gov/proposal/deposits/{proposalId}", func(w http.ResponseWriter, r *http.Request) {
		services.GetDeposits(db, w, r)
	})
}
