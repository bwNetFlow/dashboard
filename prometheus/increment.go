package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
)

// Increment updates the counters by one flow
func (exporter *Exporter) Increment(flow *flow.FlowMessage) {

	srcPort := uint32(0)
	if isPopularPort(flow.GetSrcPort()) {
		srcPort = flow.GetSrcPort()
	}
	dstPort := uint32(0)
	if isPopularPort(flow.GetDstPort()) {
		dstPort = flow.GetDstPort()
	}

	labels := prometheus.Labels{
		"src_port":  fmt.Sprint(srcPort),
		"dst_port":  fmt.Sprint(dstPort),
		"proto":     fmt.Sprint(flow.GetProto()),
		"direction": fmt.Sprint(flow.GetDirection()),
		"cid":       fmt.Sprint(flow.GetCid()),
		"peer":      flow.GetPeer(),
	}

	flowNumber.With(labels).Add(float64(flow.GetSamplingRate()))
	flowPackets.With(labels).Add(float64(flow.GetPackets()))
	flowBytes.With(labels).Add(float64(flow.GetBytes()))
}

func isPopularPort(port uint32) bool {
	switch port {
	// www
	case 80, 443:
		return true
	// ssh
	case 22:
		return true
	// dns
	case 53:
		return true
	// smtp
	case 25, 465:
		return true
	// pop3
	case 110, 995:
		return true
	// imap
	case 143, 993:
		return true
	}
	return false
}
