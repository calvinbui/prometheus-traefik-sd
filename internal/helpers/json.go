package helpers

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/calvinbui/blackbox-traefik-sd/internal/logger"
	"github.com/calvinbui/blackbox-traefik-sd/internal/prometheus"
)

func CreateJSON(tgs []prometheus.TargetGroups, folder string) error {
	for _, tg := range tgs {
		logger.Debug("Marshalling JSON")
		if config, err := json.MarshalIndent(tg, "", "  "); err != nil {
			logger.Fatal("Error creating JSON data for Prometheus", err)
		} else {
			filePath := createFileName(folder, tg[0].Targets)
			logger.Debug("Write JSON to file " + filePath)
			if err = ioutil.WriteFile(filePath, config, 0755); err != nil {
				logger.Fatal("Error writing JSON to file "+filePath, err)
			}
		}
	}

	return nil
}

func createFileName(folder string, targets []string) string {
	for i := range targets {
		targets[i] = strings.TrimPrefix(targets[i], SCHEME)
	}

	return path.Join(folder, strings.Join(targets, "_"))
}
