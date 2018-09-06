package main

import (
	"fmt"
	"log"
	"net"
	"runtime/debug"

	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
)

func runKafkaListener() {
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

		ipSrc := net.IP(flow.GetSrcIP())
		ipSrcStr := fmt.Sprintf("%v", ipSrc)
		ipDst := net.IP(flow.GetDstIP())
		ipDstStr := fmt.Sprintf("%v", ipDst)
		tophostExporter.Consider(ipSrcStr, ipDstStr, flow.GetBytes())
	}

}
