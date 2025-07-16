package configs

import "time"

type AuthConfig struct {
	JwtExpired time.Duration `mapstructure:"jwt_expired"`
	Secret     string        `mapstructure:"secret"`
}
