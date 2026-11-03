package config

import (
	"uuid"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        int    `envconfig:"PORT"`
}

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg, nil
}
