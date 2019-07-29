package counters

import "github.com/bwNetFlow/dashboard/cache"

// Meta Monitoring Data, to be added to default /metrics
var (
	Msgcount     Counter
	KafkaOffsets Counter
)

// Flow Data, to be exported to Prometheus' /flowdata or to InfluxDB
var (
	// general counters
	FlowNumber  Counter
	FlowBytes   Counter
	FlowPackets Counter

	// TOP HOSTS
	HostBytes       Counter
	HostConnections Counter
)

func InitializeCounters(cache *cache.Cache) {
	// meta monitoring
	Msgcount = NewLocalCounter("kafka_messages_total")
	KafkaOffsets = NewLocalCounter("kafka_offset_current")

	// general counters
	FlowNumber = NewRemoteCounter("flow_number_total", cache)
	FlowBytes = NewRemoteCounter("flow_bytes", cache)
	FlowPackets = NewRemoteCounter("flow_packets", cache)

	// TOP HOSTS
	HostBytes = NewLocalCounter("host_bytes")
	HostConnections = NewLocalCounter("host_connections_total")
}
