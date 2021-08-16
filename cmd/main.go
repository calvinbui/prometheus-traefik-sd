package main

import (
	"fmt"
	"time"

	"github.com/calvinbui/prometheus-traefik-sd/internal/config"
	"github.com/calvinbui/prometheus-traefik-sd/internal/helpers"
	"github.com/calvinbui/prometheus-traefik-sd/internal/logger"
	"github.com/calvinbui/prometheus-traefik-sd/internal/prometheus"
	"github.com/calvinbui/prometheus-traefik-sd/internal/traefik"
)

func main() {
	logger.Init()

	logger.Info("Loading config")
	conf, err := config.New()
	if err != nil {
		logger.Fatal("Error parsing config", err)
	}

	logger.Info("Setting log level to " + conf.LogLevel)
	if err = logger.SetLevel(conf.LogLevel); err != nil {
		logger.Fatal("Error setting log level", err)
	}

	var graceFiles []helpers.GraceFile

	for {
		logger.Info("Getting Traefik routers and hosts")
		logger.Debug("Getting Traefik routes")
		routes, err := traefik.GetRoutingRules(conf.TraefikUrl, conf.TraefikUsername, conf.TraefikPassword)
		if err != nil {
			logger.Fatal("Error getting Traefik routing rules", err)
		}
		logger.Debug("Getting Traefik rules from the routes")
		rules, err := helpers.GetRouteRules(routes)
		if err != nil {
			logger.Fatal("Error parsing Traefik routes", err)
		}
		logger.Debug(fmt.Sprintf("Rules: %+v", rules))

		logger.Debug("Getting hosts from rules")
		hosts := helpers.GetHostsFromRules(rules)
		logger.Debug(fmt.Sprintf("All hosts: %+v", hosts))

		logger.Info("Creating Prometheus JSON files")
		logger.Debug("Generating Prometheus data")
		allTargets := []helpers.PromTargetFile{}
		for _, t := range hosts {
			logger.Debug(fmt.Sprintf("Adding targets %+v", t))
			allTargets = append(allTargets, helpers.PromTargetFile{
				FilePath: helpers.CreateFileName(conf.OutputDir, t),
				Data: prometheus.TargetGroups{
					{
						Targets: t,
					},
				},
			})
		}
		logger.Debug(fmt.Sprintf("Prometheus data: %+v", allTargets))

		logger.Debug("Put JSON files")
		if err = helpers.CreateJSON(allTargets, conf.OutputDir); err != nil {
			logger.Fatal("Error writing to JSON file", err)
		}

		logger.Info(fmt.Sprintf("Finding and deleting JSON files past grace period (%v)", conf.GracePeriod))
		graceFiles, err = helpers.DeleteOldTargets(allTargets, conf.OutputDir, graceFiles, conf.GracePeriod)
		if err != nil {
			logger.Fatal("Error deleting JSON files", err)
		}

		logger.Info(fmt.Sprintf("Sleeping %v seconds until next run", conf.RunInterval))
		time.Sleep(time.Duration(conf.RunInterval) * time.Second)
	}
}
