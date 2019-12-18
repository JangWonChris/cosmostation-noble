package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Node *NodeConfig
	DB   *DBConfig
	Web  *WebConfig
	JWT  *JWTConfig
}

type (
	NodeConfig struct {
		GaiadURL string
		LCDURL   string
	}

	DBConfig struct {
		Host     string
		User     string
		Password string
		Table    string
	}

	WebConfig struct {
		Port string
	}

	JWTConfig struct {
		Token string
	}
)

// NewConfig configures configuration
func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}

	config := &Config{}
	nodeConfig := &NodeConfig{}
	dbConfig := &DBConfig{}
	webConfig := &WebConfig{}
	jwtConfig := &JWTConfig{}

	// configuration for prod, dev, testnet
	switch viper.GetString("active") {
	case "prod":
		nodeConfig.GaiadURL = viper.GetString("prod.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("prod.node.lcd_url")
		dbConfig.Host = viper.GetString("prod.database.host")
		dbConfig.User = viper.GetString("prod.database.user")
		dbConfig.Password = viper.GetString("prod.database.password")
		dbConfig.Table = viper.GetString("prod.database.table")
		webConfig.Port = viper.GetString("prod.port")
		jwtConfig.Token = viper.GetString("prod.jwt.token")
	case "dev":
		nodeConfig.GaiadURL = viper.GetString("dev.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("dev.node.lcd_url")
		dbConfig.Host = viper.GetString("dev.database.host")
		dbConfig.User = viper.GetString("dev.database.user")
		dbConfig.Password = viper.GetString("dev.database.password")
		dbConfig.Table = viper.GetString("dev.database.table")
		webConfig.Port = viper.GetString("dev.port")
		jwtConfig.Token = viper.GetString("dev.jwt.token")
	case "testnet":
		nodeConfig.GaiadURL = viper.GetString("testnet.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("testnet.node.lcd_url")
		dbConfig.Host = viper.GetString("testnet.database.host")
		dbConfig.User = viper.GetString("testnet.database.user")
		dbConfig.Password = viper.GetString("testnet.database.password")
		dbConfig.Table = viper.GetString("testnet.database.table")
		webConfig.Port = viper.GetString("testnet.port")
		jwtConfig.Token = viper.GetString("testnet.jwt.token")
	default:
		fmt.Println("Define active params in config.yaml")
	}

	config.Node = nodeConfig
	config.DB = dbConfig
	config.Web = webConfig
	config.JWT = jwtConfig

	return config
}
