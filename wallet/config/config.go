package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config wraps all configs that are used in this project.
type Config struct {
	Node  NodeConfig  `mapstructure:"node"`
	DB    DBConfig    `mapstructure:"database"`
	Web   WebConfig   `mapstructure:"web"`
	Alarm AlarmConfig `mapstructure:"alarm"`
}

// NodeConfig wraps both endpoints for Tendermint RPC Node and REST API Server.
type NodeConfig struct {
	RPCNode     string `mapstructure:"rpc_node"`
	LCDEndpoint string `mapstructure:"lcd_endpoint"`
	NetworkType string
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
	Port     string `mapstructure:"port"`
	JWTToken string `mapstructure:"jwt_token"`
}

// AlarmConfig wraps push notification server endpoint of this project.
type AlarmConfig struct {
	PushServerURL string
}

// ParseConfig attempts to read and parse config.yaml from the given path
// An error reading or parsing the config results in a panic.
func ParseConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")                                      // for test cases
	viper.AddConfigPath("/home/ubuntu/cosmostation-cosmos/wallet/") // for production

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	if viper.GetString("network_type") == "" {
		log.Fatal("define network_type param in your config file.")
	}

	var config Config
	sub := viper.Sub(viper.GetString("network_type"))
	sub.Unmarshal(&config)

	// This code is used in main.go to log network type when starting server.
	if viper.GetString("network_type") == "mainnet" {
		config.Node.NetworkType = viper.GetString("network_type")
	} else {
		config.Node.NetworkType = viper.GetString("network_type")
	}

	return &config
}
