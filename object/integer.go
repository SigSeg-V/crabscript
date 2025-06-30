package object

import "fmt"

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%v", i.Value)
}

func (i *Integer) Type() ObjectType {
	return IntegerObj
}
