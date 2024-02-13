package object

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BuiltinObj
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}
