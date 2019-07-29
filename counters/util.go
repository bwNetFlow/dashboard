package counters

func removeFromArray(s []uint32, r uint32) []uint32 {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
