package prometheus

import (
	"github.com/bwNetFlow/dashboard/counters"
	"github.com/prometheus/client_golang/prometheus"
)

// implemented according to https://godoc.org/github.com/prometheus/client_golang/prometheus#example-Collector

// NewFacadeCollector returns a new instance of FacadeCollector
func NewFacadeCollector(counter counters.Counter, labels []string) FacadeCollector {
	return FacadeCollector{
		counter: counter,
		labels:  labels,
		desc: prometheus.NewDesc(
			counter.GetName(),
			"consumer_dashboard export of "+counter.GetName(),
			labels,
			nil,
		),
	}
}

// FacadeCollector implements the Collector interface.
type FacadeCollector struct {
	counter counters.Counter
	labels  []string
	desc    *prometheus.Desc
}

// Describe returns the prometheus metric description
func (collector FacadeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

// Collect returns the counter values and labels from the internal counting system
func (collector FacadeCollector) Collect(ch chan<- prometheus.Metric) {
	fields := collector.counter.GetFields()
	for _, item := range fields {
		value := float64(item.Value)
		labels := make([]string, len(item.Label.Fields))
		for i, labelName := range collector.labels {
			label, ok := item.Label.Fields[labelName]
			if ok {
				labels[i] = label
			}
		}

		ch <- prometheus.MustNewConstMetric(
			collector.desc,
			prometheus.CounterValue,
			value,
			labels...,
		)
	}

}
