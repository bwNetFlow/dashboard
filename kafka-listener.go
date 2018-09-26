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
		flow := <-kafkaConn.ConsumerChannel()
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

	// consider flow when filter applies: 0 (all flows) OR customerID matches
	if *filterCustomerID == uint64(0) || uint64(flow.GetCid()) == *filterCustomerID {
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
