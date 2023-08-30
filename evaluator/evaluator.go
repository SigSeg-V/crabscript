package evaluator

import (
	"crabscript.rs/ast"
	"crabscript.rs/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return boolToObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)

	if isTruthy(condition) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
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
	default:
		return Null
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
		return Null
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
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
		return Null
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
		return Null
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
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
