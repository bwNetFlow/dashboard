package counters

// Meta Monitoring Data, to be added to default /metrics
var (
	Msgcount     = NewCounter("kafka_messages_total")
	KafkaOffsets = NewCounter("kafka_offset_current")
)

// Flow Data, to be exported to Prometheus' /flowdata or to InfluxDB
var (
	// general counters
	FlowNumber  = NewCounter("flow_number_total")
	FlowBytes   = NewCounter("flow_bytes")
	FlowPackets = NewCounter("flow_packets")

	// TOP HOSTS
	HostBytes       = NewCounter("host_bytes")
	HostConnections = NewCounter("host_connections_total")
)
