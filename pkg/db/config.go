package db

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `default:"localhost" envconfig:"POSTGRES_HOST"`
	Port     int    `default:"5433" envconfig:"POSTGRES_PORT"`
	User     string `default:"any" envconfig:"POSTGRES_USER"`
	Password string `default:"any" envconfig:"POSTGRES_PASSWORD"`
	DBName   string `default:"bookmark" envconfig:"POSTGRES_DB"`
	SSLMode  string `default:"disable" envconfig:"POSTGRES_SSL_MODE"`
	TimeZone string `default:"UTC" envconfig:"POSTGRES_TIMEZONE"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.TimeZone,
	)
}
