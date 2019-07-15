package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewAPIConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("/home/ubuntu/cosmostation-cosmos/app/wallet/") // call multiple times to add many search paths

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config := &Config{}
	nodeConfig := &NodeConfig{}
	dbConfig := &DBConfig{}
	webConfig := &WebConfig{}
	jwtConfig := &JWTConfig{}

	// Production or Development
	switch viper.GetString("active") {
	case "prod":
		nodeConfig.GaiadURL = viper.GetString("prod.node.gaiad_url")
		nodeConfig.LcdURL = viper.GetString("prod.node.lcd_url")
		dbConfig.Host = viper.GetString("prod.database.host")
		dbConfig.User = viper.GetString("prod.database.user")
		dbConfig.Password = viper.GetString("prod.database.password")
		dbConfig.Table = viper.GetString("prod.database.table")
		webConfig.Port = viper.GetString("prod.port")
		jwtConfig.Token = viper.GetString("prod.jwt.token")
	case "dev":
		nodeConfig.GaiadURL = viper.GetString("dev.node.gaiad_url")
		nodeConfig.LcdURL = viper.GetString("dev.node.lcd_url")
		dbConfig.Host = viper.GetString("dev.database.host")
		dbConfig.User = viper.GetString("dev.database.user")
		dbConfig.Password = viper.GetString("dev.database.password")
		dbConfig.Table = viper.GetString("dev.database.table")
		webConfig.Port = viper.GetString("dev.port")
		jwtConfig.Token = viper.GetString("dev.jwt.token")
	default:
		fmt.Println("Define active params in config.yaml")
	}

	config.Node = nodeConfig
	config.DB = dbConfig
	config.Web = webConfig
	config.JWT = jwtConfig

	return config
}
