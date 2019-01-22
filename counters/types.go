package counters

import (
	"bytes"
	"hash/fnv"
	"sync"
)

// LABEL

func NewLabel(fields map[string]string) Label {
	return Label{
		Fields: fields,
	}
}

func NewEmptyLabel() Label {
	return Label{
		Fields: make(map[string]string),
	}
}

type Label struct {
	Fields map[string]string
}

func (label *Label) hash() uint32 {
	var buffer bytes.Buffer
	for k, v := range label.Fields {
		buffer.WriteString(k + ":" + v + ";")
	}
	h := fnv.New32a()
	h.Write(buffer.Bytes())
	return h.Sum32()
}

// COUNTER

func NewCounter(name string) Counter {
	return Counter{
		Fields:        make(map[uint32]CounterItems),
		CustomerIndex: make(map[string][]uint32),
		Name:          name,
		Access:        &sync.Mutex{},
	}
}

type Counter struct {
	Fields        map[uint32]CounterItems
	CustomerIndex map[string][]uint32
	Name          string
	Access        *sync.Mutex
}

func (counter *Counter) Add(label Label, value uint64) {
	counter.Access.Lock()
	h := label.hash()
	item, ok := counter.Fields[h]
	if !ok {
		counter.Fields[h] = CounterItems{
			Label: label,
			Value: uint64(0),
		}
		item = counter.Fields[h]
	}
	item.Value += value
	counter.Fields[h] = item
	counter.addCustomerIndex(label, h)
	counter.Access.Unlock()
}

func (counter *Counter) addCustomerIndex(label Label, hash uint32) {
	cid, ok := label.Fields["cid"]
	if !ok {
		// no cid found in labels
		return
	}
	hashList, ok := counter.CustomerIndex[cid]
	if !ok {
		hashList = make([]uint32, 0)
	}
	hashList = append(hashList, hash)
	counter.CustomerIndex[cid] = hashList
}

func (counter *Counter) Delete(label Label) {
	counter.Access.Lock()
	h := label.hash()
	_, ok := counter.Fields[h]
	if ok {
		counter.delCustomerIndex(label, h)
		delete(counter.Fields, h)
	}
	counter.Access.Unlock()
}

func (counter *Counter) delCustomerIndex(label Label, hash uint32) {
	cid, ok := label.Fields["cid"]
	if !ok {
		// no cid found in labels
		return
	}
	hashList, ok := counter.CustomerIndex[cid]
	if !ok {
		// nothing to do
		return
	}

	hashList = removeFromArray(hashList, hash)
	counter.CustomerIndex[cid] = hashList
}

// COUNTER ITEMS

type CounterItems struct {
	Label Label
	Value uint64
}

// UTIL FN

func removeFromArray(s []uint32, r uint32) []uint32 {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
