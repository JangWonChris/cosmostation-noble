package custom

import (
	"github.com/cosmostation/cosmostation-cosmos/handler"

	"github.com/gorilla/mux"
)

// s is shorten for handler Session
var s *handler.Session

// RegisterHandlers registers all common query HTTP REST handlers on the provided mux router
func RegisterHandlers(session *handler.Session, r *mux.Router) {
	s = session

	// r.HandleFunc("/account/tokens/{accAddr}", GetTokensBalances).Methods("GET")
	r.HandleFunc("/market/chart", GetCoinMarketChartData).Methods("GET")
	r.HandleFunc("/market/{id}", GetSimpleCoinPrice).Methods("GET")
}
