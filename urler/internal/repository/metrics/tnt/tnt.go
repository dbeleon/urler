package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CollisionsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "urler",
		Subsystem: "tnt_urls",
		Name:      "collisions_total",
	},
		[]string{},
	)
)
