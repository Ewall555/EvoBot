package configs

type ServerConfig struct {
	SystemConfig  SystemConfig   `mapstructure:"system"`
	AuthConfig    *AuthConfig    `mapstructure:"auth"`
	LogConfig     *LogConfig     `mapstructure:"log"`
	RedisConfig   *RedisConfig   `mapstructure:"redis"`
	DBConfig      *DBConfig      `mapstructure:"db"`
	SupportConfig *SupportConfig `mapstructure:"support"`
}
