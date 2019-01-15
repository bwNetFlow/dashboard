package counters

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var ( // Meta Monitoring Data, to be added to default /metrics
	Msgcount = Counter{
		Fields: make(map[uint32]CounterItems),
		Name:   "kafka_messages_total",
		Access: &sync.Mutex{},
	}

	Kafkalabels = []string{
		"topic",
		"partition",
	}
	KafkaOffsets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_offset_current",
			Help: "Current Kafka Offset of the consumer",
		}, Kafkalabels)
)

var ( // Flow Data, to be exported on /flowdata
	FlowNumber = Counter{
		Fields: make(map[uint32]CounterItems),
		Name:   "flow_number_total",
		Access: &sync.Mutex{},
	}
	FlowBytes = Counter{
		Fields: make(map[uint32]CounterItems),
		Name:   "flow_bytes",
		Access: &sync.Mutex{},
	}
	FlowPackets = Counter{
		Fields: make(map[uint32]CounterItems),
		Name:   "flow_packets",
		Access: &sync.Mutex{},
	}

	// TOP HOSTS
	Hostlabels = []string{
		"cid",
		"ipSrc",
		"ipDst",
		"peer",
	}

	HostBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "host_bytes",
			Help: "Number of Bytes for top N hosts.",
		}, Hostlabels)

	HostConnections = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "host_connections_total",
			Help: "Number of Connections for top N hosts.",
		}, Hostlabels)
)
