package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/calvinbui/prometheus-traefik-sd/internal/logger"
)

func CreateJSON(tgs []PromTargetFile, folder string) error {
	for _, tg := range tgs {
		if _, err := os.Stat(tg.FilePath); os.IsNotExist(err) {
			logger.Debug("Marshalling JSON to " + tg.FilePath)
			if config, err := json.MarshalIndent(tg.Data, "", "  "); err != nil {
				logger.Fatal("Error creating JSON data for Prometheus", err)
			} else {
				logger.Debug(fmt.Sprintf("Write to JSON file %s: %s", tg.FilePath, config))
				if err = ioutil.WriteFile(tg.FilePath, config, 0755); err != nil {
					logger.Fatal("Error writing JSON to file "+tg.FilePath, err)
				}
			}
		} else if err != nil {
			return err
		} else {
			logger.Debug(fmt.Sprintf("The JSON file %s exists, no actions will be performed", tg.FilePath))
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
