package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Exporter provides export features to Prometheus
type Exporter struct {
}

// Initialize Prometheus Exporter, listen on addr with path /metrics
func (exporter *Exporter) Initialize(addr string) {
	// export prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(addr, nil)
	}()
	log.Println("Enabled Prometheus metrics endpoint.")

	// register collectors
	exporter.registerCollectors()
}

func (exporter *Exporter) registerCollectors() {
	// TODO make more dynamic
	prometheus.MustRegister(flowNumber, flowBytes, flowPackets)
}
