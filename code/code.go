package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// byte code instructions used for the VM

type Instructions []byte

type Opcode byte

// Enum of opcodes in use
const (
	OpConstant Opcode = iota // max of 65536 constants in constant pa
	OpAdd                    // add the topmost 2 elem of stack
	OpSub                    // add the topmost 2 elem of stack
	OpMul                    // add the topmost 2 elem of stack
	OpDiv                    // add the topmost 2 elem of stack
	OpPop                    // cleans the stack after an expression
	OpTrue                   // represents `true` literal
	OpFalse                  // represents `false` literal
	OpEq                     // equals comparator
	OpNe                     // not equals comparator
	OpGt                     // greater than comparator
	OpNeg                    // negation operator
	OpBang                   // `not` operator
	OpJmp                    // jump operator, for conditionals and
	OpJmpNt                  // jump when not true, for conditionals
	OpNull                   // *NULL*
)

// Definition - debugging info and humand readable opcode for the operation
type Definition struct {
	Name          string // readable name for operation
	OperandWidths []int  // number of bytes each operand takes up
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // max of 65536 constants in constant pool
	OpAdd:      {"OpAdd", []int{}},       // add the topmost 2 elem of stack
	OpSub:      {"OpSub", []int{}},       // add the topmost 2 elem of stack
	OpMul:      {"OpMul", []int{}},       // add the topmost 2 elem of stack
	OpDiv:      {"OpDiv", []int{}},       // add the topmost 2 elem of stack
	OpPop:      {"OpPop", []int{}},       // cleans the stack after an expression is evaluated
	OpTrue:     {"OpTrue", []int{}},      // represents `true` literal
	OpFalse:    {"OpFalse", []int{}},     // represents `false` literal
	OpEq:       {"OpEq", []int{}},        // equals comparator
	OpNe:       {"OpNe", []int{}},        // not equals comparator
	OpGt:       {"OpGt", []int{}},        // greater than comparator
	OpNeg:      {"OpNeg", []int{}},       // negation operator
	OpBang:     {"OpBang", []int{}},      // `not` operator
	OpJmp:      {"OpJmp", []int{2}},      // jump operator, for conditionals and functions
	OpJmpNt:    {"OpJmpNt", []int{2}},    // jump when not true, for conditionals
	OpNull:     {"OpNull", []int{}},      // *NULL*
}

// Lookup returns relevant debugging info for op if available
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d not defined", op)
	}
	return def, nil
}

func (instructions Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(instructions) {
		def, err := Lookup(instructions[i])
		if err != nil {
			_, _ = fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, instructions[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, instructions.fmtInstruction(def, operands))

		i += read + 1
	}

	return out.String()
}

func (instructions Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand length %d does not match defined %d", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func ReadOperands(def *Definition, in Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(in[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(instructions Instructions) uint16 {
	return binary.BigEndian.Uint16(instructions)
}
