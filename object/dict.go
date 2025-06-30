package object

import (
	"bytes"
	"hash/fnv"
	"strings"
)

type Hashable interface {
	DictKey() DictKey
}

// Dict stuff
type Dict struct {
	Pairs map[DictKey]DictPair
}

func (d *Dict) Type() ObjectType {
	return DictObj
}

func (d *Dict) Inspect() string {
	var out bytes.Buffer

	elem := []string{}
	for _, val := range d.Pairs {
		elem = append(elem, val.Key.Inspect()+":"+val.Value.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elem, ", "))
	out.WriteString("]")

	return out.String()
}

// DictPair stuff
type DictPair struct {
	Key   Object
	Value Object
}

// DictKey stuff
type DictKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) DictKey() DictKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 0
	}

	return DictKey{Type: b.Type(), Value: val}
}

func (i *Integer) DictKey() DictKey {
	return DictKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) DictKey() DictKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return DictKey{Type: s.Type(), Value: h.Sum64()}
}

func (a *Array) DictKey() DictKey {
	h := fnv.New64a()
	hashString := a.Inspect()
	h.Write([]byte(hashString))

	return DictKey{Type: a.Type(), Value: h.Sum64()}
}
