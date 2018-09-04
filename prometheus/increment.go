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
		"src_if":    fmt.Sprint(flow.GetSrcIf()),
		"dst_if":    fmt.Sprint(flow.GetDstIf()),
		"direction": fmt.Sprint(flow.GetDirection()),
		"cid":       fmt.Sprint(flow.GetCid()),
	}

	flowNumber.With(labels).Inc()
	flowPackets.With(labels).Add(float64(flow.GetPackets()))
	flowBytes.With(labels).Add(float64(flow.GetBytes()))
}

func isPopularPort(port uint32) bool {
	switch port {
	// www
	case 80, 443:
	// ssh
	case 22:
	// dns
	case 53:
	// smtp
	case 25, 465:
	// pop3
	case 110, 995:
	// imap
	case 143, 993:
		return true
	}
	return false
}
