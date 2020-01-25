package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

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
	app := &App{
		codec:     ceCodec.Codec,
		config:    config,
		db:        db.Connect(config),
		router:    setRouter(),
		rpcClient: client.NewHTTP(config.Node.RPCNode, "/websocket"), // Tendermint RPC client
	}

	// Ping database to verify connection is succeeded
	err := app.db.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to ping database."))
	}

	controllers.AccountController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.BlockController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.DistributionController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.GovController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.MintingController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.TxController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.ValidatorController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.StatusController(app.codec, app.config, app.db, app.router, app.rpcClient)
	controllers.StatsController(app.codec, app.config, app.db, app.router, app.rpcClient)

	resty.SetTimeout(5 * time.Second) // sets timeout for request.

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
