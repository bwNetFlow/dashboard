package tophost

import "sync"

var rawHostsBytes sync.Map       // concurrent map of map[string]uint64
var rawHostsConnections sync.Map // concurrent map of map[string]uint64

func countHostTraffic(ip string, bytes uint64) {
	if value, ok := rawHostsBytes.Load(ip); ok {
		currentBytes := value.(uint64)
		bytes = bytes + currentBytes
	}
	rawHostsBytes.Store(ip, bytes)
}

func countHostConnections(ip string) {
	counter := uint64(1)
	if value, ok := rawHostsConnections.Load(ip); ok {
		current := value.(uint64)
		counter = counter + current
	}
	rawHostsConnections.Store(ip, counter)
}
