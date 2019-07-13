package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("/home/ubuntu/cosmostation-cosmos/stats-exporter") // call multiple times to add many search paths

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config := &Config{}
	config.KeybaseURL = viper.GetString("keybase_url")

	nodeConfig := &NodeConfig{}
	dbConfig := &DBConfig{}

	// Production or Development
	switch viper.GetString("active") {
	case "prod":
		nodeConfig.GaiadURL = viper.GetString("prod.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("prod.node.lcd_url")
		dbConfig.Host = viper.GetString("prod.database.host")
		dbConfig.User = viper.GetString("prod.database.user")
		dbConfig.Password = viper.GetString("prod.database.password")
		dbConfig.Table = viper.GetString("prod.database.table")

	case "dev":
		nodeConfig.GaiadURL = viper.GetString("dev.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("dev.node.lcd_url")
		dbConfig.Host = viper.GetString("dev.database.host")
		dbConfig.User = viper.GetString("dev.database.user")
		dbConfig.Password = viper.GetString("dev.database.password")
		dbConfig.Table = viper.GetString("dev.database.table")

	case "testnet":
		nodeConfig.GaiadURL = viper.GetString("testnet.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("testnet.node.lcd_url")
		dbConfig.Host = viper.GetString("testnet.database.host")
		dbConfig.User = viper.GetString("testnet.database.user")
		dbConfig.Password = viper.GetString("testnet.database.password")
		dbConfig.Table = viper.GetString("testnet.database.table")

	default:
		fmt.Println("Define active params in config.yaml")
	}

	config.Node = nodeConfig
	config.DB = dbConfig

	return config
}
