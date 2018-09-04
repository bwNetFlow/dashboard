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

	// handle kafka flow messages in foreground
	for {
		flow := <-kafkaConn.Messages()
		handleFlow(flow)
	}
}

func handleFlow(flow *flow.FlowMessage) {

	if uint64(flow.GetCid()) == *filterCustomerID {
		promExporter.Increment(flow)
	}

	/*
		srcIP := decodeIP("", flow.GetSrcIP())
		dstIP := decodeIP("", flow.GetDstIP())
		if srcIP == "134.60.30.246" || dstIP == "134.60.30.246" {
			// fmt.Printf("flow: %v -> %v, dir: %v, cid: %v, norm: %v\n", net.IP(flow.SrcIP), net.IP(flow.DstIP), flow.Direction, flow.Cid, flow.Normalized)
			// fmt.Printf("      %v - %v (%v) -> %v - %v (%v)\n", flow.SrcIfName, flow.SrcIfDesc, flow.SrcIfSpeed, flow.DstIfName, flow.DstIfDesc, flow.DstIfSpeed)
			promExporter.Increment(flow)
		}
	*/

}

func decodeIP(name string, value []byte) string {
	str := ""
	ipconv := net.IP{}
	if value != nil {
		invvalue := make([]byte, len(value))
		for i := range value {
			invvalue[len(value)-i-1] = value[i]
		}
		ipconv = value
		str += name + ipconv.String()
	}
	return str
}
