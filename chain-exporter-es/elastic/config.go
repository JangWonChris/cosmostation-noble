package elastic

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ElasticHost string
	Region      string
	AccessKey   string
	SecretKey   string
	Sniff       bool
	RPCEndPoint string
}

func validateBasic(config *Config) error {
	if len(config.ElasticHost) == 0 {
		return errors.New("ElasticHost not set")
	}
	if len(config.Region) == 0 {
		return errors.New("ES Region not set")
	}
	if len(config.AccessKey) == 0 {
		return errors.New("ES AccessKey not set")
	}
	if len(config.SecretKey) == 0 {
		return errors.New("ES SecretKey not set")
	}
	if len(config.RPCEndPoint) == 0 {
		return errors.New("ES RPCEndPoint not set")
	}
	return nil
}

func InitConfig(network string, env string) (*Config, error)  {
	logrus.Info("viper test", viper.GetString(fmt.Sprintf("%s.%s.ElasticHost", network, env)))
	config := &Config{
		ElasticHost:viper.GetString(fmt.Sprintf("%s.%s.ElasticHost", network, env)),
		Region:viper.GetString(fmt.Sprintf("%s.%s.Region", network, env)),
		AccessKey:viper.GetString(fmt.Sprintf("%s.%s.AccessKey", network, env)),
		SecretKey:viper.GetString(fmt.Sprintf("%s.%s.SecretKey", network, env)),
		RPCEndPoint:viper.GetString(fmt.Sprintf("%s.%s.RPCEndPoint", network, env)),
		Sniff:viper.GetBool(fmt.Sprintf("%s.%s.Sniff", network, env)),
	}

	err := validateBasic(config)
	if err != nil {
		return nil, err
	}
	return config, err
}