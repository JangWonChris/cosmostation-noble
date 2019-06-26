package cmd

import (
	"fmt"
	"github.com/cosmostation-cosmos/chain-exporter-es/app"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

// prod/dev
var env string

func init()  {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use: "server [env]",
	Short: "Start ElasticSearch Crawler deamon",
	Args:cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectenv := args[0]
		a, err := app.New(projectenv)
		if err != nil {
			return err
		}

		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		go func() {
			for sig := range exit {
				a.Logger.Error(fmt.Sprintf("captured %v, exiting...", sig))
				if a.IsRunning() {
					a.Stop()
				}
				os.Exit(1)
			}
		}()

		if err := a.Start(); err != nil {
			a.Logger.Error(fmt.Sprintf("Failed to start: %v", err))
			os.Exit(1)
		}
		a.Logger.Info("Started http/ws client")

		// Run forever
		select {}
		return nil
	},
}

