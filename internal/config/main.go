package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	TraefikUrl      string `env:"TRAEFIK_URL"`
	TraefikUsername string `env:"TRAEFIK_USERNAME"`
	TraefikPassword string `env:"TRAEFIK_PASSWORD"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"Info"`

	// set twice as long as any blackbox alerts
	RunInterval int `env:"INTERVAL" envDefault:"600"`

	TargetsFile string `env:"TARGETS_FILE" envDefault:"/blackbox-traefik-sd/targets.json"`
}

func New() (Config, error) {
	conf := Config{}

	if err := env.Parse(&conf); err != nil {
		return Config{}, fmt.Errorf("Error parsing config from env: %+v\n", err)
	}

	return conf, nil
}
