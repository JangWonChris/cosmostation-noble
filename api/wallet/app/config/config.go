package config

// Configuration for Cosmos Mainnet
func GetMainnetConfig() *Config {
	return &Config{
		Node: &NodeConfig{
			GaiadURL: "http://52.78.163.36:26657",
			LCDUrl:   "https://lcd.cosmostation.io",
		},
		DB: &DBConfig{
			Host:     "cosmostation.ci6bhjszmrb3.ap-northeast-2.rds.amazonaws.com:5432",
			User:     "root",
			Password: "CosmosGo00!!",
			Table:    "mainnet",
		},
		JWT: &JWTConfig{
			Token: "SecureCosmostation!@#50291230728",
		},
	}
}

// Configuration for Cosmos Mainnet
func GetMainnetDevConfig() *Config {
	return &Config{
		Node: &NodeConfig{
			GaiadURL: "http://52.78.163.36:26657",
			LCDUrl:   "https://lcd-mainnet-dev.cosmostation.io",
		},
		DB: &DBConfig{
			Host:     "cosmostation-dev.ci6bhjszmrb3.ap-northeast-2.rds.amazonaws.com:5432",
			User:     "root",
			Password: "dnpfqldwhr00!!",
			Table:    "dev-test",
		},
		JWT: &JWTConfig{
			Token: "SecureCosmostation!@#50291230728",
		},
	}
}
