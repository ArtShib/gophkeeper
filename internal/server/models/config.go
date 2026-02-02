package models

import "time"

type ConfigJWT struct {
	TokenTTLMIN time.Duration `env:"TOKEN_TTL_MIN"`
	SecretKey   []byte
}

type ConfigStore struct {
	DatabaseDSN string `env:"DATABASE_DSN"`
}

type ConfigGRPC struct {
	Port int `env:"PORT"`
}
