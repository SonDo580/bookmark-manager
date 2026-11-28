package config

import (
	"time"
	"uuid"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark_service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        int    `envconfig:"PORT" default:"8080"`
	BasePath    string `envconfig:"BASE_PATH" default:"/"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"INFO"`

	PrivateKeyPath string        `envconfig:"PRIVATE_KEY_PATH" default:"./private.pem"`
	PublicKeyPath  string        `envconfig:"PUBLIC_KEY_PATH" default:"./public.pem"`
	TokenExp       time.Duration `envconfig:"TOKEN_EXP" default:"24h"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg, nil
}
