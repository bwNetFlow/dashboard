package tophost

import "sort"

type host struct {
	ip    string
	value uint64
}

type topHosts []host

func (t topHosts) addHost(ip string, value uint64) {
	host := host{
		ip:    ip,
		value: value,
	}

	if host.value <= t[0].value {
		// Doesn't belong on the list
		return
	}

	// Find the insertion position
	pos := sort.Search(len(t), func(a int) bool {
		return t[a].value > host.value
	})

	// Shift lower elements
	copy(t[:pos-1], t[1:pos])

	// Insert it
	t[pos-1] = host
}
