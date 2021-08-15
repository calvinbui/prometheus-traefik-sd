package config

import (
	"github.com/jessevdk/go-flags"
)

type Config struct {
	TraefikUrl      string `long:"traefik-url" env:"TRAEFIK_URL"`
	TraefikUsername string `long:"traefik-username" env:"TRAEFIK_USERNAME"`
	TraefikPassword string `long:"traefik-password" env:"TRAEFIK_PASSWORD"`

	LogLevel string `long:"log-level" env:"LOG_LEVEL" default:"Info"`

	RunInterval int `long:"interval" env:"INTERVAL" default:"600"`

	OutputDir string `short:"o" env:"OUTPUT_DIR" default:"/prometheus-traefik-sd/"`
}

func New() (Config, error) {
	conf := Config{}
	parser := flags.NewParser(&conf, flags.Default)

	if _, err := parser.Parse(); err != nil {
		return Config{}, err
	}

	return conf, nil
}
