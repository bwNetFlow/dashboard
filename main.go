package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	kafka "kafka/consumer_dashboard/kafka"
	"kafka/consumer_dashboard/prometheus"
	"kafka/consumer_dashboard/util"

	"github.com/Shopify/sarama"
)

// KafkaConn holds the global kafka connection
var kafkaConn = kafka.Connector{}
var promExporter = prometheus.Exporter{}

func main() {

	flag.Parse()
	util.InitLogger(*logFile)
	util.WritePid(*pidFile)

	// catch termination signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(signals)
	go func() {
		<-signals
		shutdown(0)
	}()

	// Enable Prometheus Export
	promExporter.Initialize(":8080")

	// Establish Kafka Connection
	kafkaConn.Connect(*kafkaBroker, *kafkaInTopic, *kafkaConsumerGroup, sarama.OffsetNewest)
	defer kafkaConn.Close()
	runKafkaListener()
}

func shutdown(exitcode int) {
	kafkaConn.Close()
	os.Remove("./consumer_dashboard.pid")
	log.Println("Received exit signal, kthxbye.")
	// return exit code
	os.Exit(exitcode)
}
