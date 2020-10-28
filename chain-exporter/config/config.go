package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config wraps all configs that are used in this project.
type Config struct {
	Node       Node     `mapstructure:"node"`
	DB         Database `mapstructure:"database"`
	Alarm      Alarm    `mapstructure:"alarm"`
	KeybaseURL string   `mapstructure:"keybase_url"`
}

// Node wraps both endpoints for Tendermint RPC Node and REST API Server.
type Node struct {
	RPCNode      string `mapstructure:"rpc_node"`
	LCDEndpoint  string `mapstructure:"lcd_endpoint"`
	GRPCEndpoint string `mapstructure:"grpc_endpoint"`
	NetworkType  string
}

// Database wraps PostgreSQL database config.
type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Table    string `mapstructure:"table"`
}

// Alarm wraps push notification alarm config.
type Alarm struct {
	PushServerEndpoint string `mapstructure:"push_server_endpoint"`
	Switch             bool   `mapstructure:"switch"`
}

// ParseConfig attempts to read and parse config.yaml from the given path.
// An error reading or parsing the config results in a panic.
func ParseConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")                                                     // for test cases
	viper.AddConfigPath(os.Getenv("HOME") + "/cosmostation-cosmos/chain-exporter") // for production

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	if viper.GetString("network_type") == "" {
		log.Fatal("define network_type param in your config file.")
	}

	var config Config
	sub := viper.Sub(viper.GetString("network_type"))
	sub.Unmarshal(&config)

	config.KeybaseURL = viper.GetString("keybase_endpoint")

	// This code is used in main.go to log network type when starting server.
	if viper.GetString("network_type") == "mainnet" {
		config.Node.NetworkType = viper.GetString("network_type")
	} else {
		config.Node.NetworkType = viper.GetString("network_type")
	}

	return &config
}
