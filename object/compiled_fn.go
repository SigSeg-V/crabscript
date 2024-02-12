package object

import (
	"crabscript.rs/code"
	"fmt"
)

type CompFn struct {
	Instructions  code.Instructions // set of instructions to call when fn is called
	LocalVarCount int               // count of variables bound inside the fn
	ParamCount    int               // count of params expected in the fn
}

func (cf *CompFn) Type() ObjectType {
	return CompFnObj
}

func (cf *CompFn) Inspect() string {
	return fmt.Sprintf("CompiledFn[%p]", cf)
}
