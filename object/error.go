package object

import "fmt"

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ErrorObj
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %q", e.Message)
}
