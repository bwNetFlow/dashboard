package tophost

import (
	"time"

	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/prometheus"
)

// Exporter handles top host counting and export
type Exporter struct {
	promExporter        prometheus.Exporter
	hostlistBytes       topHosts
	hostlistConnections topHosts
	maxHosts            int
}

// Initialize the top host exporter
func (exporter *Exporter) Initialize(promExporter prometheus.Exporter, maxHosts int, exportInterval time.Duration) {
	exporter.promExporter = promExporter
	exporter.maxHosts = maxHosts
	ticker := time.NewTicker(exportInterval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				exporter.exportTraffic()
				exporter.exportConnections()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Consider adds new flow to tophost exporter
func (exporter *Exporter) Consider(ipSrc string, ipDst string, bytes uint64) {
	countHostTraffic(ipSrc, bytes)
	countHostTraffic(ipDst, bytes)
	countHostConnections(ipSrc)
	countHostConnections(ipDst)
}

// runs one export cycle of current snapshot
func (exporter *Exporter) exportTraffic() {
	// create empty top host list
	exporter.hostlistBytes = make(topHosts, exporter.maxHosts)

	// walk through rawHost list
	length := 0
	rawHostsBytes.Range(func(key, value interface{}) bool {
		length++

		// check if in top N
		currentIP := key.(string)
		currentValue := value.(uint64)
		exporter.hostlistBytes.addHost(currentIP, currentValue)

		// remove from rawHosts list
		rawHostsBytes.Delete(currentIP)

		return true
	})

	// push hostlist to promExporter
	for _, host := range exporter.hostlistBytes {
		exporter.promExporter.TopHostTraffic(host.ip, host.value)
	}
}

// runs one export cycle of current snapshot
func (exporter *Exporter) exportConnections() {
	// create empty top host list
	exporter.hostlistConnections = make(topHosts, exporter.maxHosts)

	// walk through rawHost list
	length := 0
	rawHostsConnections.Range(func(key, value interface{}) bool {
		length++

		// check if in top N
		currentIP := key.(string)
		currentValue := value.(uint64)
		exporter.hostlistConnections.addHost(currentIP, currentValue)

		// remove from rawHosts list
		rawHostsConnections.Delete(currentIP)

		return true
	})

	// push hostlist to promExporter
	for _, host := range exporter.hostlistConnections {
		exporter.promExporter.TopHostConnections(host.ip, host.value)
	}
}
