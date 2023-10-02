package compiler

import (
	"fmt"
	"sort"

	"crabscript.rs/ast"
	"crabscript.rs/code"
	"crabscript.rs/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	lastInstruction     EmittedInstruction // last emitted op
	previousInstruction EmittedInstruction // 2nd last emitted op

	symbolTable *SymbolTable // storing variables
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
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
		c.emit(code.OpPop) // Remove result from stack when not needed

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpNeg)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.InfixExpression:
		// reorder left and right values to reuse OpGt
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGt)
			return nil
		}

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
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGt)
		case "==":
			c.emit(code.OpEq)
		case "!=":
			c.emit(code.OpNe)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		// using a number fresh from my ass that will be back patched later
		jmpNtPos := c.emit(code.OpJmpNt, 9999)
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		// remove extra pop so that if blocks can be used for assignment
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}
		// yet another number fresh from my ass
		jmpPos := c.emit(code.OpJmp, 9999)
		// get point to jmp to if condition is not true
		afterConsequencePos := len(c.instructions)
		// back patch the jmp length
		c.changeOperand(jmpNtPos, afterConsequencePos)
		// add jump target
		if node.Alternative == nil {
			c.emit(code.OpNull) // jmp target is null when no branch
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}
		afterAlternativePos := len(c.instructions)
		c.changeOperand(jmpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, st := range node.Statements {
			err := c.Compile(st)
			if err != nil {
				return err
			}
		}

		// binding a variable
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGbl, symbol.Index)

		// retrieving a bound variable from the store
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("unresolved symbol: %v", node.Value)
		}
		c.emit(code.OpGetGbl, symbol.Index)

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, e := range node.Elements {
			err := c.Compile(e)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))

	case *ast.DictLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		// sorting for testability
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		// get final values for the key and value
		for _, k := range keys {
			if err := c.Compile(k); err != nil {
				return err
			}
			if err := c.Compile(node.Pairs[k]); err != nil {
				return err
			}
		}
		c.emit(code.OpDict, len(node.Pairs)*2)

	case *ast.IndexExpression:
		// get expression of subscript and item
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(code.OpIdx)
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
	pos := c.addInstruction(in)

	c.setLastInstruction(op, pos)
	return pos
}

// Add instructions to stack
// returns position of the added instruction
func (c *Compiler) addInstruction(instructions []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, instructions...)
	return posNewInstruction
}

// populate the cache of previous instructions - used for jmps
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.previousInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{
		Opcode:   op,
		Position: pos,
	}
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

// back patch instructions at position
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

// back patching operation at opPos
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}
