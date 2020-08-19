package metrics

import (
	"github.com/czerwonk/bird_exporter/calico"
)

type DefaultLabelStrategy struct {
}

func (*DefaultLabelStrategy) LabelNames() []string {
	return []string{"name", "cidr", "selector", "ipip"}
}

func (*DefaultLabelStrategy) LabelValues(p *calico.IPPResult) []string {
	return []string{p.Name, p.CIDR, p.Selector, p.IPIP}
}
