package object

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "len",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got %d, want 1", len(args))
				}

				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		Name: "first",
		Builtin: &Builtin{
			func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got %d, want 1", len(args))
				}

				if args[0].Type() != ArrayObj && args[0].Type() != StringObj {
					return newError("argument to `first` is invalid, got %s", args[0].Type())
				}

				switch args[0].(type) {
				case *Array:
					arr := args[0].(*Array)
					if len(arr.Elements) > 0 {
						return arr.Elements[0]
					}
					return nil
				case *String:
					str := args[0].(*String)
					if len(str.Value) > 0 {
						foundRune, _ := utf8.DecodeRune([]byte(str.Value))
						return &String{Value: string(foundRune)}
					}
					return nil
				}
				return nil
			},
		},
	},
	{
		Name: "last",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got %d, want 1", len(args))
				}

				if args[0].Type() != ArrayObj && args[0].Type() != StringObj {
					return newError("argument to `last` is invalid, got %s", args[0].Type())
				}

				switch args[0].(type) {
				case *Array:
					arr := args[0].(*Array)
					if len(arr.Elements) > 0 {
						return arr.Elements[len(arr.Elements)-1]
					}
					return nil
				case *String:
					str := args[0].(*String)
					if len(str.Value) > 0 {
						foundRune, _ := utf8.DecodeLastRune([]byte(str.Value))
						return &String{Value: string(foundRune)}
					}
					return nil
				}
				return nil
			},
		},
	},
	{
		Name: "tail",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got %d, want 1", len(args))
				}

				if args[0].Type() != ArrayObj && args[0].Type() != StringObj {
					return newError("argument to `tail` is invalid, got %s", args[0].Type())
				}

				switch args[0].(type) {
				case *Array:
					arr := args[0].(*Array)
					if len(arr.Elements) > 1 {
						return &Array{Elements: arr.Elements[1:]}
					}
					return nil
				case *String:
					str := args[0].(*String)
					if len(str.Value) > 0 {
						_, runeSize := utf8.DecodeRune([]byte(str.Value))
						return &String{Value: str.Value[runeSize:]}
					}
					return nil
				}
				return nil
			},
		},
	},
	{
		Name: "push",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got %d want %d", len(args), 2)
				}

				if args[0].Type() != ArrayObj {
					return newError("argument to `push` must be Array, got %s", args[0].Type())
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)

				newElements := make([]Object, length+1, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &Array{Elements: newElements}
			},
		},
	},
	{
		Name: "puts",
		Builtin: &Builtin{
			Fn: func(args ...Object) Object {
				var out bytes.Buffer
				for _, obj := range args {
					out.WriteString(obj.Inspect())
				}
				println(out.String())
				return nil
			},
		},
	},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
