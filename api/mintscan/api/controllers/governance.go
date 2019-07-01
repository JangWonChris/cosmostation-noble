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
func GovernanceController(r *mux.Router, RPCClient *client.HTTP, DB *pg.DB, Config *config.Config) {
	r.HandleFunc("/gov/proposals", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.GetProposals(DB, Config, w, r)
	})
	r.HandleFunc("/gov/proposal/{proposalId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.GetProposal(DB, Config, w, r)
	})
	r.HandleFunc("/gov/proposal/votes/{proposalId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.GetProposalVotes(DB, Config, w, r)
	})
	r.HandleFunc("/gov/proposal/deposits/{proposalId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.GetProposalDeposits(DB, w, r)
	})
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services.Test(RPCClient, DB, w, r)
	})
}
