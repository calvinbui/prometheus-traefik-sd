package helpers

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/calvinbui/prometheus-traefik-sd/internal/logger"
)

func CreateJSON(tgs []PromTargetFile, folder string) error {
	for _, tg := range tgs {
		logger.Debug("Marshalling JSON")
		if config, err := json.MarshalIndent(tg.Data, "", "  "); err != nil {
			logger.Fatal("Error creating JSON data for Prometheus", err)
		} else {
			logger.Debug("Write JSON to file " + tg.FilePath)
			if err = ioutil.WriteFile(tg.FilePath, config, 0755); err != nil {
				logger.Fatal("Error writing JSON to file "+tg.FilePath, err)
			}
		}
	}

	return nil
}

func CreateFileName(folder string, targets []string) string {
	for i := range targets {
		targets[i] = strings.TrimPrefix(targets[i], SCHEME)
	}

	return path.Join(folder, strings.Join(targets, "_")+".json")
}
