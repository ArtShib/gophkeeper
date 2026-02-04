package models

type ConfigCrypto struct {
	SaltClient []byte
	SaltServer []byte
}

type ConfigTLS struct {
	Cert string `env:"TLS_CERT"`
	Key  string `env:"TLS_KEY"`
}
