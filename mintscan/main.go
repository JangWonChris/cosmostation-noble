package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/log"
	cfg "github.com/cosmostation/mintscan-backend-library/config"

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
	config := cfg.ParseConfig()

	// Create new client with node configruation.
	// Client is used for requesting any type of network data from RPC full node and REST Server.
	client := client.NewClient(&config.Client)

	// Create connection with PostgreSQL database and
	// Ping database to verify connection is success.
	db := db.Connect(&config.DB)
	err := db.Ping()
	if err != nil {
		zap.L().Error("failed to ping database", zap.Error(err))
		return
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	// account prefix를 가진 모든 REST API는
	// account 별 잔액을 조회
	// account 별 위임 상세 내역
	// account 별 위임 해제 상세 내역
	// account 별 tx 상세 내역
	r.HandleFunc("/auth/accounts/{accAddr}", handler.GetAccount).Methods("GET")                                               //kava에서만 사용중 (vesting account를 뽑기 위해)
	r.HandleFunc("/account/balances/{accAddr}", handler.GetBalance).Methods("GET")                                            // return all assets of given account
	r.HandleFunc("/account/validator/commission/{accAddr}", handler.GetValidatorCommission).Methods("GET")                    // 현재 사용중이나 /account/balances에 포함시킬 예정
	r.HandleFunc("/distribution/delegators/{delAddr}/withdraw_address", handler.GetDelegatorWithdrawalAddress).Methods("GET") //위임 내역을 반환할 때, 같이 포함시킨다

	r.HandleFunc("/account/delegations/{accAddr}", handler.GetDelegatorDelegations).Methods("GET")                    //
	r.HandleFunc("/account/unbonding_delegations/{delAddr}", handler.GetDelegatorUnbondingDelegations).Methods("GET") //moved to staking

	// 모바일 API
	r.HandleFunc("/account/txs/{accAddr}", handler.GetAccountTxs).Methods("GET")
	r.HandleFunc("/account/txs/{accAddr}/{valAddr}", handler.GetTxsBetweenDelegatorAndValidator).Methods("GET")
	r.HandleFunc("/account/transfer_txs/{accAddr}", handler.GetAccountTransferTxs).Methods("GET")

	r.HandleFunc("/blocks", handler.GetBlocks).Methods("GET")
	r.HandleFunc("/blocks/{proposer}", handler.GetBlocksByProposer).Methods("GET")

	//사용 안하는 중
	r.HandleFunc("/distribution/community_pool", handler.GetCommunityPool).Methods("GET")
	//end

	r.HandleFunc("/gov/proposals", handler.GetProposals).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposal_id}", handler.GetProposal).Methods("GET")
	r.HandleFunc("/gov/proposal/deposits/{proposal_id}", handler.GetDeposits).Methods("GET")
	r.HandleFunc("/gov/proposal/votes/{proposal_id}", handler.GetVotes).Methods("GET")

	r.HandleFunc("/minting/inflation", handler.GetMintingInflation).Methods("GET")
	r.HandleFunc("/market/chart", handler.GetCoinMarketChartData).Methods("GET")
	r.HandleFunc("/market/{id}", handler.GetSimpleCoinPrice).Methods("GET")
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
	r.HandleFunc("/staking/delegator/{delAddr}/redelegations", handler.GetRedelegations).Methods("GET")
	r.HandleFunc("/staking/delegators/{delAddr}/unbonding_delegations", handler.GetDelegatorUnbondingDelegations).Methods("GET")

	// 다음 버전에 업데이트 될 APIs
	r.HandleFunc("/account/total/balance/{accAddr}", handler.GetTotalBalance).Methods("GET")

	// These APIs will be deprecated in next update.
	r.HandleFunc("/staking/redelegations", handler.GetRedelegationsLegacy).Methods("GET")                             //staking/delegator/{delAddr}/redelegations 로 변경 됨
	r.HandleFunc("/account/balance/{accAddr}", handler.GetBalance).Methods("GET")                                     // /account/balances/{accAddr}
	r.HandleFunc("/account/commission/{accAddr}", handler.GetValidatorCommission).Methods("GET")                      // /account/validator/commission/{accAddr}
	r.HandleFunc("/account/unbonding-delegations/{delAddr}", handler.GetDelegatorUnbondingDelegations).Methods("GET") // /acount/unbonding_delegations/{accAddr}
	r.HandleFunc("/tx/{hash}", handler.GetLegacyTransactionFromDB).Methods("GET")                                     // /tx?hash={hash}

	// These APIs will need to be added in next update.
	// r.HandleFunc("/module/accounts", handler.GetModuleAccounts).Methods("GET")

	// Session will wrap both client and database and be used throughout all handlers.
	handler.SetSession(client, db)

	sm := &http.Server{
		Addr:         ":" + config.Web.Port,
		Handler:      handler.Middleware(r),
		ReadTimeout:  10 * time.Second, // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
	}

	go func() {
		for {
			if err := handler.SetStatus(); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Start the Mintscan API server.
	go func() {
		zap.S().Infof("Server is running on http://localhost:%s", config.Web.Port)
		zap.S().Infof("Version: %s | Commit: %s", Version, Commit)

		err := sm.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()

	TrapSignal(sm)
}

//TrapSignal trap the signal from os
func TrapSignal(sm *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,  // interrupt
		syscall.SIGKILL, // kill
		syscall.SIGTERM, // terminate
		syscall.SIGHUP)  // hangup(reload)

	for {
		sig := <-c // Block until a signal is received.
		switch sig {
		case syscall.SIGHUP:
			// cfg.ReloadConfig()
		default:
			terminate(sm, sig)
			break
		}
	}
}

func terminate(sm *http.Server, sig os.Signal) {
	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sm.Shutdown(ctx)

	zap.S().Infof("Gracefully shutting down the server: %s", sig)
}
