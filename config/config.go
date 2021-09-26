package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AppName     string `envconfig:"APP_NAME" default:"walrus"`
	Host        string `envconfig:"HOST" default:"0.0.0.0"`
	Port        string `envconfig:"PORT" default:"3000"`
	Environment string `envconfig:"ENVIRONMENT" default:"dev"`
}

func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
