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

	logger.Debug("Loading config")
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
		logger.Info("Getting Traefik routing rules")
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

		logger.Info("Getting hosts from rules")
		hosts := helpers.GetHostsFromRules(rules)
		logger.Debug(fmt.Sprintf("All hosts: %+v", hosts))

		logger.Info("Generating Prometheus data")
		allTargets := []helpers.PromTargetFile{}
		for _, t := range hosts {
			logger.Info(fmt.Sprintf("Adding targets %+v", t))
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

		logger.Info("Creating Prometheus JSON target file")
		if err = helpers.CreateJSON(allTargets, conf.OutputDir); err != nil {
			logger.Fatal("Error writing to JSON file", err)
		}

		logger.Info("Deleting old target files past grace period")
		graceFiles, err = helpers.DeleteOldTargets(allTargets, conf.OutputDir, graceFiles, conf.GracePeriod)
		if err != nil {
			logger.Fatal("Error deleting old target files", err)
		}

		logger.Info(fmt.Sprintf("Sleeping %v seconds until next run", conf.RunInterval))
		time.Sleep(time.Duration(conf.RunInterval) * time.Second)
	}
}
