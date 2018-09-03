package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	flowNumber = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_number",
			Help: "Number of Flows received.",
		}, []string{"src_port", "dst_port", "src_cc", "dst_cc", "proto", "afi", "src_if", "dst_if"})
	flowBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_bytes",
			Help: "Number of Bytes received across Flows.",
		}, []string{"src_port", "dst_port", "src_cc", "dst_cc", "proto", "afi", "src_if", "dst_if"})
	flowPackets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_packets",
			Help: "Number of Packets received across Flows.",
		}, []string{"src_port", "dst_port", "src_cc", "dst_cc", "proto", "afi", "src_if", "dst_if"})
)
