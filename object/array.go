package object

import (
	"bytes"
	"strings"
)

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType {
	return ArrayObj
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elem := []string{}
	for _, e := range ao.Elements {
		elem = append(elem, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elem, ", "))
	out.WriteString("]")

	return out.String()
}
