package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
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

	logger.Debug("Building Traefik rules regex")
	var re = regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)

	for {
		logger.Info("Getting Traefik routing rules")
		rules, err := traefik.GetRoutingRules(conf.TraefikUrl, conf.TraefikUsername, conf.TraefikPassword)
		if err != nil {
			logger.Fatal("Error getting Traefik routing rules", err)
		}
		logger.Debug(fmt.Sprintf("Rules: %+v", rules))

		logger.Info("Getting hosts from rules")
		hosts := [][]string{}
		for _, r := range rules {
			logger.Debug("Finding hosts in the rule:" + r)
			if match := re.FindAllStringSubmatch(r, -1); len(match) > 0 {
				var t []string
				for _, m := range match {
					logger.Debug(fmt.Sprintf("Found host: %s", m[0]))
					// assume https://
					t = append(t, "https://"+m[0])
				}

				logger.Debug(fmt.Sprintf("Processed all targets on rule and found: %+v", t))
				hosts = append(hosts, t)
			}
		}
		logger.Debug(fmt.Sprintf("All hosts: %+v", hosts))

		logger.Info("Generating Prometheus data")
		tg := []prometheus.TargetGroup{}
		for _, t := range hosts {
			tg = append(tg, prometheus.TargetGroup{Targets: t})
		}

		logger.Debug("Creating config folder if it does not exist")
		err = helpers.InitFolder(conf.TargetsFile)
		if err != nil {
			logger.Fatal("Error creating config folder", err)
		}

		logger.Info("Creating Prometheus target file")
		logger.Debug("Marshalling JSON")
		if file, err := json.MarshalIndent(tg, "", "  "); err != nil {
			logger.Fatal("Error creating JSON data for Prometheus", err)
		} else {
			logger.Debug("Write JSON to file")
			if err = ioutil.WriteFile(conf.TargetsFile, file, 0755); err != nil {
				logger.Fatal("Error writing to JSON file", err)
			}
		}

		logger.Info(fmt.Sprintf("Sleeping %v seconds until next run", conf.RunInterval))
		time.Sleep(time.Duration(conf.RunInterval) * time.Second)
	}
}
