package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/chain-exporter/server"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/chain-exporter/subscribe"
	"github.com/cosmostation/cosmostation-cosmos/api/wallet/app/config"

	"github.com/tendermint/tendermint/libs/log"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func main() {
	// Start the syncing task
	sync()

}

func sync() {
	config := config.GetMainnetDevConfig()

	subscriber := subscribe.NewSubscriber(logger)
	ss := server.NewSubServer(logger, subscriber, config.Node.GaiadURL, "/websocket")

	// Stop upon receiving SIGTERM or CTRL-C
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range exit {
			logger.Error(fmt.Sprintf("captured %v, exiting...", sig))
			if ss.IsRunning() {
				ss.Stop()
			}
			os.Exit(1)
		}
	}()

	if err := ss.Start(); err != nil {
		logger.Error(fmt.Sprintf("Failed to start: %v", err))
		os.Exit(1)
	}
	logger.Info("Started http/ws client")

	// Run forever
	select {}
}
