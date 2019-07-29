package counters

import (
	"sort"
	"sync"

	"github.com/jinzhu/copier"
)

// LocalCounter returns a new counter
func NewLocalCounter(name string) *LocalCounter {
	return &LocalCounter{
		Fields:        make(map[uint32]CounterItems),
		CustomerIndex: make(map[string][]uint32),
		Name:          name,
		Access:        &sync.Mutex{},
	}
}

type LocalCounter struct {
	Fields        map[uint32]CounterItems
	CustomerIndex map[string][]uint32
	Name          string
	Access        *sync.Mutex
}

func (localCounter *LocalCounter) GetName() string {
	return localCounter.Name
}

func (localCounter *LocalCounter) Add(label Label, value uint64) {
	localCounter.Access.Lock()
	h := label.hash()
	item, ok := localCounter.Fields[h]
	if !ok {
		localCounter.Fields[h] = CounterItems{
			Label: label,
			Value: uint64(0),
		}
		item = localCounter.Fields[h]
	}
	item.Value += value
	localCounter.Fields[h] = item
	localCounter.addCustomerIndex(label, h)
	localCounter.Access.Unlock()
}

func (localCounter *LocalCounter) addCustomerIndex(label Label, hash uint32) {
	cid, ok := label.Fields["cid"]
	if !ok {
		// no cid found in labels
		return
	}
	if localCounter.isHashInCustomerIndex(cid, hash) {
		// hash is present. nothing to do.
		return
	}

	hashList, ok := localCounter.CustomerIndex[cid]
	if !ok {
		hashList = make([]uint32, 0)
	}

	// not present, add the new hash
	hashList = append(hashList, hash)

	// sort hashList
	sort.Slice(hashList, func(i, j int) bool { return hashList[i] < hashList[j] })

	// store back the extended hash list
	localCounter.CustomerIndex[cid] = hashList
}

func (localCounter *LocalCounter) isHashInCustomerIndex(cid string, hash uint32) bool {
	hashList, ok := localCounter.CustomerIndex[cid]
	if !ok {
		return false
	}
	i := sort.Search(len(hashList), func(i int) bool { return hashList[i] >= hash })
	if i < len(hashList) && hashList[i] == hash {
		return true
	}
	return false
}

func (localCounter *LocalCounter) Delete(label Label) {
	localCounter.Access.Lock()
	h := label.hash()
	_, ok := localCounter.Fields[h]
	if ok {
		localCounter.delCustomerIndex(label, h)
		delete(localCounter.Fields, h)
	}
	localCounter.Access.Unlock()
}

func (localCounter *LocalCounter) delCustomerIndex(label Label, hash uint32) {
	cid, ok := label.Fields["cid"]
	if !ok {
		// no cid found in labels
		return
	}
	hashList, ok := localCounter.CustomerIndex[cid]
	if !ok {
		// nothing to do
		return
	}

	hashList = removeFromArray(hashList, hash)
	localCounter.CustomerIndex[cid] = hashList
}

func (localCounter *LocalCounter) GetFields() map[uint32]CounterItems {
	localCounter.Access.Lock()
	defer localCounter.Access.Unlock()
	// deep copy the field to avoid concurrency issues
	var fields map[uint32]CounterItems
	copier.Copy(localCounter.Fields, &fields)
	return fields
}

func (localCounter *LocalCounter) GetField(hash uint32) CounterItems {
	localCounter.Access.Lock()
	defer localCounter.Access.Unlock()
	// deep copy the field to avoid concurrency issues
	var field CounterItems
	copier.Copy(localCounter.Fields[hash], &field)
	return field
}

func (localCounter *LocalCounter) GetCids() []string {
	localCounter.Access.Lock()
	defer localCounter.Access.Unlock()
	cids := make([]string, 0, len(localCounter.CustomerIndex))
	for k := range localCounter.CustomerIndex {
		cids = append(cids, k)
	}
	return cids
}

func (localCounter *LocalCounter) GetFieldHashesByCid(cid string) []uint32 {
	hashes := localCounter.CustomerIndex[cid]
	return hashes
}
