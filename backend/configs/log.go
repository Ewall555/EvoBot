package configs

type LogConfig struct {
	Lang       string `mapstructure:"lang"`
	Model      string `mapstructure:"model"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
