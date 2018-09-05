package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
)

// Increment updates the counters by one flow
func (exporter *Exporter) Increment(flow *flow.FlowMessage) {

	var application string
	srcPort, appGuess1 := filterPopularPorts(flow.GetSrcPort())
	dstPort, appGuess2 := filterPopularPorts(flow.GetDstPort())
	if appGuess1 != "" {
		application = appGuess1
	} else if appGuess2 != "" {
		application = appGuess2
	}

	labels := prometheus.Labels{
		"src_port":    fmt.Sprint(srcPort),
		"dst_port":    fmt.Sprint(dstPort),
		"application": application,
		"proto":       fmt.Sprint(flow.GetProto()),
		"direction":   fmt.Sprint(flow.GetDirection()),
		"cid":         fmt.Sprint(flow.GetCid()),
		"peer":        flow.GetPeer(),
	}

	flowNumber.With(labels).Add(float64(flow.GetSamplingRate()))
	flowPackets.With(labels).Add(float64(flow.GetPackets()))
	flowBytes.With(labels).Add(float64(flow.GetBytes()))
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
