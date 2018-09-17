package tophost

import (
	"reflect"
	"sync"
)

var rawWindowBytes sync.Map       // concurrent map of map[string]uint64
var rawWindowConnections sync.Map // concurrent map of map[string]uint64
var rawTotalBytes sync.Map        // concurrent map of map[string]uint64
var rawTotalConnections sync.Map  // concurrent map of map[string]uint64

func countHostTraffic(identifier string, bytes uint64) {
	addToMap(&rawWindowBytes, identifier, bytes)
	addToMap(&rawTotalBytes, identifier, bytes)
}

func countHostConnections(identifier string) {
	addToMap(&rawWindowConnections, identifier, 1)
	addToMap(&rawTotalConnections, identifier, 1)
}

func addToMap(rawmap *sync.Map, identifier string, value uint64) {
	if currentValueRaw, ok := rawmap.Load(identifier); ok {
		if reflect.TypeOf(currentValueRaw).Kind() == reflect.Uint64 {
			currentValue := currentValueRaw.(uint64)
			value = value + uint64(currentValue)
		} else {
			// fmt.Printf("cannot add to map: %v, type %v not of type uint64 \n", currentValueRaw, reflect.TypeOf(currentValueRaw))
			rawmap.Store(identifier, 0)
		}
	}
	rawmap.Store(identifier, value)
}
