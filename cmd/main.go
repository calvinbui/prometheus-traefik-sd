package main

import (
	"fmt"
	"time"

	"github.com/calvinbui/blackbox-traefik-sd/internal/config"
	"github.com/calvinbui/blackbox-traefik-sd/internal/helpers"
	"github.com/calvinbui/blackbox-traefik-sd/internal/logger"
	"github.com/calvinbui/blackbox-traefik-sd/internal/prometheus"
	"github.com/calvinbui/blackbox-traefik-sd/internal/traefik"
)

func main() {
	logger.Init()

	logger.Debug("Loading config")
	conf, err := config.New()
	if err != nil {
		logger.Fatal("Error parsing config", err)
	}

	logger.Info("Setting log level to " + conf.LogLevel)
	if err = logger.SetLevel(conf.LogLevel); err != nil {
		logger.Fatal("Error setting log level", err)
	}

	for {
		logger.Info("Getting Traefik routing rules")
		rules, err := traefik.GetRoutingRules(conf.TraefikUrl, conf.TraefikUsername, conf.TraefikPassword)
		if err != nil {
			logger.Fatal("Error getting Traefik routing rules", err)
		}
		logger.Debug(fmt.Sprintf("Rules: %+v", rules))

		logger.Info("Getting hosts from rules")
		hosts := helpers.GetHostsFromRules(rules)
		logger.Debug(fmt.Sprintf("All hosts: %+v", hosts))

		logger.Info("Generating Prometheus data")
		tg := []prometheus.TargetGroup{}
		for _, t := range hosts {
			tg = append(tg, prometheus.TargetGroup{Targets: t})
		}

		logger.Debug("Creating config folder if it does not exist")
		if err = helpers.InitFolder(conf.TargetsFile); err != nil {
			logger.Fatal("Error creating config folder", err)
		}

		logger.Info("Creating Prometheus JSON target file")
		if err = helpers.CreateTargetsJSON(tg, conf.TargetsFile); err != nil {
			logger.Fatal("Error writing to JSON file", err)
		}

		logger.Info(fmt.Sprintf("Sleeping %v seconds until next run", conf.RunInterval))
		time.Sleep(time.Duration(conf.RunInterval) * time.Second)
	}
}
