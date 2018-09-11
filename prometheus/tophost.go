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

// RemoveTopHostTraffic removes the host from the counter vector
func (exporter *Exporter) RemoveTopHostTraffic(ip string) {
	labels := prometheus.Labels{
		"ip": ip,
	}
	hostBytes.Delete(labels)
}

// TopHostConnections updates one entry for Top Hosts
func (exporter *Exporter) TopHostConnections(ip string, connections uint64) {
	labels := prometheus.Labels{
		"ip": ip,
	}
	hostConnections.With(labels).Add(float64(connections))
}

// RemoveTopHostConnections removes the host from the counter vector
func (exporter *Exporter) RemoveTopHostConnections(ip string) {
	labels := prometheus.Labels{
		"ip": ip,
	}
	hostConnections.Delete(labels)
}
