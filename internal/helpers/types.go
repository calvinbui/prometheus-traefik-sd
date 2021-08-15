package helpers

import "github.com/calvinbui/prometheus-traefik-sd/internal/prometheus"

type PromTargetFile struct {
	FilePath string
	Data     []prometheus.TargetGroups
}
