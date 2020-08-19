package metrics

import (
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricExporter interface {
	Describe(ch chan<- *prometheus.Desc)
	Export(p *protocol.IPPool, ch chan<- prometheus.Metric)
}
