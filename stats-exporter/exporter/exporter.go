package exporter

import (
	"context"
	"crypto/tls"
	"os"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/databases"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/utils"

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
		codec:     utils.MakeCodec(), // register Cosmos SDK codecs
		config:    config,
		db:        databases.ConnectDatabase(config), // connect to PostgreSQL
		wsCtx:     context.Background(),
		rpcClient: client.NewHTTP(config.Node.GaiadURL, "/websocket"), // connect to Tendermint RPC client
	}
	// create database schema
	databases.CreateSchema(ses.db)

	// SetTimeout method sets timeout for request.
	resty.SetTimeout(5 * time.Second)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // for local test

	return ses
}

// Override method for BaseService, which starts a service
func (ses *StatsExporterService) OnStart() {
	// Cron jobs
	// c := cron.New()

	// Every hour
	// 0 * * * * = every minute
	// 0 */60 * * * = every hour
	// 0 0 * * * * = every hour
	// c.AddFunc("0 0 * * * *", func() { ses.SaveValidatorsStats1H() })
	// c.AddFunc("0 0 * * * *", func() { ses.SaveNetworkStats1H() })
	// c.AddFunc("0 0 * * * *", func() { ses.SaveCoinGeckoMarketStats1H() })
	// c.AddFunc("0 0 * * * *", func() { ses.SaveCoinMarketCapMarketStats1H() })

	// // Every day at 2:00 AM (UTC zone) which equals 11:00 AM in Seoul
	// c.AddFunc("0 0 2 * * *", func() { ses.SaveValidatorsStats24H() })
	// c.AddFunc("0 0 2 * * *", func() { ses.SaveNetworkStats24H() })
	// c.AddFunc("0 0 2 * * *", func() { ses.SaveCoinGeckoMarketStats24H() })
	// c.AddFunc("0 0 2 * * *", func() { ses.SaveCoinMarketCapMarketStats24H() })
	// go c.Start()

	// // Allow graceful closing of the governance loop
	// signalCh := make(chan os.Signal, 1)
	// signal.Notify(signalCh, os.Interrupt)
	// <-signalCh

	ses.SaveValidatorsStats1H()
	ses.SaveValidatorsStats24H()

	ses.SaveNetworkStats1H()
	ses.SaveNetworkStats24H()

	ses.SaveCoinGeckoMarketStats1H()
	ses.SaveCoinGeckoMarketStats24H()

	ses.SaveCoinMarketCapMarketStats1H()
	ses.SaveCoinMarketCapMarketStats24H()

}
