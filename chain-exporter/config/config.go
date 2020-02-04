package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config wraps all configs that are used in this project
type Config struct {
	Node       NodeConfig
	DB         DBConfig
	Alarm      AlarmConfig
	ES         ESConfig
	Market     MarketConfig
	KeybaseURL string
}

// NodeConfig wraps both endpoints for Tendermint RPC Node and REST API Server
type NodeConfig struct {
	RPCNode     string `yaml:"rpc_node"`
	LCDEndpoint string `yaml:"lcd_endpoint"`
}

// DBConfig wraps PostgreSQL database config
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Table    string `yaml:"table"`
}

// AlarmConfig wraps push notification alarm config
type AlarmConfig struct {
	PushServerEndpoint string
	Switch             bool
}

// MarketConfig wraps endpoints for CoinmarketCap and CoinGecko
// In this project, we primarily use CoinGecko
type MarketConfig struct {
	CoinGecko struct {
		Endpoint string `yaml:"endpoint"`
	} `yaml:"coingecko"`
	CoinmarketCap struct {
		Endpoint string `yaml:"endpoint"`
		CoinID   string `yaml:"coin_id"`
		APIKey   string `yaml:"api_key"`
	} `yaml:"coinmarketcap"`
}

type (
	// [NOT USED]
	ESConfig struct {
		Endpoint  string `yaml:"endpoint"`
		Region    string `yaml:"region"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Sniff     bool   `yaml:"sniff"`
	}

	// [NOT USED]
	RavenConfig struct {
		RavenDSN string `yaml:"raven_dsn"`
		Address  string `yaml:"address"`
		Period   string `yaml:"period"`
	}
)

// ParseConfig configures configuration
func ParseConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	if viper.GetString("active") == "" {
		log.Fatal("define active param in your config file.")
	}

	cfg := Config{}

	switch viper.GetString("active") {
	case "mainnet":
		cfg.Node = NodeConfig{
			RPCNode:     viper.GetString("mainnet.node.rpc_node"),
			LCDEndpoint: viper.GetString("mainnet.node.lcd_endpoint"),
		}
		cfg.DB = DBConfig{
			Host:     viper.GetString("mainnet.database.host"),
			Port:     viper.GetString("mainnet.database.port"),
			User:     viper.GetString("mainnet.database.user"),
			Password: viper.GetString("mainnet.database.password"),
			Table:    viper.GetString("mainnet.database.table"),
		}
		cfg.Alarm = AlarmConfig{
			PushServerEndpoint: viper.GetString("mainnet.alarm.push_server_endpoint"),
			Switch:             viper.GetBool("mainnet.alarm.switch"),
		}
	case "dev":
		cfg.Node = NodeConfig{
			RPCNode:     viper.GetString("dev.node.rpc_node"),
			LCDEndpoint: viper.GetString("dev.node.lcd_endpoint"),
		}
		cfg.DB = DBConfig{
			Host:     viper.GetString("dev.database.host"),
			Port:     viper.GetString("dev.database.port"),
			User:     viper.GetString("dev.database.user"),
			Password: viper.GetString("dev.database.password"),
			Table:    viper.GetString("dev.database.table"),
		}
		cfg.Alarm = AlarmConfig{
			PushServerEndpoint: viper.GetString("dev.alarm.push_server_endpoint"),
			Switch:             viper.GetBool("dev.alarm.switch"),
		}
	case "testnet":
		cfg.Node = NodeConfig{
			RPCNode:     viper.GetString("testnet.node.rpc_node"),
			LCDEndpoint: viper.GetString("testnet.node.lcd_endpoint"),
		}
		cfg.DB = DBConfig{
			Host:     viper.GetString("testnet.database.host"),
			Port:     viper.GetString("testnet.database.port"),
			User:     viper.GetString("testnet.database.user"),
			Password: viper.GetString("testnet.database.password"),
			Table:    viper.GetString("testnet.database.table"),
		}
		cfg.Alarm = AlarmConfig{
			PushServerEndpoint: viper.GetString("testnet.alarm.push_server_endpoint"),
			Switch:             viper.GetBool("testnet.alarm.switch"),
		}
	default:
		fmt.Println("define active params in config.yaml")
	}

	return &cfg
}
