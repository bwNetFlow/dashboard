package exporter

import (
	"fmt"

	"github.com/bwNetFlow/dashboard/counters"
	flow "github.com/bwNetFlow/protobuf/go"
	"github.com/bwNetFlow/protobuf_helpers/go"
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

	flowh := protobuf_helpers.NewFlowHelper(flow)
	labels := counters.NewLabel(map[string]string{
		// "src_port":      fmt.Sprint(srcPort),
		// "dst_port":      fmt.Sprint(dstPort),
		"ipversion":   flowh.IPVersionString(),
		"application": application,
		"protoname":   fmt.Sprint(flow.GetProtoName()),
		"direction":   flowh.FlowDirectionString(),
		"cid":         fmt.Sprint(flow.GetCid()),
		"peer":        flowh.Peer(),
		// "remotecountry": flow.GetRemoteCountry(),
	})

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
