package counters

// Counter defines a local or remote counter
type Counter interface {
	Add(label Label, value uint64)
	Delete(label Label)
	GetName() string
	GetFields() map[uint32]CounterItems
	GetField(hash uint32) CounterItems
	GetCids() []string
	GetFieldHashesByCid(cid string) []uint32
}
