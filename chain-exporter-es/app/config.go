package app

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	RPCEndPoint string
}

func validateBasic(config *Config) error {
	if len(config.RPCEndPoint) == 0 {
		return errors.New("ES RPCEndPoint not set")
	}
	return nil
}

func InitConfig(network string, env string) (*Config, error)  {
	config := &Config{
		RPCEndPoint:viper.GetString(fmt.Sprintf("%s.%s.RPCEndPoint", network, env)),
	}



	err := validateBasic(config)
	if err != nil {
		return nil, err
	}

	return config, err
}