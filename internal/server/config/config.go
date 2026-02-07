package config

import (
	"context"
	"flag"
	"log/slog"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Config структура конфига
type Config struct {
	ConfigStore *models.ConfigStore
	ConfigJWT   *models.ConfigJWT
	ConfigGRPC  *models.ConfigGRPC
	logger      *slog.Logger
}

// LoadConfigEnv загрузка данных в конфиг из env
func (c *Config) LoadConfigEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := env.Parse(c.ConfigStore); err != nil {
		return err
	}
	if err := env.Parse(c.ConfigGRPC); err != nil {
		return err
	}
	if err := env.Parse(c.ConfigJWT); err != nil {
		return err
	}
	return nil
}

// LoadConfigFlag загрузка данных в конфиг из cmd
func (c *Config) LoadConfigFlag() {
	if c.ConfigGRPC.Port == 0 {
		flag.IntVar(&c.ConfigGRPC.Port, "p", 3030, "gRPC port")
	}
	if c.ConfigStore.DatabaseDSN == "" {
		flag.StringVar(&c.ConfigStore.DatabaseDSN, "d", "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable", "DataBase connection string")
	}
	flag.Parse()
}

// MustLoadConfig конструктор Config
func MustLoadConfig(ctx context.Context, logger *slog.Logger) (*Config, error) {
	logHelper := loghelper.New(logger, "config.MustLoadConfig")
	var err error
	cfg := Config{
		ConfigGRPC:  &models.ConfigGRPC{},
		ConfigStore: &models.ConfigStore{},
		ConfigJWT: &models.ConfigJWT{
			TokenTTLMIN: time.Minute * 30,
			SecretKey:   []byte("048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21"),
		},
	}

	err = cfg.LoadConfigEnv()
	if err != nil {
		logHelper.LogError(ctx, "LoadConfigEnv", err)
	}
	cfg.LoadConfigFlag()

	return &cfg, err
}
