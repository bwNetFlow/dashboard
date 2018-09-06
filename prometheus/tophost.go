package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// TopHostTraffic updates one entry for Top Hosts
func (exporter *Exporter) TopHostTraffic(ip string, bytes uint64) {
	labels := prometheus.Labels{
		"ip": ip,
	}
	hostBytes.With(labels).Add(float64(bytes))
}

// TopHostConnections updates one entry for Top Hosts
func (exporter *Exporter) TopHostConnections(ip string, connections uint64) {
	labels := prometheus.Labels{
		"ip": ip,
	}
	hostConnections.With(labels).Add(float64(connections))
}
