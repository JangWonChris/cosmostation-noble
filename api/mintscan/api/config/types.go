package config

type Config struct {
	Node *NodeConfig
	DB   *DBConfig
	Web  *WebConfig
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

