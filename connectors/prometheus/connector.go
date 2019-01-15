package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Connector provides export features to Prometheus
type Connector struct {
	Addr string
}

// Initialize Prometheus Exporter, listen on addr with path /metrics and /flowdata
func (connector *Connector) Initialize() {
	//prometheus.MustRegister(counters.Msgcount, counters.KafkaOffsets)

	flowReg := prometheus.NewRegistry()
	//flowReg.MustRegister(counters.FlowNumber, counters.FlowBytes, counters.FlowPackets, counters.HostBytes, counters.HostConnections)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/flowdata", promhttp.HandlerFor(flowReg, promhttp.HandlerOpts{}))

	go func() {
		http.ListenAndServe(connector.Addr, nil)
	}()
	log.Println("Enabled Prometheus /metrics and /flowdata endpoints.")
}
