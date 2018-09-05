package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// TopHost updates one entry for Top Hosts
func (exporter *Exporter) TopHost(ip string, bytes uint64) {
	labels := prometheus.Labels{
		"ip": ip,
	}
	hostBytes.With(labels).Add(float64(bytes))
}
