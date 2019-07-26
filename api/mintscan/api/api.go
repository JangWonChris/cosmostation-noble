package api

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	gaiaApp "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/config"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/controllers"
	"github.com/cosmostation/cosmostation-cosmos/api/mintscan/api/databases"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"

	resty "gopkg.in/resty.v1"
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
	// Configuration
	a.config = config

	// Connect to Tendermint RPC client through websocket
	a.rpcClient = client.NewHTTP(a.config.Node.GaiadURL, "/websocket")

	// Connect to PostgreSQL
	a.db = databases.ConnectDatabase(config)

	// Register Cosmos SDK codecs
	a.codec = gaiaApp.MakeCodec()

	// Register routers
	a.setRouters()

	// SetTimeout method sets timeout for request.
	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // Local 환경에서 테스트를 위해
}

// Sets the all required routers
func (a *App) setRouters() {
	a.router = mux.NewRouter()
	a.router = a.router.PathPrefix("/v1").Subrouter()
	controllers.AccountController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.BlockController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.DistributionController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.GovernanceController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.MintingController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.TransactionController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.ValidatorController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.StatusController(a.codec, a.config, a.db, a.router, a.rpcClient)
	controllers.StatsController(a.codec, a.config, a.db, a.router, a.rpcClient)
}

// Run the app
func (a *App) Run(host string) {
	fmt.Print("Server is Starting on http://localhost", host, "\n")
	log.Fatal(http.ListenAndServe(host, a.router))
}
