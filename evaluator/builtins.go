package evaluator

import (
	"crabscript.rs/object"
	"unicode/utf8"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d, want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d, want 1", len(args))
			}

			if args[0].Type() != object.ArrayObj && args[0].Type() != object.StringObj {
				return newError("argument to `first` is invalid, got %s", args[0].Type())
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return Null
			case *object.String:
				str := args[0].(*object.String)
				if len(str.Value) > 0 {
					foundRune, _ := utf8.DecodeRune([]byte(str.Value))
					return &object.String{Value: string(foundRune)}
				}
				return Null
			}
			return Null
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d, want 1", len(args))
			}

			if args[0].Type() != object.ArrayObj && args[0].Type() != object.StringObj {
				return newError("argument to `last` is invalid, got %s", args[0].Type())
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[len(arr.Elements)-1]
				}
				return Null
			case *object.String:
				str := args[0].(*object.String)
				if len(str.Value) > 0 {
					foundRune, _ := utf8.DecodeLastRune([]byte(str.Value))
					return &object.String{Value: string(foundRune)}
				}
				return Null
			}
			return Null
		},
	},
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d, want 1", len(args))
			}

			if args[0].Type() != object.ArrayObj && args[0].Type() != object.StringObj {
				return newError("argument to `tail` is invalid, got %s", args[0].Type())
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				if len(arr.Elements) > 1 {
					return &object.Array{Elements: arr.Elements[1:]}
				}
				return Null
			case *object.String:
				str := args[0].(*object.String)
				if len(str.Value) > 0 {
					_, runeSize := utf8.DecodeRune([]byte(str.Value))
					return &object.String{Value: str.Value[runeSize:]}
				}
				return Null
			}
			return Null
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got %d want %d", len(args), 2)
			}

			if args[0].Type() != object.ArrayObj {
				return newError("argument to `push` must be Array, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
}
