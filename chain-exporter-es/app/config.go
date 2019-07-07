package app

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	kavaApp "github.com/kava-labs/kava/app"

)

type Config struct {
	RPCEndPoint string
	sdkConfig *sdkTypes.Config
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
		sdkConfig:sdkTypes.GetConfig(),
	}

	if network == "kava" {
		sdkConfig := sdkTypes.GetConfig()
		kavaApp.SetBech32AddressPrefixes(sdkConfig)
		sdkConfig.Seal()

		config.sdkConfig = sdkConfig
	}


	err := validateBasic(config)
	if err != nil {
		return nil, err
	}

	return config, err
}