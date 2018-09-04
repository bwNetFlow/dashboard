package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	flow "omi-gitlab.e-technik.uni-ulm.de/bwnetflow/bwnetflow_api/go"
)

// Increment updates the counters by one flow
func (exporter *Exporter) Increment(flow *flow.FlowMessage) {
	labels := prometheus.Labels{
		// "src_port":  fmt.Sprint(flow.GetSrcPort()),
		// "dst_port":  fmt.Sprint(flow.GetDstPort()),
		// "proto":     fmt.Sprint(flow.GetProto()),
		// "src_if": fmt.Sprint(flow.GetSrcIf()),
		// "dst_if": fmt.Sprint(flow.GetDstIf()),
		// "direction": fmt.Sprint(flow.GetDirection()),
		"cid": fmt.Sprint(flow.GetCid()),
	}

	flowNumber.With(labels).Inc()
	flowPackets.With(labels).Add(float64(flow.GetPackets()))
	flowBytes.With(labels).Add(float64(flow.GetBytes()))
}
