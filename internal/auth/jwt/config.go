package jwt

import (
	"time"
)

type Config struct {
	SecretKey       string        `yaml:"secret_key"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
}
