package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/counters"
)

// IncrementCtrl updates the kafka offset counters
func (exporter *Exporter) IncrementCtrl(topic string, partition int32, offset int64) {
	labels := prometheus.Labels{
		"topic":     topic,
		"partition": fmt.Sprint(partition),
	}
	counters.KafkaOffsets.With(labels).Add(float64(offset))
}
