package main

import (
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/tophost"
)

var validCustomerIDs []int

func runKafkaListener() {
	go handleControlMessages()

	if *filterCustomerIDs != "" {
		stringIDs := strings.Split(*filterCustomerIDs, ",")
		for _, stringID := range stringIDs {
			customerID, err := strconv.Atoi(stringID)
			if err != nil {
				continue
			}
			validCustomerIDs = append(validCustomerIDs, customerID)
		}
		sort.Ints(validCustomerIDs)

		validCustomerIDsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(validCustomerIDs)), ","), "[]")
		log.Printf("Filter flows for customer ids %s\n", validCustomerIDsStr)
	} else {
		log.Printf("No customer filter enabled.\n")
	}

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

	// consider flow when filter applies
	if len(validCustomerIDs) == 0 || isValidCustomerID(int(flow.GetCid())) {
		promExporter.Increment(flow)
		tophostExporter.Consider(tophost.Input{
			Cid:       flow.GetCid(),
			IPSrc:     fmt.Sprintf("%v", net.IP(flow.GetSrcIP())),
			IPDst:     fmt.Sprintf("%v", net.IP(flow.GetDstIP())),
			Peer:      flow.GetPeer(),
			Direction: flow.GetDirection().String(),
			Packets:   flow.GetPackets(),
			Bytes:     flow.GetBytes(),
		})
	}

}

func isValidCustomerID(cid int) bool {
	pos := sort.SearchInts(validCustomerIDs, cid)
	if pos == len(validCustomerIDs) {
		return false
	}
	return validCustomerIDs[pos] == cid
}

func handleControlMessages() {
	ctrlChan := kafkaConn.GetConsumerControlMessages()
	var offsetPerPartition []int64
	for {
		ctrlMsg, ok := <-ctrlChan
		if !ok {
			kafkaConn.CancelConsumerControlMessages()
			return
		}
		partition := ctrlMsg.Partition

		// extend offsetPerPartition array if needed
		if len(offsetPerPartition) < int(partition) {
			n := int(partition) - len(offsetPerPartition) + 1
			newArr := make([]int64, n)
			offsetPerPartition = append(offsetPerPartition, newArr...)
		}

		offsetDiff := ctrlMsg.Offset - offsetPerPartition[partition]
		offsetPerPartition[partition] = ctrlMsg.Offset

		promExporter.IncrementCtrl(*kafkaInTopic, partition, offsetDiff)
	}
}
