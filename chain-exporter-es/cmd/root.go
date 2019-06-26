package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// 모든 앱 설정

var rootCmd = &cobra.Command{
	Use:"ElasticSearch Crawler",
	Short:"ElasticSearch Crawler Application",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute()  {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		//viper.AddConfigPath("/home/ubuntu/cosmos-proxy-api/api/proxy/")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to read config: %v\n", err)
		os.Exit(1)
	}
}
