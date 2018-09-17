package main

import (
	"fmt"
	"log"
	"net"
	"runtime/debug"

	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/tophost"
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
		tophostExporter.Consider(tophost.Input{
			IPSrc:     fmt.Sprintf("%v", net.IP(flow.GetSrcIP())),
			IPDst:     fmt.Sprintf("%v", net.IP(flow.GetDstIP())),
			Peer:      flow.GetPeer(),
			Direction: flow.GetDirection().String(),
			Packets:   flow.GetPackets(),
			Bytes:     flow.GetBytes(),
		})
	}

}
