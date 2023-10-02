package object

import (
	"crabscript.rs/code"
	"fmt"
)

type CompFn struct {
	Instructions code.Instructions // set of instructions to call when fn is called
}

func (cf *CompFn) Type() ObjectType {
	return CompFnObj
}

func (cf *CompFn) Inspect() string {
	return fmt.Sprintf("CompiledFn[%p]", cf)
}
