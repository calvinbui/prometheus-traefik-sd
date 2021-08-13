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

	logger.Debug("Building Traefik rules regex")
	var re = regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)

	for {
		logger.Info("Getting Traefik routing rules")
		rules, err := cTraefik.GetRules()
		if err != nil {
			logger.Fatal("Error getting Traefik routes", err)
		}
		logger.Debug(fmt.Sprintf("Rules: %+v", rules))

		targets := [][]string{}
		for _, r := range rules {
			logger.Debug(fmt.Sprintf("Working on rule: %s", r))
			// get all hosts for the route
			match := re.FindAllStringSubmatch(r, -1)
			if len(match) == 0 {
				continue
			}

			var t []string
			for _, m := range match {
				logger.Debug(fmt.Sprintf("Found host: %s", m[0]))
				// assume https://
				t = append(t, "https://"+m[0])
			}

			logger.Debug(fmt.Sprintf("Processed all targets on rule and found: %+v", t))
			targets = append(targets, t)
		}
		logger.Debug(fmt.Sprintf("Targets: %+v", targets))

		logger.Debug("Generating Prometheus data")
		tg := []prometheus.TargetGroup{}
		for _, t := range targets {
			tg = append(tg, prometheus.TargetGroup{
				Targets: t,
			})
		}

		logger.Debug("Creating config folder if it does not exist")
		err = helpers.InitFolder(conf.TargetsFile)
		if err != nil {
			logger.Fatal("", err)
		}

		logger.Info("Creating Prometheus target file")
		logger.Debug("Marshalling JSON")
		file, err := json.MarshalIndent(tg, "", "  ")
		if err != nil {
			logger.Fatal("", err)
		}

		logger.Debug("Write JSON to file")
		err = ioutil.WriteFile(conf.TargetsFile, file, 0755)
		if err != nil {
			logger.Fatal("", err)
		}

		logger.Info(fmt.Sprintf("Sleeping %v seconds until next run", conf.RunInterval))
		time.Sleep(time.Duration(conf.RunInterval) * time.Second)
	}
}
