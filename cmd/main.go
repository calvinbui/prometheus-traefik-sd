package main

import (
	"fmt"
	"time"

	"github.com/calvinbui/blackbox-traefik-sd/internal/config"
	"github.com/calvinbui/blackbox-traefik-sd/internal/logger"
	"github.com/calvinbui/blackbox-traefik-sd/internal/traefik"
)

func main() {
	logger.Init()

	logger.Debug("Loading internal config")
	conf, err := config.New()
	if err != nil {
		logger.Fatal("Error parsing config", err)
	}

	logger.Info("Setting log level to " + conf.LogLevel)
	err = logger.SetLevel(conf.LogLevel)
	if err != nil {
		logger.Fatal("Error setting log level", err)
	}

	logger.Debug("Creating Traefik client")
	cTraefik := traefik.Client{
		Url:      conf.TraefikUrl,
		Username: conf.TraefikUsername,
		Password: conf.TraefikPassword,
	}

	for {
		logger.Info("Getting Traefik routes")
		routes, err := cTraefik.GetRoutes()
		if err != nil {
			logger.Fatal("Error getting Traefik routes", err)
		}

		targets := traefik.GetHostsFromRouter(routes)

		logger.Debug(fmt.Sprintf("Targets: %+v", targets))

		logger.Info("Sleeping until next run")
		time.Sleep(60 * time.Second)
	}
}
