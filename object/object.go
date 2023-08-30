package object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	IntegerObj = "Integer"
	BooleanObj = "Boolean"
	NullObj    = "Null"
	ReturnObj  = "Return"
	ErrorObj   = "Error"
)
