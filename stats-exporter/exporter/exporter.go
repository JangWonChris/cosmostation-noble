package exporter

import (
	"context"
	"crypto/tls"
	"os"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/databases"

	gaiaApp "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/go-pg/pg"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"

	// "github.com/robfig/cron"

	resty "gopkg.in/resty.v1"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

// Monitor wraps Tendermint RPC client and PostgreSQL database
type StatsExporterService struct {
	cmn.BaseService
	codec     *codec.Codec
	config    *config.Config
	db        *pg.DB
	wsCtx     context.Context
	rpcClient *client.HTTP
}

// Initializes all the required configs
func NewStatsExporterService(config *config.Config) *StatsExporterService {
	ses := &StatsExporterService{
		codec:     gaiaApp.MakeCodec(), // Register Cosmos SDK codecs
		config:    config,
		db:        databases.ConnectDatabase(config), // Connect to PostgreSQL
		wsCtx:     context.Background(),
		rpcClient: client.NewHTTP(config.Node.GaiadURL, "/websocket"), // Connect to Tendermint RPC client
	}
	// Setup database schema
	databases.CreateSchema(ses.db)

	// SetTimeout method sets timeout for request.
	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // Test locally

	return ses
}

// Override method for BaseService, which starts a service
func (ses *ChainExporterService) OnStart() error {
	// Cron jobs every 1 hour
	// c := cron.New()
	// c.AddFunc("0 */60 * * *", func() { ses.SaveValidatorStats() })
	// c.AddFunc("0 */60 * * *", func() { ses.SaveNetworkStats() })
	// c.AddFunc("0 */60 * * *", func() { ses.SaveCoinMarketCapMarketStats() })
	// c.AddFunc("0 */60 * * *", func() { ses.SaveCoinGeckoMarketStats() })
	// c.AddFunc("0 */15 * * *", func() { ses.SaveValidatorKeyBase() })
	// go c.Start()

	// // Allow graceful closing of the governance loop
	// signalCh := make(chan os.Signal, 1)
	// signal.Notify(signalCh, os.Interrupt)
	// <-signalCh

	ses.SaveValidatorStats()
	ses.SaveNetworkStats()
	ses.SaveCoinMarketCapMarketStats()
	ses.SaveCoinGeckoMarketStats()
	ses.SaveValidatorKeyBase()
}
