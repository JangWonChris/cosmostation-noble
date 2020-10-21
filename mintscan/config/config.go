package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// CommonConfig wraps common configs that are used in this project.
type CommonConfig struct {
	NetworkType string `mapstructure:"network_type"`
	Maintenance bool   `mapstructure:"maintenance"`
}

// Config wraps all configs that are used in this project.
type Config struct {
	Node   NodeConfig   `mapstructure:"node"`
	DB     DBConfig     `mapstructure:"database"`
	Web    WebConfig    `mapstructure:"web"`
	Market MarketConfig `mapstructure:"market"`
}

// NodeConfig wraps both endpoints for Tendermint RPC Node and REST API Server.
type NodeConfig struct {
	RPCNode     string `mapstructure:"rpc_node"`
	LCDEndpoint string `mapstructure:"lcd_endpoint"`
}

// DBConfig wraps PostgreSQL database config.
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Table    string `mapstructure:"table"`
}

// WebConfig wraps port number of this project.
type WebConfig struct {
	Port string `mapstructure:"port"`
}

// MarketConfig wraps endpoints for CoinmarketCap and CoinGecko.
// In this project, we primarily use CoinGecko.
type MarketConfig struct {
	CoinGeckoEndpoint     string `mapstructure:"coingecko_endpoint"`
	CoinmarketCapEndpoint string `mapstructure:"coinmarketcap_endpoint"`
	CoinmarketCapCoinID   string `mapstructure:"coinmarketcap_coin_id"`
	CoinmarketCapAPIKey   string `mapstructure:"coinmarketcap_api_key"`
}

//Common is store common information read from config file
var Common CommonConfig

// ParseConfig attempts to read and parse config.yaml from the given path
// An error reading or parsing the config results in a panic.
func ParseConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")                                               // for test cases
	viper.AddConfigPath(os.Getenv("HOME") + "/cosmostation-cosmos/mintscan") // for production

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	if viper.GetString("network_type") == "" {
		zap.S().Fatal("define network_type param in your config file.")
	}

	var config Config

	Common.NetworkType = viper.GetString("network_type")
	Common.Maintenance = viper.GetBool("maintenance")

	sub := viper.Sub(Common.NetworkType)
	if err := sub.Unmarshal(&config); err != nil {
		zap.S().Fatal("error occurs while reading config file")
	}

	return &config
}

// ReloadConfig is reload Maintenance param when receive syscall.SIGHUP from os
func ReloadConfig() {
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	Common.Maintenance = viper.GetBool("maintenance")
	zap.S().Info("reload complete, Maintenance :", Common.Maintenance)
}
