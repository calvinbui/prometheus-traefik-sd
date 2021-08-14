package helpers

import (
	"fmt"
	"regexp"

	"github.com/calvinbui/blackbox-traefik-sd/internal/logger"
)

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
				t = append(t, "https://"+m[0])
			}

			logger.Debug(fmt.Sprintf("Processed all targets on rule and found: %+v", t))
			hosts = append(hosts, t)
		}
	}

	logger.Debug(fmt.Sprintf("All hosts: %+v", hosts))

	return hosts
}
