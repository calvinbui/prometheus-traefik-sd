package helpers

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/calvinbui/prometheus-traefik-sd/internal/logger"
)

const SCHEME = "https://"

var re = regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)

func GetHostsFromRules(rules []string) [][]string {
	hosts := [][]string{}

	for _, r := range rules {
		logger.Debug("Finding hosts in the rule:" + r)
		if match := re.FindAllStringSubmatch(r, -1); len(match) > 0 {
			var t []string
			for _, m := range match {
				logger.Debug(fmt.Sprintf("Found host: %s", m[0]))
				// assume https://
				t = append(t, SCHEME+m[0])
			}

			logger.Debug("Sorting targets in alphabetical order")
			sort.Strings(t)

			logger.Debug(fmt.Sprintf("Processed all targets on rule and found: %+v", t))
			hosts = append(hosts, t)
		}
	}

	return hosts
}
