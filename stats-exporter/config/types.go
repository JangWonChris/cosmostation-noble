package config

type Config struct {
	Node   *NodeConfig
	Raven  *RavenConfig
	DB     *DBConfig
	ES     *ESConfig
	Market *MarketConfig
}

type (
	NodeConfig struct {
		GaiadURL string
		LCDURL   string
	}

	RavenConfig struct {
		RavenDSN string
		Address  string
		Period   string
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

	MarketConfig struct {
		CoinmarketCap struct {
			URL    string
			CoinID string
			APIKey string
		}
		CoinGecko struct {
			URL string
		}
	}
)
