package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/prometheus"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/tophost"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/util"
	kafka "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/kafkaconnector"

	"github.com/Shopify/sarama"
)

// KafkaConn holds the global kafka connection
var kafkaConn = kafka.Connector{}
var promExporter = prometheus.Exporter{}
var tophostExporter = tophost.Exporter{}

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

	// Enable TopHost Counter
	var maxHosts = 10
	var exportInterval = 15 * time.Second
	tophostExporter.Initialize(promExporter, maxHosts, exportInterval)

	// Set kafka auth
	if *kafkaUser != "" {
		kafkaConn.SetAuth(*kafkaUser, *kafkaPass)
	} else {
		kafkaConn.SetAuthAnon()
	}

	// Establish Kafka Connection
	kafkaConn.StartConsumer(*kafkaBroker, []string{*kafkaInTopic}, *kafkaConsumerGroup, sarama.OffsetNewest)
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
