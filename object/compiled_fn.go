package object

import (
	"crabscript.rs/code"
	"fmt"
)

type CompFn struct {
	Instructions  code.Instructions // set of instructions to call when fn is called
	LocalVarCount int               // Count of variables bound inside the fn
}

func (cf *CompFn) Type() ObjectType {
	return CompFnObj
}

func (cf *CompFn) Inspect() string {
	return fmt.Sprintf("CompiledFn[%p]", cf)
}
