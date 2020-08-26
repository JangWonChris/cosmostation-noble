package exporter

import (
	"os"
	"os/signal"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/client"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/config"
	"github.com/cosmostation/cosmostation-cosmos/stats-exporter/db"

	"go.uber.org/zap"

	cron "github.com/robfig/cron/v3"
)

var (
	// Version is a project's version string.
	Version = "Development"

	// Commit is commit hash of this project.
	Commit = ""
)

// Exporter wraps all required parameters to create exporter jobs
type Exporter struct {
	client *client.Client
	db     *db.Database
}

// NewExporter creates new exporter
func NewExporter() *Exporter {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)

	// Parse config from configuration file (config.yaml).
	config := config.ParseConfig()

	// Create new client with node configruation.
	// Client is used for requesting any type of network data from RPC full node and REST Server.
	client, err := client.NewClient(config.Node, config.Market)
	if err != nil {
		zap.S().Errorf("failed to create new client", err)
		return &Exporter{}
	}

	// Create connection with PostgreSQL database and
	// Ping database to verify connection is success.
	db := db.Connect(&config.DB)
	err = db.Ping()
	if err != nil {
		zap.S().Errorf("failed to ping database: %s", err)
		return &Exporter{}
	}

	// Create database tables if not exist already
	db.CreateTables()

	return &Exporter{client, db}
}

// Start starts to create cron jobs and start them
// Cron spec format can be found here https://godoc.org/gopkg.in/robfig/cron.v3#hdr-Intervals
func (ex *Exporter) Start() error {
	zap.S().Info("Starting Stat Exporter...")
	zap.S().Infof("Version: %s Commit: %s", Version, Commit)

	ex.SaveValidatorsStats1H()

	c := cron.New(
		cron.WithLocation(time.UTC),
	)

	// Run every 5 minutes: @every 5m
	c.AddFunc("@every 5m", func() {
		ex.SaveStatsMarket5M()
		zap.S().Info("successfully saved data @every 5m ")
	})

	// Run once an hour: @hourly or @every 1h
	c.AddFunc("@hourly", func() { // same as 0 * * * * *
		ex.SaveStatsMarket1H()
		ex.SaveNetworkStats1H()
		ex.SaveValidatorsStats1H()
		zap.S().Info("successfully saved data @hourly ")
	})

	// Run once a day: @daily or @midnight
	c.AddFunc("@midnight", func() { // same as 0 0 * * *
		ex.SaveStatsMarket1D()
		ex.SaveNetworkStats1D()
		ex.SaveValidatorsStats1D()
		zap.S().Info("successfully saved data @midnight ")
	})

	c.Start()

	// Allow graceful closing of cron jobs
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	return nil
}
