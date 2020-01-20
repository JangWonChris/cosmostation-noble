package api

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	ceCodec "github.com/cosmostation/cosmostation-cosmos/mintscan/api/codec"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/controllers"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/api/db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"

	resty "gopkg.in/resty.v1"
)

// App wraps up the required variables that are needed in this app
type App struct {
	codec     *codec.Codec
	config    *config.Config
	db        *db.Database
	router    *mux.Router
	rpcClient *client.HTTP
}

// NewApp initializes the app with predefined configuration
func NewApp(config *config.Config) *App {
	// sets timeout for request.
	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // local test

	app := &App{
		ceCodec.Codec,
		config,
		db.Connect(config),
		setRouter(),
		client.NewHTTP(config.Node.GaiadURL, "/websocket"), // Tendermint RPC client
	}

	controllers.AccountController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.BlockController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.DistributionController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.GovernanceController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.MintingController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.TransactionController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.ValidatorController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.StatusController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.StatsController(app.codec, app.config, app.db, app.router, app.rpcClient)

	return app
}

// setRouter sets the all required routers
func setRouter() *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()

	return r
}

// Run the app
func (a *App) Run(host string) {
	fmt.Print("Server is Starting on http://localhost", host, "\n")
	log.Fatal(http.ListenAndServe(host, a.router))
}
