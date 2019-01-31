package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/connectors"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/exporter"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/tophost"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/util"
	kafka "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/kafkaconnector"

	"github.com/Shopify/sarama"
)

// KafkaConn holds the global kafka connection
var kafkaConn = kafka.Connector{}
var mainExporter = exporter.Exporter{}
var metaExporter = exporter.Exporter{}
var tophostExporter = tophost.TophostExporter{}

func main() {

	flag.Parse()
	util.InitLogger(*logFile)

	// catch termination signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signals
		log.Println("Received exit signal, kthxbye.")
		os.Exit(0)
	}()

	// Initialize connectors to TSDBs
	if *exportPrometheus {
		prometheusConnector := connectors.NewPrometheusConnector(":8383")
		prometheusConnector.Initialize()
	}
	if *exportInflux {
		influxConnector := connectors.NewInfluxConnector(*exportInfluxURL, *exportInfluxUser, *exportInfluxPass, *exportInfluxDatabase, *exportInfluxExportFreq, *exportInfluxPerCid)
		influxConnector.Initialize()
		defer influxConnector.Close()
	}

	// Enable TopHost Counter
	var maxHosts = 10
	var exportInterval = 15 * time.Second
	var hostMaxAge = -20 * time.Minute // 20 minutes old
	tophostExporter.Initialize(mainExporter, maxHosts, exportInterval, hostMaxAge)

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
