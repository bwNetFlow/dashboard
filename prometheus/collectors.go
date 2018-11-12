package prometheus

import "github.com/prometheus/client_golang/prometheus"

var ( // Meta Monitoring Data, to be added to default /metrics
	msgcount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "msgcount",
			Help: "Number of Kafka messages",
		})

	kafkalabels = []string{
		"topic",
		"partition",
	}
	kafkaOffsets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_offset",
			Help: "Current Kafka Offset of the consumer",
		}, kafkalabels)
)

var ( // Flow Data, to be exported on /flowdata
	labels = []string{
		// "src_port",
		// "dst_port",
		"ipversion",
		"application",
		"protoname",
		"direction",
		"cid",
		"peer",
		// "remotecountry",
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

	// TOP HOSTS
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
