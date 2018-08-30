package main

import (
	"flag"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	proto "github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sarama "gopkg.in/Shopify/sarama.v1"
	"io"
	"io/ioutil"
	flow "kafka-processor/api/flow-message-enriched"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	LogFile            = flag.String("log", "./dashboard.log", "Location of the log file.")
	KafkaConsumerGroup = flag.String("kafka.consumer_group", "dashboard", "Kafka Consumer Group")
	KafkaInTopic       = flag.String("kafka.in_topic", "flow-messages-enriched", "Kafka topic to consume from")
	KafkaBroker        = flag.String("kafka.brokers", "127.0.0.1:9092,[::1]:9092", "Kafka brokers separated by commas")
	flowNumber         = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_number",
			Help: "Number of Flows received.",
		}, []string{"src_port", "dst_port", "src_cc", "dst_cc", "proto", "afi", "src_if", "dst_if"})
	flowBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_bytes",
			Help: "Number of Bytes received across Flows.",
		}, []string{"src_port", "dst_port", "src_cc", "dst_cc", "proto", "afi", "src_if", "dst_if"})
	flowPackets = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flow_packets",
			Help: "Number of Packets received across Flows.",
		}, []string{"src_port", "dst_port", "src_cc", "dst_cc", "proto", "afi", "src_if", "dst_if"})
)

func main() {
	flag.Parse()
	var err error

	// All initialization below is done globally, as there is no need for
	// each partition worker routine to do it for themselves.

	// initialize logger
	logfile, err := os.OpenFile(*LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		println("Error opening file for logging: %v", err)
		return
	}
	defer logfile.Close()
	mw := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mw)

	log.Println("-------------------------- Started.")

	// write a pid file too avoid searching the processes:
	// 'kill -SIGUSR1 $(cat analysis.pid)'
	log.Printf("Writing PID %d to file '%s'.", os.Getpid(), "./consumer_dashboard.pid")
	err = ioutil.WriteFile("./consumer_dashboard.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		log.Fatal("Could not write PID File.")
	}

	// export prometheus metrics
	prometheus.MustRegister(flowNumber)
	prometheus.MustRegister(flowBytes)
	prometheus.MustRegister(flowPackets)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	log.Println("Enabled Prometheus metrics endpoint.")

	// init consumer and connect
	brokers := strings.Split(*KafkaBroker, ",")
	consConf := cluster.NewConfig()
	consConf.Consumer.Return.Errors = true
	consConf.Consumer.Offsets.Initial = sarama.OffsetOldest
	consConf.Group.Return.Notifications = true
	topics := []string{*KafkaInTopic}
	consumer, err := cluster.NewConsumer(brokers, *KafkaConsumerGroup, topics, consConf)
	if err != nil {
		log.Fatalf("Error initializing the Kafka Consumer: %v", err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Kafka Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Kafka Notification: %+v\n", ntf)
		}
	}()

	for {
		flowmsg := &flow.FlowMessage{}

		select {
		case msg, ok := <-consumer.Messages():
			if !ok {
				// TODO: determine when this happens
				log.Println("Message channel closed.")
			}
			consumer.MarkOffset(msg, "") // mark message as processed
			err = proto.Unmarshal(msg.Value, flowmsg)
			if err != nil {
				log.Printf("Received broken message. Unmarshalling error: %v", err)
				continue
			}
		case sig := <-signals:
			switch sig {
			case syscall.SIGUSR1: // do custom stuff
				continue
			case syscall.SIGINT, syscall.SIGTERM:
				os.Remove("./consumer_dashboard.pid")
				log.Print("Received exit signal, kthxbye.")
				os.Exit(0)
			}
		}

		// TODO: do the below stuff correctly
		flowNumber.With(prometheus.Labels{"src_if": string(flowmsg.SrcIf), "afi": "4"}).Inc()
	}
}
