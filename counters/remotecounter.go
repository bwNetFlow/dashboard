package counters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwNetFlow/dashboard/cache"
)

// NewRemoteCounter returns a new counter
func NewRemoteCounter(name string, cache *cache.Cache) *RemoteCounter {
	return &RemoteCounter{
		Name:  name,
		cache: cache,
	}
}

type RemoteCounter struct {
	Name  string
	cache *cache.Cache
}

func (remotecounter *RemoteCounter) GetName() string {
	return remotecounter.Name
}

func (remotecounter *RemoteCounter) getKey(label Label) string {
	h := label.hash()
	cid := remotecounter.getCidFromLabel(label)
	key := fmt.Sprintf("%s:%s:%d", remotecounter.Name, cid, h)
	return key
}

func (remotecounter *RemoteCounter) Add(label Label, value uint64) {
	key := remotecounter.getKey(label)
	remotecounter.cache.IncreaseBy(key, value)
	// counter.addCustomerIndex(label, h)
}

func (remotecounter *RemoteCounter) getCidFromLabel(label Label) string {
	cid, ok := label.Fields["cid"]
	if !ok {
		// no cid found in labels
		return ""
	}
	return cid
}

func (remotecounter *RemoteCounter) Delete(label Label) {
	key := remotecounter.getKey(label)
	remotecounter.cache.Delete(key)
}

func (remotecounter *RemoteCounter) GetFields() map[uint32]CounterItems {
	pattern := fmt.Sprintf("%s:*", remotecounter.Name)
	keys := remotecounter.cache.FindKeys(pattern)
	fields := make(map[uint32]CounterItems)
	for _, key := range keys {
		keyParts := strings.Split(key, ":")
		if len(keyParts) < 3 {
			// invalid key
			continue
		}
		tmp, err := strconv.ParseUint(keyParts[2], 10, 32)
		if err != nil {
			// failed to parse
			continue
		}
		hash := uint32(tmp)
		fields[hash] = remotecounter.GetField(hash)
	}
	return fields
}

func (remotecounter *RemoteCounter) getField(key string) CounterItems {
	value := remotecounter.cache.Get(key)
}

func (remotecounter *RemoteCounter) GetField(hash uint32) CounterItems {
	pattern := fmt.Sprintf("%s:*:%d", remotecounter.Name, hash)
	keys := remotecounter.cache.FindKeys(pattern)
	if len(keys) <= 0 {
		// nothing found
		fmt.Printf("No field found for pattern %s\n", pattern)
		return CounterItems{}
	}

	return CounterItems{}
}

func (remotecounter *RemoteCounter) GetCids() []string {
	return make([]string, 0)
}

func (remotecounter *RemoteCounter) GetFieldHashesByCid(cid string) []uint32 {
	return nil
}
