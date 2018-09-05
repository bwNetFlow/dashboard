package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	labels = []string{
		"src_port",
		"dst_port",
		"proto",
		"direction",
		"cid",
		"peer",
	}

	flowNumber = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_number",
			Help: "Number of Flows received.",
		}, labels)
	flowBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_bytes",
			Help: "Number of Bytes received across Flows.",
		}, labels)
	flowPackets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_packets",
			Help: "Number of Packets received across Flows.",
		}, labels)
)
