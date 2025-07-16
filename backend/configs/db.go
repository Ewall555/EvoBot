package configs

type DBConfig struct {
	Engine   string `mapstructure:"engine"` // postgres/mysql
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"` // Only PostgreSQL
}
