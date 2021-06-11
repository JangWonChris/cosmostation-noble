package custom

import (
	"github.com/cosmostation/cosmostation-cosmos/app"
	"github.com/gorilla/mux"
)

// RegisterHandlers registers all common query HTTP REST handlers on the provided mux router
func RegisterHandlers(a *app.App, r *mux.Router) {
	r.HandleFunc("/market/chart", GetCoinMarketChartData(a)).Methods("GET")
	r.HandleFunc("/market/{id}", GetSimpleCoinPrice(a)).Methods("GET")
}
