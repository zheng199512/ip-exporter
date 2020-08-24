package metrics

import (
	"github.com/czerwonk/bird_exporter/calico"
	"github.com/czerwonk/bird_exporter/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type IPPoolExporter struct {
	labelStrategy LabelStrategy
	IPPResultList []calico.IPPResult
	ippoolDesc    *prometheus.Desc
}

func NewGenericIPPoolExporter(prefix string, newNaming bool, labelStrategy LabelStrategy, IPPResultList []calico.IPPResult) *IPPoolExporter {
	m := &IPPoolExporter{labelStrategy: labelStrategy, IPPResultList: IPPResultList}
	return m
}

func (m *IPPoolExporter) Describe(ch chan<- *prometheus.Desc) {

}

func (m *IPPoolExporter) Export(p *protocol.IPPool, ch chan<- prometheus.Metric) {

	for _, ipp := range m.IPPResultList {
		labels := m.labelStrategy.LabelNames()

		ippTotal := prometheus.NewDesc("calico_ippool"+"_total", "total of ippool", labels, nil)
		ippInuse := prometheus.NewDesc("calico_ippool"+"_inuse", "inuse of ippool", labels, nil)
		l := m.labelStrategy.LabelValues(&ipp)
		ch <- prometheus.MustNewConstMetric(ippTotal, prometheus.GaugeValue, ipp.Total, l...)
		ch <- prometheus.MustNewConstMetric(ippInuse, prometheus.GaugeValue, ipp.Inuse, l...)
	}
}
