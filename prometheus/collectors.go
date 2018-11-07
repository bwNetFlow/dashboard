package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	labels = []string{
		// "src_port",
		// "dst_port",
		"ipversion",
		"application",
		"protoname",
		"direction",
		"cid",
		"peer",
		"remotecountry",
	}

	msgcount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "msgcount",
			Help: "Number Kafka messages.",
		}, labels)

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

// TOP HOSTS
var (
	hostlabels = []string{
		"cid",
		"ipSrc",
		"ipDst",
		"peer",
	}

	hostBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "host_bytes",
			Help: "Number of Bytes for top N hosts.",
		}, hostlabels)

	hostConnections = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "host_connections",
			Help: "Number of Connections for top N hosts.",
		}, hostlabels)
)

// KAFKA METRICS
var (
	kafkalabels = []string{
		"topic",
		"partition",
	}
	kafkaOffsets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_offset",
			Help: "Kafka Offset of the consumer",
		}, kafkalabels)
)
