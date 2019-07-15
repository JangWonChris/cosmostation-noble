package config

type Config struct {
	Node *NodeConfig
	DB   *DBConfig
	JWT  *JWTConfig
	Web  *WebConfig
}

type NodeConfig struct {
	GaiadURL string
	LcdURL   string
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

type JWTConfig struct {
	Token string
}
