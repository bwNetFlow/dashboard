package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
)

// TOPN defines how many hosts are exported
var TOPN = 50
var hostlist topHosts

var byHostBytes sync.Map // concurrent map of map[string]uint64
func countHostTraffic(ip string, bytes uint64) {
	if value, ok := byHostBytes.Load(ip); ok {
		currentBytes := value.(uint64)
		bytes = bytes + currentBytes
	}
	byHostBytes.Store(ip, bytes)
}

func printTopHostList() {
	length := 0
	byHostBytes.Range(func(key, value interface{}) bool {
		length++

		// check if in top N
		currentIP := key.(string)
		currentBytes := value.(uint64)
		hostlist.addToTopN(host{
			ip:    currentIP,
			bytes: currentBytes,
		})

		return true
	})
	fmt.Printf("byHostBytes length: %d\n", length)

	// push hostlist to promExporter
	for _, host := range hostlist {
		promExporter.TopHost(host.ip, host.bytes)
	}
	printMemUsage()
}

func startPeriodicHostExport() {

	hostlist = make(topHosts, TOPN)

	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				printTopHostList()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

/////////
// see https://play.golang.org/p/e6SQCR4cu-1
/////////
type host struct {
	ip    string
	bytes uint64
}

type topHosts []host

func (t topHosts) addToTopN(host host) {
	if host.bytes <= t[0].bytes {
		// Doesn't belong on the list
		return
	}

	// Find the insertion position
	pos := sort.Search(len(t), func(a int) bool {
		return t[a].bytes > host.bytes
	})

	// Shift lower elements
	copy(t[:pos-1], t[1:pos])

	// Insert it
	t[pos-1] = host
}
