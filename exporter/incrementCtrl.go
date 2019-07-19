package exporter

import (
	"fmt"

	"github.com/bwNetFlow/dashboard/counters"
)

// IncrementCtrl updates the kafka offset counters
func (exporter *Exporter) IncrementCtrl(topic string, partition int32, offset int64) {
	labels := counters.NewLabel(map[string]string{
		"topic":     topic,
		"partition": fmt.Sprint(partition),
	})
	counters.KafkaOffsets.Add(labels, uint64(offset))
}
