package mobile

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/handler"
	"github.com/gorilla/mux"
)

// s is shorten for handler Session
var s *handler.Session

// RegisterHandlers registers all common query HTTP REST handlers on the provided mux router
func RegisterHandlers(session *handler.Session, r *mux.Router) {
	s = session

	r.HandleFunc("/account/txs/{accAddr}", GetAccountTxs).Methods("GET")
	r.HandleFunc("/account/transfer_txs/{accAddr}", GetAccountTransferTxs).Methods("GET")
	r.HandleFunc("/account/txs/{accAddr}/{valAddr}", GetTxsBetweenDelegatorAndValidator).Methods("GET")
}
