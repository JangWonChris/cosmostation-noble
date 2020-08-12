package db

import (
	"crypto/tls"
	"os"
	"time"

	"github.com/cosmostation/cosmostation-cosmos/chain-exporter/config"
	"github.com/go-pg/pg"
)

func pgOptions() *pg.Options {
	config := config.ParseConfig()

	return &pg.Options{
		TLSConfig: getTLSConfig(),

		Addr:     config.DB.Host + ":" + config.DB.Port,
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Table,

		MaxRetries:      1,
		MinRetryBackoff: -1,

		DialTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		PoolSize:           10,
		MaxConnAge:         10 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

func getTLSConfig() *tls.Config {
	pgSSLMode := os.Getenv("PGSSLMODE")
	if pgSSLMode == "disable" {
		return nil
	}
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}
