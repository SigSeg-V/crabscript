package evaluator

import (
	"crabscript.rs/ast"
	"crabscript.rs/object"
	"fmt"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		eval := Eval(node.ReturnValue, env)
		if isError(eval) {
			return eval
		}
		return &object.ReturnValue{Value: eval}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// adding / modifying val on heap
		env.Set(node.Name.String(), val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return boolToObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return callFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.DictLiteral:
		return evalDictLiteral(node, env)
	}

	return nil
}

func evalDictLiteral(node *ast.DictLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.DictKey]object.DictPair)

	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		dictKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable hash key: %s", key.Type())
		}

		val := Eval(valNode, env)
		if isError(val) {
			return val
		}

		hashed := dictKey.DictKey()
		pairs[hashed] = object.DictPair{Key: key, Value: val}
	}

	return &object.Dict{Pairs: pairs}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.DictObj:
		return evalDictIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalDictIndexExpression(left object.Object, index object.Object) object.Object {
	dictObject := left.(*object.Dict)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := dictObject.Pairs[key.DictKey()]
	if !ok {
		return Null
	}

	return pair.Value
}

func evalArrayIndexExpression(left object.Object, index object.Object) object.Object {
	arrayObj := left.(*object.Array)
	inx := index.(*object.Integer).Value
	imax := int64(len(arrayObj.Elements) - 1)

	if inx < 0 || inx > imax {
		return Null
	}

	return arrayObj.Elements[inx]
}

func callFunction(function object.Object, args []object.Object) object.Object {
	switch function := function.(type) {
	// user defined fns
	case *object.Function:
		extendedEnv := extendFnEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturnVal(evaluated)

	// builtin interpreter fns
	case *object.Builtin:
		if res := function.Fn(args...); res != nil {
			return res
		}
		return Null

	default:
		return newError("not a function: %s", function.Type())
	}
}

func unwrapReturnVal(obj object.Object) object.Object {
	if ret, ok := obj.(*object.ReturnValue); ok {
		return ret.Value
	}
	return obj
}

func extendFnEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for pi, param := range fn.Parameters {
		env.Set(param.Value, args[pi])
	}

	return env
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range expressions {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin := object.GetBuiltinByName(node.String()); builtin != nil {
		return builtin
	}

	return newError("identifier not found: %s", node.String())
}

// used to evaluate blocks and nested blocks to solve returning wrong value
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		// retrieve inner return value if any
		if result != nil {
			rt := result.Type()
			if rt == object.ReturnObj || rt == object.ErrorObj {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return Null
	}
}

func isTruthy(condition object.Object) bool {
	switch condition {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	// int -> int ops
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringInfixExpression(operator, left, right)
	// any -> bool ops
	case operator == "==":
		return boolToObject(left == right)
	case operator == "!=":
		return boolToObject(left != right)
	// error handling
	case left.Type() != right.Type():
		return newError("types not matching: %s and %s", left.Type(), right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {

	// int ops returning ints
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}

	// int ops returning bools
	case "<":
		return boolToObject(leftValue < rightValue)
	case ">":
		return boolToObject(leftValue > rightValue)
	case "==":
		return boolToObject(leftValue == rightValue)
	case "!=":
		return boolToObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalProgram(pgm *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range pgm.Statements {
		result = Eval(statement, env)

		// unwrapping return values
		switch result.(type) {
		case *object.ReturnValue:
			return result.(*object.ReturnValue).Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator : %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return newError("unknown operator: %s%s", "-", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == object.ErrorObj
}

// Cache for common simple objects
var (
	// Boolean objects
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}

	Null = &object.Null{}
)

func boolToObject(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}
