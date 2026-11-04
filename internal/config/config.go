package config

import (
	"uuid"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark_service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        int    `envconfig:"PORT" default:"8080"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return &cfg, nil
}
