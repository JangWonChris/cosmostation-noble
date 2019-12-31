package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Node  *NodeConfig
	DB    *DBConfig
	Web   *WebConfig
	Alarm *AlarmConfig
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
		Port     string
		JWTToken string
	}

	AlarmConfig struct {
		PushServerURL string
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
	alarmConfig := &AlarmConfig{}

	// configuration for prod, dev, testnet
	switch viper.GetString("active") {
	case "prod":
		nodeConfig.GaiadURL = viper.GetString("prod.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("prod.node.lcd_url")
		dbConfig.Host = viper.GetString("prod.database.host")
		dbConfig.User = viper.GetString("prod.database.user")
		dbConfig.Password = viper.GetString("prod.database.password")
		dbConfig.Table = viper.GetString("prod.database.table")
		webConfig.Port = viper.GetString("prod.web.port")
		webConfig.JWTToken = viper.GetString("prod.web.jwt_token")
		alarmConfig.PushServerURL = viper.GetString("prod.alarm.push_server_url")
	case "dev":
		nodeConfig.GaiadURL = viper.GetString("dev.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("dev.node.lcd_url")
		dbConfig.Host = viper.GetString("dev.database.host")
		dbConfig.User = viper.GetString("dev.database.user")
		dbConfig.Password = viper.GetString("dev.database.password")
		dbConfig.Table = viper.GetString("dev.database.table")
		webConfig.Port = viper.GetString("dev.web.port")
		webConfig.JWTToken = viper.GetString("dev.web.jwt_token")
		alarmConfig.PushServerURL = viper.GetString("dev.alarm.push_server_url")
	case "testnet":
		nodeConfig.GaiadURL = viper.GetString("testnet.node.gaiad_url")
		nodeConfig.LCDURL = viper.GetString("testnet.node.lcd_url")
		dbConfig.Host = viper.GetString("testnet.database.host")
		dbConfig.User = viper.GetString("testnet.database.user")
		dbConfig.Password = viper.GetString("testnet.database.password")
		dbConfig.Table = viper.GetString("testnet.database.table")
		webConfig.Port = viper.GetString("testnet.web.port")
		webConfig.JWTToken = viper.GetString("testnet.web.jwt_token")
		alarmConfig.PushServerURL = viper.GetString("testnet.alarm.push_server_url")
	default:
		fmt.Println("define active params in config.yaml")
	}

	config.Node = nodeConfig
	config.DB = dbConfig
	config.Web = webConfig
	config.Alarm = alarmConfig

	return config
}
