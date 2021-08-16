package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/calvinbui/prometheus-traefik-sd/internal/logger"
	"github.com/calvinbui/prometheus-traefik-sd/internal/prometheus"
)

func CreateJSON(tgs []PromTargetFile, folder string) error {
	for _, tg := range tgs {
		logger.Trace(fmt.Sprintf("Marshalling JSON for %+v", tg.Data))
		config, err := json.MarshalIndent(tg.Data, "", "  ")
		if err != nil {
			logger.Fatal("Error creating JSON data for Prometheus", err)
		}

		logger.Debug("Checking if " + tg.FilePath + " exists")
		if _, err := os.Stat(tg.FilePath); os.IsNotExist(err) {
			logger.Debug(tg.FilePath + " does not exists")
			logger.Info(fmt.Sprintf("Creating target file %s for %+v", tg.FilePath, tg.Data[0].Targets))
			if err = writeFile(tg.FilePath, config); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			logger.Debug(tg.FilePath + " exists")
			if currentData, err := ioutil.ReadFile(tg.FilePath); err != nil {
				return err
			} else {
				var currentConfig prometheus.TargetGroups
				if err := json.Unmarshal(currentData, &currentConfig); err != nil {
					return err
				} else {
					if reflect.DeepEqual(currentConfig, tg.Data) {
						logger.Debug(fmt.Sprintf("The JSON file %s exists, no actions will be performed", tg.FilePath))
					} else {
						logger.Info(tg.FilePath + " is outdated. Updating.")
						if err = writeFile(tg.FilePath, config); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func CreateFileName(folder string, targets []string) string {
	var noSchemeTargets []string
	for i := range targets {
		noSchemeTargets = append(noSchemeTargets, strings.TrimPrefix(targets[i], SCHEME))
	}

	return path.Join(folder, strings.Join(noSchemeTargets, "_")+".json")
}

func writeFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0755)
}
