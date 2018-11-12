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

// Initialize Prometheus Exporter, listen on addr with path /metrics and /flowdata
func (exporter *Exporter) Initialize(addr string) {
	prometheus.MustRegister(msgcount, kafkaOffsets)

	flowReg := prometheus.NewRegistry()
	flowReg.MustRegister(flowNumber, flowBytes, flowPackets, hostBytes, hostConnections)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/flowdata", promhttp.HandlerFor(flowReg, promhttp.HandlerOpts{}))

	go func() {
		http.ListenAndServe(addr, nil)
	}()
	log.Println("Enabled Prometheus /metrics and /flowdata endpoints.")
}
