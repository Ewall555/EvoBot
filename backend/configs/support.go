package configs

type SupportConfig struct {
	Baseurl  string `mapstructure:"baseurl"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
