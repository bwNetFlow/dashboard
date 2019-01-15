package exporter

import (
	"fmt"

	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/counters"
)

// Increment updates the counters by one flow
func (exporter *Exporter) Increment(flow *flow.FlowMessage) {
	application := ""
	_, appGuess1 := filterPopularPorts(flow.GetSrcPort())
	_, appGuess2 := filterPopularPorts(flow.GetDstPort())
	if appGuess1 != "" {
		application = appGuess1
	} else if appGuess2 != "" {
		application = appGuess2
	}

	labels := counters.Label{
		Fields: map[string]string{
			// "src_port":      fmt.Sprint(srcPort),
			// "dst_port":      fmt.Sprint(dstPort),
			"ipversion":   flow.GetIPversion().String(),
			"application": application,
			"protoname":   fmt.Sprint(flow.GetProtoName()),
			"direction":   fmt.Sprint(flow.GetDirection()),
			"cid":         fmt.Sprint(flow.GetCid()),
			"peer":        flow.GetPeer(),
			// "remotecountry": flow.GetRemoteCountry(),
		},
	}

	counters.Msgcount.Add(counters.NewEmptyLabel(), 1)

	counters.FlowNumber.Add(labels, flow.GetSamplingRate())
	counters.FlowPackets.Add(labels, flow.GetPackets())
	counters.FlowBytes.Add(labels, flow.GetBytes())
}

func filterPopularPorts(port uint32) (uint32, string) {
	switch port {
	case 80, 443:
		return port, "www"
	case 22:
		return port, "ssh"
	case 53:
		return port, "dns"
	case 25, 465:
		return port, "smtp"
	case 110, 995:
		return port, "pop3"
	case 143, 993:
		return port, "imap"
	}
	return 0, ""
}
