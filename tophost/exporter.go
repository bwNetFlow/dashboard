package tophost

import (
	"sync"
	"time"

	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/prometheus"
)

// Exporter handles top host counting and export
type Exporter struct {
	promExporter prometheus.Exporter
	maxHosts     int
	counters     sync.Map
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
				// periodical export for each cid
				exporter.counters.Range(func(key, counterRaw interface{}) bool {
					cid, _ := key.(uint32)
					counter := counterRaw.(*Counter)

					exporter.export(prometheus.TopHostTypeBytes, cid, counter)
					exporter.export(prometheus.TopHostTypeConnections, cid, counter)
					return true
				})
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Consider adds new flow to tophost exporter
func (exporter *Exporter) Consider(input Input) {
	// split by cid
	var counter *Counter
	if counterRaw, ok := exporter.counters.Load(input.Cid); ok {
		counter = counterRaw.(*Counter)
	} else {
		counter = NewCounter()
	}
	counter.countHostTraffic(input.getIdentifier(), input.Bytes)
	counter.countHostConnections(input.getIdentifier())
	exporter.counters.Store(input.Cid, counter)
}

// runs one export cycle of current snapshot
func (exporter *Exporter) export(topHostType prometheus.TopHostType, cid uint32, counter *Counter) {
	// get current tophosts
	tophosts := counter.toplist[topHostType]

	// copy all previous hosts
	previousHosts := make([]string, exporter.maxHosts)
	for i, host := range tophosts {
		previousHosts[i] = host.identifier
	}

	// create new, empty top host list
	tophosts = make(topHosts, exporter.maxHosts)

	// walk through rawHost list
	length := 0
	counter.window[topHostType].Range(func(key, value interface{}) bool {
		length++

		// check if in top N
		currentIdentifier := key.(string)
		currentValue := value.(uint64)
		tophosts.addHost(currentIdentifier, currentValue)

		// remove from rawHosts list
		counter.window[topHostType].Delete(currentIdentifier)

		return true
	})

	// push hostlist to promExporter
	for _, host := range tophosts {
		hostInput := splitIdentifier(host.identifier)
		counterValueRaw, ok := counter.total[topHostType].Load(host.identifier)
		if !ok {
			// fmt.Printf("Skipping hostInput %v for host %v ... tophosts: %+v \n", hostInput, host, tophosts)
			continue
		}
		counterValue, ok := counterValueRaw.(uint64)
		if !ok {
			// counterValueRaw not uint64 - skipping
			continue
		}
		exporter.promExporter.TopHost(topHostType, cid, hostInput.IPSrc, hostInput.IPDst, hostInput.Peer, counterValue)
		counter.total[topHostType].Store(host.identifier, 0) // Reset total counter since exported
		for i, hostIdentifier := range previousHosts {
			if hostIdentifier == host.identifier {
				previousHosts[i] = ""
			}
		}
	}

	// save new tophosts
	counter.toplist[topHostType] = tophosts

	// find and report removed hosts
	deleted := 0
	for _, hostIdentifier := range previousHosts {
		if hostIdentifier != "" {
			hostInput := splitIdentifier(hostIdentifier)
			exporter.promExporter.RemoveTopHost(topHostType, cid, hostInput.IPSrc, hostInput.IPDst, hostInput.Peer)
			deleted++
		}
	}

}
