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
	PrePareMsgExp()
	r.HandleFunc("/account/txs/{accAddr}", GetAccountTxsHistory).Methods("GET")
	r.HandleFunc("/account/transfer_txs/{accAddr}", GetAccountTransferTxsHistory).Methods("GET")
	r.HandleFunc("/account/txs/{accAddr}/{valAddr}", GetTxsHistoryBetweenDelegatorAndValidator).Methods("GET")
}
