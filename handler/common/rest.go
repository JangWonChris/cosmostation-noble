package common

import (
	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/gorilla/mux"
)

// RegisterHandlers registers all common query HTTP REST handlers on the provided mux router
func RegisterHandlers(a *app.App, r *mux.Router) {
	PrePareMsgExp(a)
	r.HandleFunc("/auth/accounts/{accAddr}", GetAuthAccount(a)).Methods("GET")
	r.HandleFunc("/account/module_accounts", GetModuleAccounts(a)).Methods("GET")
	// r.HandleFunc("/account/balances/{accAddr}", GetBalance).Methods("GET")
	r.HandleFunc("/account/balances/{accAddr}", GetAllBalances(a)).Methods("GET")
	r.HandleFunc("/account/delegations/{accAddr}", GetDelegatorDelegations(a)).Methods("GET")
	r.HandleFunc("/account/undelegations/{accAddr}", GetDelegatorUnbondingDelegations(a)).Methods("GET")

	r.HandleFunc("/blocks", GetBlocks(a)).Methods("GET") // ?from=&limit=
	r.HandleFunc("/block/id/{id}", GetBlocksByID(a)).Methods("GET")
	r.HandleFunc("/block/hash/{hash}", GetBlocksByHash(a)).Methods("GET")
	r.HandleFunc("/block/{chainid}/{height}", GetBlockByChainIDHeight(a)).Methods("GET")

	r.HandleFunc("/gov/proposals", GetProposals(a)).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposal_id}", GetProposal(a)).Methods("GET")
	r.HandleFunc("/gov/proposal/deposits/{proposal_id}", GetDeposits(a)).Methods("GET")
	r.HandleFunc("/gov/proposal/votes/{proposal_id}", GetVotes(a)).Methods("GET")

	r.HandleFunc("/minting/inflation", GetMintingInflation(a)).Methods("GET")

	r.HandleFunc("/txs", GetTransactions(a)).Methods("GET") // &from=&limit=
	r.HandleFunc("/tx/hash/{hash}", GetTransactionByHash(a)).Methods("GET")
	r.HandleFunc("/tx/id/{id}", GetTransactionByID(a)).Methods("GET")
	r.HandleFunc("/txs", GetTransactionsList(a)).Methods("POST")

	// r.HandleFunc("/txs", GetTransactions).Methods("GET") //단종
	r.HandleFunc("/tx/broadcast/{signed_tx}", BroadcastTx(a)).Methods("GET")

	r.HandleFunc("/staking/validators", GetValidators(a)).Methods("GET")
	r.HandleFunc("/staking/validator/{address}", GetValidator(a)).Methods("GET")
	r.HandleFunc("/staking/validator/blocks/{proposer}", GetValidatorProposedBlocks(a)).Methods("GET")
	r.HandleFunc("/staking/validator/uptime/{address}", GetValidatorUptime(a)).Methods("GET")
	r.HandleFunc("/staking/validator/uptime/range/{address}", GetValidatorUptimeRange(a)).Methods("GET")
	r.HandleFunc("/staking/validator/delegations/{address}", GetValidatorDelegations(a)).Methods("GET")
	r.HandleFunc("/staking/validator/events/{address}", GetValidatorPowerHistoryEvents(a)).Methods("GET")
	r.HandleFunc("/staking/redelegations", GetRedelegationsLegacy(a)).Methods("GET")
	// r.HandleFunc("/staking/redelegations", GetRedelegations).Methods("GET")
	r.HandleFunc("/status", GetStatus).Methods("GET")
	r.HandleFunc("/stats/market", GetMarketStats(a)).Methods("GET")
	// r.HandleFunc("/stats/network", GetNetworkStats).Methods("GET") // 제거 2021.06.02

	// 부연님 커미션 조회할때 사용하는 api
	r.HandleFunc("/account/validator/commission/{accAddr}", GetValidatorCommission(a)).Methods("GET") // (포함) /account/balances/{accAddr}
	r.HandleFunc("/distribution/community_pool", GetCommunityPool(a)).Methods("GET")                  // (포함) /status
	/*
		변경 = 요청 Path 변경
		포함 = 다른 API에 포함
		확인 = 확실히 사용하고 있지 않은 API인지 확인 필요
		논의 = 논의가 필요한 API
		나중 = 나중에 필요할 수도 있을 만한 API
	*/

	r.HandleFunc("/account/total/balance/{accAddr}", GetTotalBalance(a)).Methods("GET")                                  // (변경) /account/balances/{accAddr}
	r.HandleFunc("/blocks/{proposer}", GetBlocksByProposer(a)).Methods("GET")                                            // (변경) /staking/validator/blocks/{proposer}
	r.HandleFunc("/account/unbonding_delegations/{accAddr}", GetDelegatorUndelegations(a)).Methods("GET")                // (변경) /account/undelegations/{accAddr}
	r.HandleFunc("/distribution/delegators/{delAddr}/withdraw_address", GetDelegatorWithdrawalAddress(a)).Methods("GET") // (포함) /account/balances/{accAddr}

	// r.HandleFunc("/account/delegations/rewards/{accAddr}", GetDelegationsRewards).Methods("GET")                                // (확인)
	// r.HandleFunc("/staking/validator/events/{address}/count", GetValidatorEventsTotalCount).Methods("GET")                      // (확인)
	// r.HandleFunc("/distribution/delegators/{delAddr}/rewards/{valAddr}", GetRewardsBetweenDelegatorAndValidator).Methods("GET") // (확인)
	r.HandleFunc("/staking/validator/uptime/range/{address}", GetValidatorUptimeRange(a)).Methods("GET") // (나중)

	r.HandleFunc("/account/txs/{accAddr}", GetAccountTxsHistory(a)).Methods("GET")
	r.HandleFunc("/account/transfer_txs/{accAddr}", GetAccountTransferTxsHistory(a)).Methods("GET")
	r.HandleFunc("/account/txs/{accAddr}/{valAddr}", GetTxsHistoryBetweenDelegatorAndValidator(a)).Methods("GET")
}
