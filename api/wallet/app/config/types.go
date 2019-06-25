package config

type Config struct {
	Node *NodeConfig
	DB   *DBConfig
	JWT  *JWTConfig
}

type NodeConfig struct {
	GaiadURL string
	LCDUrl   string
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Table    string
}
type JWTConfig struct {
	Token string
}
