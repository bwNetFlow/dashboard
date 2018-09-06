package main

import (
	"log"
	"net"
	"runtime/debug"

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

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered panic in handleFlow", r)
			debug.PrintStack()
			log.Printf("failed flow: %+v\n", flow)
		}
	}()

	if uint64(flow.GetCid()) == *filterCustomerID {
		promExporter.Increment(flow)
		ipDst := net.IP(flow.GetDstIP())
		ipDstStr := ipDst.String()
		ipSrc := net.IP(flow.GetSrcIP())
		ipSrcStr := ipSrc.String()
		countHostTraffic(ipDstStr, flow.GetBytes())
		countHostTraffic(ipSrcStr, flow.GetBytes())
	}

}
