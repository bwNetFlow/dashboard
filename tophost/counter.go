package tophost

import (
	"reflect"
	"sync"

	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/prometheus"
)

// Counter describes one counter for one cid tophost list for bytes / connections
type Counter struct {
	window  []sync.Map // array index equals prometheus.TopHostType, map of map[string]uint64 with counters reset every cycle
	total   []sync.Map // array index equals prometheus.TopHostType, map of map[string]uint64 with counters reset only when exported
	toplist []topHosts // array index equals prometheus.TopHostType, toplist contains current top n hosts
}

func (counter *Counter) countHostTraffic(identifier string, bytes uint64) {
	counter.addToMap(&counter.window[prometheus.TopHostTypeBytes], identifier, bytes)
	counter.addToMap(&counter.total[prometheus.TopHostTypeBytes], identifier, bytes)
}

func (counter *Counter) countHostConnections(identifier string) {
	counter.addToMap(&counter.window[prometheus.TopHostTypeConnections], identifier, 1)
	counter.addToMap(&counter.total[prometheus.TopHostTypeConnections], identifier, 1)
}

func (counter *Counter) addToMap(rawmap *sync.Map, identifier string, value uint64) {
	if currentValueRaw, ok := rawmap.Load(identifier); ok {
		if reflect.TypeOf(currentValueRaw).Kind() == reflect.Uint64 {
			currentValue := currentValueRaw.(uint64)
			value = value + uint64(currentValue)
		} else {
			// fmt.Printf("cannot add to map: %v, type %v not of type uint64 \n", currentValueRaw, reflect.TypeOf(currentValueRaw))
			rawmap.Store(identifier, 0)
			return
		}
	}
	rawmap.Store(identifier, value)
}

// NewCounter initializes a new counter
func NewCounter() *Counter {
	return &Counter{
		window:  make([]sync.Map, 2), // 2 ==> bytes (0) + connections (1)
		total:   make([]sync.Map, 2),
		toplist: make([]topHosts, 2),
	}
}
