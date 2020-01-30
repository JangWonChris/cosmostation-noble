package config

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config wraps all configs that are used in this project
type Config struct {
	Node   NodeConfig
	DB     DBConfig
	Web    WebConfig
	Market MarketConfig
	Denom  string
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

// WebConfig wraps port number of this project
type WebConfig struct {
	Port string `yaml:"port"`
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

// NewConfig configures configuration.
// Mainnet | Dev | Testnet
func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to read config"))
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
		cfg.Web = WebConfig{
			Port: viper.GetString("mainnet.web.port"),
		}
		cfg.Denom = "uatom"
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
		cfg.Web = WebConfig{
			Port: viper.GetString("dev.web.port"),
		}
		cfg.Denom = "uatom"
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
		cfg.Web = WebConfig{
			Port: viper.GetString("testnet.web.port"),
		}
		cfg.Denom = "muon"
	default:
		log.Fatal("active can be either mainnet or testnet.")
	}

	cfg.Market.CoinGecko.Endpoint = viper.GetString("market.coingecko.endpoint")
	cfg.Market.CoinmarketCap.Endpoint = viper.GetString("market.coinmarketcap.endpoint")
	cfg.Market.CoinmarketCap.CoinID = viper.GetString("market.coinmarketcap.coin_id")
	cfg.Market.CoinmarketCap.APIKey = viper.GetString("market.coinmarketcap.api_key")

	return &cfg
}
