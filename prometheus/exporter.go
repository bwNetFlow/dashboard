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
func (exporter *Exporter) Initialize(addr string, exporterType string) {
	// export prometheus metrics
	httpServerMux := http.NewServeMux()
	httpServerMux.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(addr, httpServerMux)
	}()
	log.Println("Enabled Prometheus metrics endpoint.")

	// register collectors
	if exporterType == "data" {
		exporter.registerDataCollectors()
	} else if exporterType == "meta" {
		exporter.registerMetaCollectors()
	} else {
		log.Printf("invalid exporterType %s!\n", exporterType)
	}
}

func (exporter *Exporter) registerDataCollectors() {
	// TODO make more dynamic
	prometheus.MustRegister(msgcount, flowNumber, flowBytes, flowPackets, hostBytes, hostConnections)
}

func (exporter *Exporter) registerMetaCollectors() {
	// TODO make more dynamic
	prometheus.MustRegister(kafkaOffsets)
}
