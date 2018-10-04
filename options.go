package main

import "flag"

var (
	// common options
	logFile = flag.String("log", "", "Location of the log file.")
	pidFile = flag.String("pid", "", "Location of the pid file.")

	// Kafka options
	kafkaConsumerGroup = flag.String("kafka.consumer_group", "dashboard", "Kafka Consumer Group")
	kafkaInTopic       = flag.String("kafka.topic", "flow-messages-enriched", "Kafka topic to consume from")
	kafkaBroker        = flag.String("kafka.brokers", "127.0.0.1:9092,[::1]:9092", "Kafka brokers separated by commas")
	kafkaUser          = flag.String("kafka.user", "", "Kafka username to authenticate with")
	kafkaPass          = flag.String("kafka.pass", "", "Kafka password to authenticate with")

	// prometheus options
	// TODO listen on addr

	// dashboard consumer specific
	filterCustomerID = flag.Uint64("customerid", 0, "If defined, only flows for this customer are considered. \"0\" to disable filter.")
)
