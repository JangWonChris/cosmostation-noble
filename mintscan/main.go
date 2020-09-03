package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/log"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var (
	// Version is a project's version string.
	Version = "Development"

	// Commit is commit hash of this project.
	Commit = ""
)

func main() {
	// Create custom logger with a combination of using uber/zap and lumberjack.v2.
	l, _ := log.NewCustomLogger()
	zap.ReplaceGlobals(l)
	defer l.Sync()

	// Parse config from configuration file (config.yaml).
	config := config.ParseConfig()

	// Create new client with node configruation.
	// Client is used for requesting any type of network data from RPC full node and REST Server.
	client, err := client.NewClient(config.Node, config.Market)
	if err != nil {
		zap.L().Error("failed to create new client", zap.Error(err))
		return
	}

	// Create connection with PostgreSQL database and
	// Ping database to verify connection is success.
	db := db.Connect(config.DB)
	err = db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	r.HandleFunc("/auth/accounts/{accAddr}", handler.GetAccount).Methods("GET")
	r.HandleFunc("/account/balances/{accAddr}", handler.GetAccountBalance).Methods("GET")
	r.HandleFunc("/account/delegations/{accAddr}", handler.GetDelegatorDelegations).Methods("GET")
	r.HandleFunc("/account/delegations/rewards/{accAddr}", handler.GetDelegationsRewards).Methods("GET")
	r.HandleFunc("/account/validator/commission/{accAddr}", handler.GetValidatorCommission).Methods("GET")
	r.HandleFunc("/account/unbonding_delegations/{accAddr}", handler.GetDelegatorUnbondingDelegations).Methods("GET")
	r.HandleFunc("/account/txs/{accAddr}", handler.GetAccountTxs).Methods("GET")                                // Mobile
	r.HandleFunc("/account/txs/{accAddr}/{valAddr}", handler.GetTxsBetweenDelegatorAndValidator).Methods("GET") // Mobile
	r.HandleFunc("/account/transfer_txs/{accAddr}", handler.GetAccountTransferTxs).Methods("GET")               // Mobile
	r.HandleFunc("/blocks", handler.GetBlocks).Methods("GET")
	r.HandleFunc("/blocks/{proposer}", handler.GetBlocksByProposer).Methods("GET")
	r.HandleFunc("/distribution/delegators/{delAddr}/rewards/{valAddr}", handler.GetRewardsBetweenDelegatorAndValidator).Methods("GET")
	r.HandleFunc("/distribution/delegators/{delAddr}/withdraw_address", handler.GetDelegatorWithdrawalAddress).Methods("GET")
	r.HandleFunc("/distribution/community_pool", handler.GetCommunityPool).Methods("GET")
	r.HandleFunc("/gov/proposals", handler.GetProposals).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposal_id}", handler.GetProposal).Methods("GET")
	r.HandleFunc("/gov/proposal/deposits/{proposal_id}", handler.GetDeposits).Methods("GET")
	r.HandleFunc("/gov/proposal/votes/{proposal_id}", handler.GetVotes).Methods("GET")
	r.HandleFunc("/minting/inflation", handler.GetMintingInflation).Methods("GET")
	r.HandleFunc("/market/{id}", handler.GetCoinPrice).Methods("GET")
	r.HandleFunc("/stats/market", handler.GetMarketStats).Methods("GET")
	r.HandleFunc("/stats/network", handler.GetNetworkStats).Methods("GET")
	r.HandleFunc("/status", handler.GetStatus).Methods("GET")
	r.HandleFunc("/txs", handler.GetTransactions).Methods("GET")
	r.HandleFunc("/txs", handler.GetTransactionsList).Methods("POST")
	r.HandleFunc("/tx", handler.GetTransaction).Methods("GET")
	r.HandleFunc("/tx/broadcast/{signed_tx}", handler.BroadcastTx).Methods("GET")
	r.HandleFunc("/staking/validators", handler.GetValidators).Methods("GET")
	r.HandleFunc("/staking/validator/{address}", handler.GetValidator).Methods("GET")
	r.HandleFunc("/staking/validator/uptime/{address}", handler.GetValidatorUptime).Methods("GET")
	r.HandleFunc("/staking/validator/uptime/range/{address}", handler.GetValidatorUptimeRange).Methods("GET")
	r.HandleFunc("/staking/validator/delegations/{address}", handler.GetValidatorDelegations).Methods("GET")
	r.HandleFunc("/staking/validator/events/{address}", handler.GetValidatorPowerHistoryEvents).Methods("GET")
	r.HandleFunc("/staking/validator/events/{address}/count", handler.GetValidatorEventsTotalCount).Methods("GET")
	r.HandleFunc("/staking/redelegations", handler.GetRedelegations).Methods("GET")

	// These APIs will be deprecated in next update.
	r.HandleFunc("/account/balance/{accAddr}", handler.GetAccountBalance).Methods("GET")                              // /account/balances/{accAddr}
	r.HandleFunc("/account/commission/{accAddr}", handler.GetValidatorCommission).Methods("GET")                      // /account/validator/commission/{accAddr}
	r.HandleFunc("/account/unbonding-delegations/{accAddr}", handler.GetDelegatorUnbondingDelegations).Methods("GET") // /acount/unbonding_delegations/{accAddr}
	r.HandleFunc("/tx/{hash}", handler.GetLegacyTransactionFromDB).Methods("GET")                                     // /tx?hash=
	r.HandleFunc("/staking/validator/misses/detail/{address}", handler.GetValidatorUptime).Methods("GET")             // /staking/validator/updatime/{address}
	r.HandleFunc("/staking/validator/misses/{address}", handler.GetValidatorUptimeRange).Methods("GET")               // /staking/validator/uptime/range/{address}

	// These APIs will need to be added in next update.
	// r.HandleFunc("/module/accounts", handler.GetModuleAccounts).Methods("GET")

	sm := &http.Server{
		Addr:         ":" + config.Web.Port,
		Handler:      handler.Middleware(r, client, db),
		ReadTimeout:  10 * time.Second, // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
	}

	// Start the Mintscan API server.
	go func() {
		zap.S().Infof("Server is running on http://localhost:%s", config.Web.Port)
		zap.S().Infof("Network Type: %s | Version: %s | Commit: %s", config.Node.NetworkType, Version, Commit)

		err := sm.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()

	TrapSignal(sm)
}

// TrapSignal traps sigterm or interupt and gracefully shutdown the server.
func TrapSignal(sm *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c

	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sm.Shutdown(ctx)

	zap.S().Infof("Gracefully shutting down the server: %s", sig)
}
