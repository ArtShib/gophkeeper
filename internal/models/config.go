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

type ConfigCrypto struct {
	SaltClient []byte
	SaltServer []byte
}

type ConfigTLS struct {
	Cert string `env:"TLS_CERT"`
	Key  string `env:"TLS_KEY"`
	CA   string `env:"TLS_CA"`
}
