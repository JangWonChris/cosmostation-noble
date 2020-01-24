package config

import (
	"log"

	"github.com/pkg/errors"

	"github.com/spf13/viper"
)

// Config defines all necessary juno configuration parameters.
type Config struct {
	Node NodeConfig `yaml:"node"`
	DB   DBConfig   `yaml:"database"`
}

// NodeConfig defines endpoints for both RPC node and LCD REST API server
type NodeConfig struct {
	RPCNode     string `yaml:"rpc_node"`
	LCDEndpoint string `yaml:"lcd_endpoint"`
}

// DBConfig defines all database connection configuration parameters.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Table    string `yaml:"table"`
}

// ParseConfig attempts to read and parse chain-exporter config from the given configPath.
// An error reading or parsing the config results in a panic.
func ParseConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

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
	default:
		log.Fatal("active can be either mainnet or dev.")
	}

	return cfg
}
