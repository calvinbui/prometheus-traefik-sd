package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/calvinbui/blackbox-traefik-sd/internal/prometheus"
)

func GetTargetsFromRules(rules []string) [][]string {
	var re = regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	var targets [][]string

	for _, r := range rules {
		// get all hosts for the route
		match := re.FindAllStringSubmatch(r, -1)

		var t []string
		for _, m := range match {
			t = append(t, m[0])
		}

		targets = append(targets, t)
	}

	return targets
}

func BuildPrometheusTargetFile(targets [][]string, filePath string) error {
	var tg []prometheus.TargetGroup

	for _, t := range targets {
		tg = append(tg, prometheus.TargetGroup{
			Targets: t,
		})
	}

	err := initFolder(filePath)
	if err != nil {
		return err
	}

	err = initConfig(filePath)
	if err != nil {
		return err
	}

	file, err := json.MarshalIndent(tg, "", "")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, file, 0755)
	if err != nil {
		return err
	}

	return nil
}

func initFolder(filePath string) error {
	folder := path.Dir(filePath)

	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		err := os.Mkdir(path.Dir(filePath), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func initConfig(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
