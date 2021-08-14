package helpers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/calvinbui/blackbox-traefik-sd/internal/logger"
	"github.com/calvinbui/blackbox-traefik-sd/internal/prometheus"
)

func CreateTargetsJSON(tg []prometheus.TargetGroup, filePath string) error {
	logger.Debug("Marshalling JSON")
	if file, err := json.MarshalIndent(tg, "", "  "); err != nil {
		logger.Fatal("Error creating JSON data for Prometheus", err)
	} else {
		logger.Debug("Write JSON to file")
		if err = ioutil.WriteFile(filePath, file, 0755); err != nil {
			logger.Fatal("Error writing to JSON file", err)
		}
	}

	return nil
}
