package main

import (
	"github.com/czerwonk/bird_exporter/calico"
	"github.com/czerwonk/bird_exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type MetricCollector struct {
	exporters      map[int][]metrics.MetricExporter
	ippExpotersMap map[string][]metrics.MetricExporter
	ippoolExporter metrics.MetricExporter
}

func NewMetricCollector() *MetricCollector {

	l := &metrics.DefaultLabelStrategy{}
	//e := metrics.NewGenericProtocolMetricExporter("IPPool", true, l)
	ipplist, err := calico.Show()
	if err != nil {
		log.Errorf("can not get ipp")
		return nil
	}
	ippoolExporter := metrics.NewGenericIPPoolExporter("IPPool", true, l, ipplist)

	return &MetricCollector{ippoolExporter: ippoolExporter}
}

func (m *MetricCollector) Describe(ch chan<- *prometheus.Desc) {

	//m.uptimeDesc = prometheus.NewDesc(prefix+"_uptime", "Uptime of the protocol in seconds", labels, nil)
	m.ippoolExporter.Describe(ch)
}

func (m *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	m.ippoolExporter.Export(nil, ch)
}
