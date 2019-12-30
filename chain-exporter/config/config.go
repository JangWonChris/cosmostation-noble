package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Node          *NodeConfig
	DB            *DBConfig
	ES            *ESConfig
	Raven         *RavenConfig
	Coinmarketcap *CoinmarketcapConfig
	KeybaseURL    string
}

type (
	NodeConfig struct {
		GaiadURL string
		LCDURL   string
	}

	ESConfig struct {
		URL       string
		Region    string
		AccessKey string
		SecretKey string
		Sniff     bool
	}

	DBConfig struct {
		Host     string
		User     string
		Password string
		Table    string
	}

	RavenConfig struct {
		RavenDSN string
		Address  string
		Period   string
	}

	CoinmarketcapConfig struct {
		CoinID string
		APIKey string
	}
)

// NewConfig configures configuration
func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/home/ubuntu/cosmostation-cosmos/chain-exporter")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	config := &Config{}
	config.KeybaseURL = viper.GetString("keybase_url")
	nodeConfig := &NodeConfig{}
	dbConfig := &DBConfig{}

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
		fmt.Println("define active params in config.yaml")
	}

	config.Node = nodeConfig
	config.DB = dbConfig

	return config
}
