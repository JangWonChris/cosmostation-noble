package controllers

import (
	"net/http"

	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/services"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
)

// Passes requests to its respective service
func TransactionController(codec *codec.Codec, config *config.Config, db *db.Database, r *mux.Router, rpcClient *client.HTTP) {
	r.HandleFunc("/txs", func(w http.ResponseWriter, r *http.Request) {
		services.GetTxs(codec, db, rpcClient, w, r)
	})
	r.HandleFunc("/tx/{hash}", func(w http.ResponseWriter, r *http.Request) {
		services.GetTx(codec, config, db, rpcClient, w, r)
	})
	r.HandleFunc("/tx/broadcast/{hash}", func(w http.ResponseWriter, r *http.Request) {
		services.BroadcastTx(codec, rpcClient, w, r)
	})

	/*
		[TODO]:
			1. /tx/{hash} 버그 수정
			2. 계정에서 발생한 트랜잭션 조회 API (getTxsByAddr)
			3. [리치리스트 방법]
	*/
}
