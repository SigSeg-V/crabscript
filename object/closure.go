package object

import "fmt"

type Closure struct {
	Fn   *CompFn
	Free []Object // free variables bound to the scope
}

func (c *Closure) Type() ObjectType {
	return ClosureObj
}

func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}
