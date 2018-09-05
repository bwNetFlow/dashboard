package main

import (
	"log"
	"net"

	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
)

func runKafkaListener() {
	// handle kafka errors
	go func() {
		for err := range kafkaConn.Errors() {
			log.Printf("Kafka Error: %s\n", err.Error())
		}
	}()

	// handle kafka notifications
	go func() {
		for ntf := range kafkaConn.Notifications() {
			log.Printf("Kafka Notification: %+v\n", ntf)
		}
	}()

	startPeriodicHostExport()

	// handle kafka flow messages in foreground
	for {
		flow := <-kafkaConn.Messages()
		handleFlow(flow)
	}
}

func handleFlow(flow *flow.FlowMessage) {

	if uint64(flow.GetCid()) == *filterCustomerID {
		promExporter.Increment(flow)
		ipDst := net.IP(flow.GetDstIP()).String()
		ipSrc := net.IP(flow.GetSrcIP()).String()
		countHostTraffic(ipDst, flow.GetBytes())
		countHostTraffic(ipSrc, flow.GetBytes())
	}

}
