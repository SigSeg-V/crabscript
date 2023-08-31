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
	}

	return nil
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)

	if !ok {
		return newError("identifier not found: %s", node.String())
	}

	return val
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
