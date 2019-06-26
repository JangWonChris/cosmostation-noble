package app

import (
	"fmt"
	"log"
	"net/http"

	gaiaApp "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/api/controllers"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/api/models"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/config"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/rpc/client"
)

// App wraps up the required variables that are needed in this app
type App struct {
	Codec     *codec.Codec
	Config    *config.Config
	DB        *pg.DB
	Router    *mux.Router
	RPCClient *client.HTTP
}

// Initializez the app with predefined configuration
func (a *App) Initialize(config *config.Config) {

	// Configuration
	a.Config = config

	// Connect to Tendermint RPC client through websocket
	a.RPCClient = client.NewHTTP(a.Config.Node.GaiadURL, "/websocket")

	// Connect to PostgreSQL
	a.DB = pg.Connect(&pg.Options{
		Addr:     a.Config.DB.Host,
		User:     a.Config.DB.User,
		Password: a.Config.DB.Password,
		Database: a.Config.DB.Table,
	})

	// Setup database schema
	_ = CreateSchema(a.DB)

	// Register Cosmos SDK codecs
	a.Codec = gaiaApp.MakeCodec()

	// Register routers
	a.setRouters()
}

// Create tables if it doesn't already exist
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*models.Account)(nil), (*models.Version)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// Sets the all required routers
func (a *App) setRouters() {
	a.Router = mux.NewRouter()
	a.Router = a.Router.PathPrefix("/v1").Subrouter()
	// 지금은 사용안함
	// a.Router.Use(auth.JwtAuthentication) // attach JWT auth middleware
	// controllers.AuthController(a.Router, a.RPCClient, a.DB)
	controllers.AccountController(a.Router, a.RPCClient, a.DB)
	controllers.AlarmController(a.Router, a.RPCClient, a.DB)
	controllers.VersionController(a.Router, a.RPCClient, a.DB)
}

// Run the app
func (a *App) Run(host string) {
	fmt.Print("Server is Starting on http://localhost", host, "\n")
	log.Fatal(http.ListenAndServe(host, a.Router))
}
