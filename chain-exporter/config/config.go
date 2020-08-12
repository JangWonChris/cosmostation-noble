package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config wraps all configs that are used in this project
type Config struct {
	Node       Node     `mapstructure:"node"`
	DB         Database `mapstructure:"database"`
	Alarm      Alarm    `mapstructure:"alarm"`
	KeybaseURL string   `mapstructure:"keybase_url"`
}

// Node wraps both endpoints for Tendermint RPC Node and REST API Server
type Node struct {
	RPCNode     string `mapstructure:"rpc_node"`
	LCDEndpoint string `mapstructure:"lcd_endpoint"`
}

// Database wraps PostgreSQL database config
type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Table    string `mapstructure:"table"`
}

// Alarm wraps push notification alarm config
type Alarm struct {
	PushServerEndpoint string `mapstructure:"push_server_endpoint"`
	Switch             bool   `mapstructure:"switch"`
}

// ParseConfig configures configuration
func ParseConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")                                              // for test cases
	viper.AddConfigPath("/home/ubuntu/cosmostation-cosmos/chain-exporter/") // for production

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	if viper.GetString("network_type") == "" {
		log.Fatal("define active param in your config file.")
	}

	var config Config
	sub := viper.Sub(viper.GetString("network_type"))
	sub.Unmarshal(&config)

	config.KeybaseURL = viper.GetString("keybase_endpoint")

	return &config
}
