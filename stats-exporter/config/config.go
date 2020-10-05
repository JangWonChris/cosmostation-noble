package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config wraps all configs that are used in this project.
type Config struct {
	Node   Node     `mapstructure:"node"`
	DB     Database `mapstructure:"database"`
	Market Market   `mapstructure:"market"`
}

// Node wraps both endpoints for Tendermint RPC Node and REST API Server.
type Node struct {
	RPCNode     string `mapstructure:"rpc_node"`
	LCDEndpoint string `mapstructure:"lcd_endpoint"`
}

// Database wraps PostgreSQL database config.
type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Table    string `mapstructure:"table"`
}

// Market wraps endpoints for CoinmarketCap and CoinGecko.
// In this project, we primarily use CoinGecko.
type Market struct {
	CoinGeckoEndpoint string `mapstructure:"coingecko_endpoint"`
}

// ParseConfig attempts to read and parse config.yaml from the given path.
// An error reading or parsing the config results in a panic.
func ParseConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")                                                      // for test cases
	viper.AddConfigPath(os.Getenv("HOME") + "/cosmostation-cosmos/stats-exporter/") // for production

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	if viper.GetString("network_type") == "" {
		log.Fatal("define network_type param in your config file.")
	}

	var config Config
	sub := viper.Sub(viper.GetString("network_type"))
	sub.Unmarshal(&config)

	return &config
}
