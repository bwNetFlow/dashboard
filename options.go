package main

import "flag"

var (
	// common options
	logFile = flag.String("log", "./consumer_dashboard.log", "Location of the log file.")

	// Kafka options
	kafkaConsumerGroup = flag.String("kafka.consumer_group", "dashboard", "Kafka Consumer Group")
	kafkaInTopic       = flag.String("kafka.topic", "flow-messages-enriched", "Kafka topic to consume from")
	kafkaBroker        = flag.String("kafka.brokers", "127.0.0.1:9092,[::1]:9092", "Kafka brokers separated by commas")
	kafkaUser          = flag.String("kafka.user", "", "Kafka username to authenticate with")
	kafkaPass          = flag.String("kafka.pass", "", "Kafka password to authenticate with")
	kafkaAuth          = flag.Bool("kafka.auth", true, "Kafka auth enable/disable")
	kafkaTLS           = flag.Bool("kafka.tls", true, "Kafka tls connection enable/disable")

	// prometheus options
	exportPrometheus     = flag.Bool("export.prometheus", false, "enable prometheus export endpoint")
	exportPrometheusBind = flag.String("export.prometheusBind", ":8080", "Host/Port to bind prometheus exporter")

	// influx options
	exportInflux           = flag.Bool("export.influx", false, "enable influxdb push (requires further endpoint params)")
	exportInfluxURL        = flag.String("export.influxUrl", "", "Path to Influx DB")
	exportInfluxUser       = flag.String("export.influxUser", "", "Username for Influx")
	exportInfluxPass       = flag.String("export.influxPass", "", "Password for Influx")
	exportInfluxDatabase   = flag.String("export.influxDatabase", "", "Database to use for Influx (Prefix in case of export.influxPerCid)")
	exportInfluxExportFreq = flag.Int("export.influxFreq", 10, "Frequency [seconds] for exports")
	exportInfluxPerCid     = flag.Bool("export.influxPerCid", false, "Export in separate customer database named export.influxDatabase-$CID")

	// dashboard consumer specific
	filterCustomerIDs = flag.String("customerid", "", "If defined, only flows for this customer are considered. Leave empty to disable filter. Provide comma separated list to filter for multiple customers.")
)
