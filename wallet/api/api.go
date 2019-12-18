package api

import (
	"fmt"
	"log"
	"net/http"

	gaiaApp "github.com/cosmos/gaia/app"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/config"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/controllers"
	"github.com/cosmostation/cosmostation-cosmos/wallet/api/databases"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"

	"github.com/tendermint/tendermint/rpc/client"
)

// App wraps up the required variables that are needed in this app
type App struct {
	codec     *codec.Codec
	config    *config.Config
	db        *pg.DB
	router    *mux.Router
	rpcClient *client.HTTP
}

// NewApp initializes the app with predefined configuration
func (a *App) NewApp(config *config.Config) {
	a.config = config

	// connect to Tendermint RPC client through websocket
	a.rpcClient = client.NewHTTP(a.config.Node.GaiadURL, "/websocket")

	// connect to PostgreSQL and create database schema
	a.db = databases.ConnectDatabase(config)
	databases.CreateSchema(a.db)

	// register Cosmos SDK codecs
	a.codec = gaiaApp.MakeCodec()

	// register routers
	a.setRouters()
}

// setRouters sets the all routers
func (a *App) setRouters() {
	a.router = mux.NewRouter()
	a.router = a.router.PathPrefix("/v1").Subrouter()
	// 지금은 사용안함
	// a.router.Use(auth.JwtAuthentication) // attach JWT auth middleware
	// controllers.AuthController(a.router, a.rpcClient, a.db)
	controllers.AccountController(a.router, a.rpcClient, a.db)
	controllers.AlarmController(a.router, a.rpcClient, a.db)
	controllers.VersionController(a.router, a.rpcClient, a.db)
}

// Run the app
func (a *App) Run(host string) {
	fmt.Print("Server is Starting on http://localhost", host, "\n")
	log.Fatal(http.ListenAndServe(host, a.router))
}
