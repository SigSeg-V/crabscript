package compiler

import (
	"crabscript.rs/ast"
	"crabscript.rs/code"
	"crabscript.rs/object"
	"fmt"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// TODO: Write compiler... lol
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, st := range node.Statements {
			err := c.Compile(st)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		// check operator in infix position
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		default:
			return fmt.Errorf("unknown opeator: %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// Adds a constant to the constant pool
// Returns the new location of the stack pointer (end of the array)
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// Generate an instruction and add to results
// Returns the starting position of the new instruction
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	in := code.Make(op, operands...)
	return c.addInstruction(in)
}

// Add instructions to stack
// returns position of the added instruction
func (c *Compiler) addInstruction(instructions []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, instructions...)
	return posNewInstruction
}
