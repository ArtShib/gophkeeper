package config

import (
	"log/slog"
	"os"

	"github.com/ArtShib/gophkeeper/internal/models"
)

type Config struct {
	ConfigStore  *models.ConfigStore
	ConfigGRPC   *models.ConfigGRPC
	ConfigTLS    *models.ConfigTLS
	ConfigCrypto *models.ConfigCrypto
	logger       *slog.Logger
}

func MustLoadConfig() *Config {
	return &Config{
		ConfigGRPC: &models.ConfigGRPC{
			Server: "localhost:3030",
		},
		ConfigStore: &models.ConfigStore{
			DatabaseDSN: "test.db",
		},
		ConfigTLS: &models.ConfigTLS{
			Cert: os.Getenv("TLS_CERT"),
			Key:  os.Getenv("TLS_KEY"),
			CA:   os.Getenv("TLS_CA"),
		},
		ConfigCrypto: &models.ConfigCrypto{
			SaltClient: []byte("saltxmcxl,mv;xcmv.,xcmv.,xmc.mxzc.zcZX"),
			SaltServer: []byte("sdlaksjdlaksdasdasd,mmm,san,dnmasd908"),
		},
	}
}
