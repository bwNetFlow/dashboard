package prometheus

import (
	"log"
	"net/http"

	"github.com/bwNetFlow/dashboard/counters"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Connector provides export features to Prometheus
type Connector struct {
	Addr string
}

// Initialize Prometheus Exporter, listen on addr with path /metrics and /flowdata
func (connector *Connector) Initialize() {

	// default flow counters
	flowLabels := []string{
		"ipversion",
		"application",
		"protoname",
		"direction",
		"cid",
		"peer",
	}
	flowNumber := NewFacadeCollector(counters.FlowNumber, flowLabels)
	flowBytes := NewFacadeCollector(counters.FlowBytes, flowLabels)
	flowPackets := NewFacadeCollector(counters.FlowPackets, flowLabels)

	// top host flow counters
	flowHostLabels := []string{
		"cid",
		"ipSrc",
		"ipDst",
		"peer",
	}
	hostBytes := NewFacadeCollector(counters.HostBytes, flowHostLabels)
	hostConnections := NewFacadeCollector(counters.HostConnections, flowHostLabels)

	flowReg := prometheus.NewPedanticRegistry()
	flowReg.MustRegister(flowNumber, flowBytes, flowPackets, hostBytes, hostConnections)

	// kafka meta counters
	metaLabels := []string{
		"topic",
		"partition",
	}
	msgcount := NewFacadeCollector(counters.Msgcount, []string{})
	KafkaOffsets := NewFacadeCollector(counters.KafkaOffsets, metaLabels)

	metaReg := prometheus.NewPedanticRegistry()
	metaReg.MustRegister(msgcount, KafkaOffsets)

	// add default counters
	metaReg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	// register paths and start http server
	http.Handle("/metrics", promhttp.HandlerFor(metaReg, promhttp.HandlerOpts{}))
	http.Handle("/flowdata", promhttp.HandlerFor(flowReg, promhttp.HandlerOpts{}))
	go func() {
		http.ListenAndServe(connector.Addr, nil)
	}()
	log.Println("Enabled Prometheus /metrics and /flowdata endpoints.")
}
