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
	OpConst   Opcode = iota // max of 65536 constants in constant pa
	OpAdd                   // add the topmost 2 elem of stack
	OpSub                   // add the topmost 2 elem of stack
	OpMul                   // add the topmost 2 elem of stack
	OpDiv                   // add the topmost 2 elem of stack
	OpPop                   // cleans the stack after an expression
	OpTrue                  // represents `true` literal
	OpFalse                 // represents `false` literal
	OpEq                    // equals comparator
	OpNe                    // not equals comparator
	OpGt                    // greater than comparator
	OpNeg                   // negation operator
	OpBang                  // `not` operator
	OpJmp                   // jump operator, for conditionals and
	OpJmpNt                 // jump when not true, for conditionals
	OpNull                  // *NULL*
	OpGetGbl                // getting bound variables from stack
	OpSetGbl                // setting bound variables from stack
	OpArray                 // list collection type
	OpDict                  // dictionary type
	OpIdx                   // index or subscript operator
	OpCall                  // call fn
	OpRet                   // return to branch point
	OpRetVal                // return value to top of stack
	OpGetLcl                // getting bound varables from the fn stack frame
	OpSetLcl                // setting bound variables from the fn stack frame
	OpGetBIn                // getting built in fns
	OpClosure               // anonymous functions
	OpGetFree               // getting variables from closures
)

// Definition - debugging info and humand readable opcode for the operation
type Definition struct {
	Name          string // readable name for operation
	OperandWidths []int  // number of bytes each operand takes up
}

var definitions = map[Opcode]*Definition{
	OpConst:  {"OpConst", []int{2}},  // max of 65536 constants in constant pool
	OpAdd:    {"OpAdd", []int{}},     // add the topmost 2 elem of stack
	OpSub:    {"OpSub", []int{}},     // add the topmost 2 elem of stack
	OpMul:    {"OpMul", []int{}},     // add the topmost 2 elem of stack
	OpDiv:    {"OpDiv", []int{}},     // add the topmost 2 elem of stack
	OpPop:    {"OpPop", []int{}},     // cleans the stack after an expression is evaluated
	OpTrue:   {"OpTrue", []int{}},    // represents `true` literal
	OpFalse:  {"OpFalse", []int{}},   // represents `false` literal
	OpEq:     {"OpEq", []int{}},      // equals comparator
	OpNe:     {"OpNe", []int{}},      // not equals comparator
	OpGt:     {"OpGt", []int{}},      // greater than comparator
	OpNeg:    {"OpNeg", []int{}},     // negation operator
	OpBang:   {"OpBang", []int{}},    // `not` operator
	OpJmp:    {"OpJmp", []int{2}},    // jump operator, for conditionals and functions
	OpJmpNt:  {"OpJmpNt", []int{2}},  // jump when not true, for conditionals
	OpNull:   {"OpNull", []int{}},    // *NULL*
	OpGetGbl: {"OpGetGbl", []int{2}}, // getting bound variables from the stack
	OpSetGbl: {"OpSetGbl", []int{2}}, // setting bound variables from the stack
	OpArray:  {"OpArray", []int{2}},  // list collection type
	OpDict:   {"OpDict", []int{2}},   // dictionary type
	OpIdx:    {"OpIdx", []int{}},     // index or subscript operator
	OpCall:   {"OpCall", []int{1}},   // call fn - holds number of arguments (max 255)
	OpRet:    {"OpRet", []int{}},     // return to branch point
	OpRetVal: {"OpRetVal", []int{}},  // push value to top of stack
	OpGetLcl: {"OpGetLcl", []int{2}}, // getting bound varables from the fn stack frame
	OpSetLcl: {"OpSetLcl", []int{2}}, // setting bound variables from the fn stack frame
	OpGetBIn: {"OpGetBIn", []int{1}}, // get built in fns
	// anonymous functions, op 1 is const index,
	// op 2 is number of free vars that need to
	//be moved with the closure
	OpClosure: {"OpClosure", []int{2, 1}},
	OpGetFree: {"OpGetFree", []int{1}}, // getting variables from closures
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
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
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
		case 1:
			operands[i] = int(ReadUint8(in[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(instructions Instructions) uint16 {
	return binary.BigEndian.Uint16(instructions)
}

func ReadUint8(instructions Instructions) uint8 {
	return uint8(instructions[0])
}
