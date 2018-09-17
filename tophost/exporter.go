package tophost

import (
	"fmt"
	"sync"
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
				exporter.export(prometheus.TopHostTypeBytes, &rawWindowBytes, &rawTotalBytes)
				exporter.export(prometheus.TopHostTypeConnections, &rawWindowConnections, &rawTotalConnections)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Consider adds new flow to tophost exporter
func (exporter *Exporter) Consider(input Input) {
	countHostTraffic(input.getIdentifier(), input.Bytes)
	countHostConnections(input.getIdentifier())
}

// runs one export cycle of current snapshot
func (exporter *Exporter) export(topHostType prometheus.TopHostType, rawWindow *sync.Map, rawTotal *sync.Map) {
	var tophosts topHosts
	if topHostType == prometheus.TopHostTypeBytes {
		tophosts = exporter.hostlistBytes
	} else if topHostType == prometheus.TopHostTypeConnections {
		tophosts = exporter.hostlistConnections
	} else {
		fmt.Printf("unknown TopHostType %v", topHostType)
	}

	// copy all previous hosts
	previousHosts := make([]string, exporter.maxHosts)
	for i, host := range tophosts {
		previousHosts[i] = host.identifier
	}

	// create empty top host list
	tophosts = make(topHosts, exporter.maxHosts)

	// walk through rawHost list
	length := 0
	rawWindow.Range(func(key, value interface{}) bool {
		length++

		// check if in top N
		currentIdentifier := key.(string)
		currentValue := value.(uint64)
		tophosts.addHost(currentIdentifier, currentValue)

		// remove from rawHosts list
		rawWindow.Delete(currentIdentifier)

		return true
	})

	// push hostlist to promExporter
	for _, host := range tophosts {
		hostInput := splitIdentifier(host.identifier)
		counterValueRaw, ok := rawTotal.Load(host.identifier)
		if !ok {
			// fmt.Printf("Skipping hostInput %v for host %v ... tophosts: %+v \n", hostInput, host, tophosts)
			continue
		}
		counterValue := counterValueRaw.(uint64)
		exporter.promExporter.TopHost(topHostType, hostInput.IPSrc, hostInput.IPDst, hostInput.Peer, counterValue)
		rawTotal.Store(host.identifier, 0) // Reset total counter
		for i, hostIdentifier := range previousHosts {
			if hostIdentifier == host.identifier {
				previousHosts[i] = ""
			}
		}
	}

	if topHostType == prometheus.TopHostTypeBytes {
		exporter.hostlistBytes = tophosts
	} else if topHostType == prometheus.TopHostTypeConnections {
		exporter.hostlistConnections = tophosts
	} else {
		fmt.Printf("unknown TopHostType %v", topHostType)
	}

	// find and report removed hosts
	for _, hostIdentifier := range previousHosts {
		if hostIdentifier != "" {
			hostInput := splitIdentifier(hostIdentifier)
			exporter.promExporter.RemoveTopHost(topHostType, hostInput.IPSrc, hostInput.IPDst, hostInput.Peer)
		}
	}

}
