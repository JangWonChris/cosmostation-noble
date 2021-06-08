package common

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	"github.com/gorilla/mux"
)

// s is shorten for handler Session
var s *handler.Session

// RegisterHandlers registers all common query HTTP REST handlers on the provided mux router
func RegisterHandlers(session *handler.Session, r *mux.Router) {
	s = session

	r.HandleFunc("/auth/accounts/{accAddr}", GetAuthAccount).Methods("GET")
	r.HandleFunc("/account/module_accounts", GetModuleAccounts).Methods("GET")
	// r.HandleFunc("/account/balances/{accAddr}", GetBalance).Methods("GET")
	r.HandleFunc("/account/balances/{accAddr}", GetAllBalances).Methods("GET")
	r.HandleFunc("/account/delegations/{accAddr}", GetDelegatorDelegations).Methods("GET")
	r.HandleFunc("/account/undelegations/{accAddr}", GetDelegatorUnbondingDelegations).Methods("GET")

	r.HandleFunc("/blocks", GetBlocks).Methods("GET") // ?from=&limit=
	r.HandleFunc("/block/id/{id}", GetBlocksByID).Methods("GET")
	r.HandleFunc("/block/hash/{hash}", GetBlocksByHash).Methods("GET")
	r.HandleFunc("/block/{chainid}/{height}", GetBlockByChainIDHeight).Methods("GET")

	r.HandleFunc("/gov/proposals", GetProposals).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposal_id}", GetProposal).Methods("GET")
	r.HandleFunc("/gov/proposal/deposits/{proposal_id}", GetDeposits).Methods("GET")
	r.HandleFunc("/gov/proposal/votes/{proposal_id}", GetVotes).Methods("GET")

	r.HandleFunc("/minting/inflation", GetMintingInflation).Methods("GET")

	r.HandleFunc("/txs", GetTransactions).Methods("GET") // &from=&limit=
	r.HandleFunc("/tx/hash/{hash}", GetTransactionByHash).Methods("GET")
	r.HandleFunc("/tx/id/{id}", GetTransactionByID).Methods("GET")
	r.HandleFunc("/txs", GetTransactionsList).Methods("POST")

	// r.HandleFunc("/txs", GetTransactions).Methods("GET") //단종
	r.HandleFunc("/tx/broadcast/{signed_tx}", BroadcastTx).Methods("GET")

	r.HandleFunc("/staking/validators", GetValidators).Methods("GET")
	r.HandleFunc("/staking/validator/{address}", GetValidator).Methods("GET")
	r.HandleFunc("/staking/validator/blocks/{proposer}", GetValidatorProposedBlocks).Methods("GET")
	r.HandleFunc("/staking/validator/uptime/{address}", GetValidatorUptime).Methods("GET")
	r.HandleFunc("/staking/validator/uptime/range/{address}", GetValidatorUptimeRange).Methods("GET")
	r.HandleFunc("/staking/validator/delegations/{address}", GetValidatorDelegations).Methods("GET")
	r.HandleFunc("/staking/validator/events/{address}", GetValidatorPowerHistoryEvents).Methods("GET")
	r.HandleFunc("/staking/redelegations", GetRedelegationsLegacy).Methods("GET")
	// r.HandleFunc("/staking/redelegations", GetRedelegations).Methods("GET")
	r.HandleFunc("/status", GetStatus).Methods("GET")
	r.HandleFunc("/stats/market", GetMarketStats).Methods("GET")
	// r.HandleFunc("/stats/network", GetNetworkStats).Methods("GET") // 제거 2021.06.02

	// 부연님 커미션 조회할때 사용하는 api
	r.HandleFunc("/account/validator/commission/{accAddr}", GetValidatorCommission).Methods("GET") // (포함) /account/balances/{accAddr}
	r.HandleFunc("/distribution/community_pool", GetCommunityPool).Methods("GET")                  // (포함) /status
	/*
		변경 = 요청 Path 변경
		포함 = 다른 API에 포함
		확인 = 확실히 사용하고 있지 않은 API인지 확인 필요
		논의 = 논의가 필요한 API
		나중 = 나중에 필요할 수도 있을 만한 API
	*/

	r.HandleFunc("/account/total/balance/{accAddr}", GetTotalBalance).Methods("GET")                                  // (변경) /account/balances/{accAddr}
	r.HandleFunc("/blocks/{proposer}", GetBlocksByProposer).Methods("GET")                                            // (변경) /staking/validator/blocks/{proposer}
	r.HandleFunc("/account/unbonding_delegations/{accAddr}", GetDelegatorUndelegations).Methods("GET")                // (변경) /account/undelegations/{accAddr}
	r.HandleFunc("/distribution/delegators/{delAddr}/withdraw_address", GetDelegatorWithdrawalAddress).Methods("GET") // (포함) /account/balances/{accAddr}

	// r.HandleFunc("/account/delegations/rewards/{accAddr}", GetDelegationsRewards).Methods("GET")                                // (확인)
	// r.HandleFunc("/staking/validator/events/{address}/count", GetValidatorEventsTotalCount).Methods("GET")                      // (확인)
	// r.HandleFunc("/distribution/delegators/{delAddr}/rewards/{valAddr}", GetRewardsBetweenDelegatorAndValidator).Methods("GET") // (확인)
	r.HandleFunc("/staking/validator/uptime/range/{address}", GetValidatorUptimeRange).Methods("GET") // (나중)
}
