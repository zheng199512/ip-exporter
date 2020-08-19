package metrics

import "github.com/czerwonk/bird_exporter/calico"

// LabelStrategy abstracts the label generation for protocol metrics
type LabelStrategy interface {
	// LabelNames is the list of label names
	LabelNames() []string

	// Label values is the list of values for the labels specified in `LabelNames()`
	LabelValues(re *calico.IPPResult) []string
}
