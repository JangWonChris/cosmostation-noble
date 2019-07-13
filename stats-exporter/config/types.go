package config

type Config struct {
	Node             *NodeConfig
	Raven            *RavenConfig
	DB               *DBConfig
	ES               *ESConfig
	Coinmarketcap    *CoinmarketcapConfig
	KeybaseURL       string
	CoinmarketcapURL string
	CoinGeckoURL     string
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
	CoinmarketcapConfig struct {
		CoinID string
		APIKey string
	}
)
