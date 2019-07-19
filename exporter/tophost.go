package exporter

import (
	"fmt"

	"github.com/bwNetFlow/dashboard/counters"
)

// TopHostType defines whether the top host is bytes or connections
type TopHostType int32

const (
	// TopHostTypeBytes defines type byte of TopHostType
	TopHostTypeBytes TopHostType = 0
	// TopHostTypeConnections defines type connections of TopHostType
	TopHostTypeConnections TopHostType = 1
)

// TopHost updates one entry for Top Hosts
func (exporter *Exporter) TopHost(topHostType TopHostType, cid uint32, ipSrc string, ipDst string, peer string, value uint64) {
	labels := counters.NewLabel(map[string]string{
		"cid":   fmt.Sprintf("%d", cid),
		"ipSrc": ipSrc,
		"ipDst": ipDst,
		"peer":  peer,
	})

	var counterVec counters.Counter
	if topHostType == TopHostTypeBytes {
		counterVec = counters.HostBytes
	} else if topHostType == TopHostTypeConnections {
		counterVec = counters.HostConnections
	} else {
		return
	}
	counterVec.Add(labels, value)
}

// RemoveTopHost removes the host from the counter vector
func (exporter *Exporter) RemoveTopHost(topHostType TopHostType, cid uint32, ipSrc string, ipDst string, peer string) {
	labels := counters.NewLabel(map[string]string{
		"cid":   fmt.Sprintf("%d", cid),
		"ipSrc": ipSrc,
		"ipDst": ipDst,
		"peer":  peer,
	})
	var counterVec counters.Counter
	if topHostType == TopHostTypeBytes {
		counterVec = counters.HostBytes
	} else if topHostType == TopHostTypeConnections {
		counterVec = counters.HostConnections
	} else {
		return
	}
	counterVec.Delete(labels)
}
