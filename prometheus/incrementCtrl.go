package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// IncrementCtrl updates the kafka offset counters
func (exporter *Exporter) IncrementCtrl(topic string, partition int32, offset int64) {
	labels := prometheus.Labels{
		"topic":     topic,
		"partition": fmt.Sprint(partition),
	}
	kafkaOffsets.With(labels).Add(float64(offset))
}
