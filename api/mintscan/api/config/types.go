package config

type Config struct {
	Node   *NodeConfig
	DB     *DBConfig
	Web    *WebConfig
	Market *MarketConfig
}

type NodeConfig struct {
	GaiadURL string
	LCDURL   string
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Table    string
}

type WebConfig struct {
	Port string
}

type MarketConfig struct {
	CoinmarketCap struct {
		URL    string
		CoinID string
		APIKey string
	}
	CoinGecko struct {
		URL string
	}
}
