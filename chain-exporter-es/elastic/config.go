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

func InitConfig(env string) (*Config, error)  {
	logrus.Info("viper test", viper.GetString(fmt.Sprintf("%s.ElasticHost", env)))
	config := &Config{
		ElasticHost:viper.GetString(fmt.Sprintf("%s.ElasticHost", env)),
		Region:viper.GetString(fmt.Sprintf("%s.Region", env)),
		AccessKey:viper.GetString(fmt.Sprintf("%s.AccessKey", env)),
		SecretKey:viper.GetString(fmt.Sprintf("%s.SecretKey", env)),
		RPCEndPoint:viper.GetString(fmt.Sprintf("%s.RPCEndPoint", env)),
		Sniff:viper.GetBool(fmt.Sprintf("%s.Sniff", env)),
	}

	err := validateBasic(config)
	if err != nil {
		return nil, err
	}
	return config, err
}