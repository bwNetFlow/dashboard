package counters

import (
	"bytes"
	"hash/fnv"
	"sort"
)

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

func (label *Label) Hash() uint32 {
	return label.hash()
}

func (label *Label) hash() uint32 {
	var buffer bytes.Buffer
	keys := make([]string, 0, len(label.Fields))
	for k := range label.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := label.Fields[k]
		buffer.WriteString(k + ":" + v + ";")
	}
	h := fnv.New32a()
	h.Write(buffer.Bytes())
	return h.Sum32()
}
