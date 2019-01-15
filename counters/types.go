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

type Counter struct {
	Fields map[uint32]CounterItems
	Name   string
	Access *sync.Mutex
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
	counter.Access.Unlock()
}

// COUNTER ITEMS

type CounterItems struct {
	Label Label
	Value uint64
}
