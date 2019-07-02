package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"github.com/cosmostation-cosmos/chain-exporter-es/app"

)


var (
	// prod, dev
	env string

	// cosmos, kava, iris
	network string
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start ElasticSearch Crawler deamon",
	RunE: serverCmdHandler,
}

func init()  {
	serverCmd.Flags().StringVar(&network, "network", "cosmos", "set network (cosmos or kava) - required")
	serverCmd.Flags().StringVar(&env, "env", "dev", "set build environment (dev or prod) - required")

	//Flags are optional by default. If instead you wish your command to report an error when a flag has not been set, mark it as required:
	serverCmd.MarkFlagRequired("env")
	serverCmd.MarkFlagRequired("network")
	rootCmd.AddCommand(serverCmd)
}

func serverCmdHandler(cmd *cobra.Command, args []string) error {

	a, err := app.NewApp(network, env)
	if err != nil {
		logrus.Error(err)
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
}