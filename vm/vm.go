package vm

import (
	"crabscript.rs/code"
	"crabscript.rs/compiler"
	"crabscript.rs/object"
	"fmt"
)

const StackSize = 2048
const GlobalSize = 65536 // maximum number of global variables allowed

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

func boolToObject(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}

type Vm struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	globals      []object.Object
	sp           int // stack pointer, always at next free slot at top of stack

}

func New(bytecode *compiler.Bytecode) *Vm {
	return &Vm{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		globals:      make([]object.Object, GlobalSize),
		sp:           0,
	}
}

// for the repl
func NewWithGblStore(bytecode *compiler.Bytecode, s []object.Object) *Vm {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

// Retrieves the object at the top of the stack
func (vm *Vm) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

// executes bytecode loaded
func (vm *Vm) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		// decoding operations
		switch op {
		case code.OpConstant:
			constIdx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			// putting new const onto stack
			err := vm.push(vm.constants[constIdx])
			if err != nil {
				return err
			}

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.execBinaryOp(op); err != nil {
				return err
			}

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpNe, code.OpEq, code.OpGt:
			if err := vm.execComparison(op); err != nil {
				return err
			}

		case code.OpNeg:
			if err := vm.execNegation(); err != nil {
				return err
			}

		case code.OpBang:
			if err := vm.execBoolNegation(); err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpJmp:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1

		case code.OpJmpNt:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2 // move to the condition (past 2 bytes of constants)
			condition := vm.pop()
			// evaluate the condition and jmp if needed
			if !isTruthy(condition) {
				ip = pos - 1
			}

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetGbl:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			// add var to global heap
			vm.globals[globalIndex] = vm.pop()

			// putting the binding to the top of stack
		case code.OpGetGbl:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

			// initialising array
		case code.OpArray:
			numElem := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			array := vm.buildArray(vm.sp-numElem, vm.sp)
			vm.sp = vm.sp - numElem

			if err := vm.push(array); err != nil {
				return err
			}

		case code.OpDict:
			// get operand
			numElem := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			dict, err := vm.buildDict(vm.sp-numElem, vm.sp)
			if err != nil {
				return err
			}

			// reduce stack pointer now elems are popped onto hash
			vm.sp = vm.sp - numElem

			if err := vm.push(dict); err != nil {
				return err
			}

		}
	}

	return nil
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value

	case *object.Null:
		return false

	default:
		return true
	}
}

// push object onto stack
func (vm *Vm) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *Vm) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

// Returns item popped from stack last
func (vm *Vm) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *Vm) execBinaryOp(op code.Opcode) error {
	right := vm.pop()
	rightType := right.Type()
	left := vm.pop()
	leftType := left.Type()

	switch {
	case leftType == object.IntegerObj && rightType == object.IntegerObj:
		return vm.execBinaryIntOp(op, left.(*object.Integer), right.(*object.Integer))

	case leftType == object.StringObj && rightType == object.StringObj:
		return vm.execBinaryStringOp(op, left.(*object.String), right.(*object.String))
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *Vm) execBinaryStringOp(op code.Opcode, left *object.String, right *object.String) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	return vm.push(&object.String{Value: left.Value + right.Value})
}

func (vm *Vm) execBinaryIntOp(op code.Opcode, left *object.Integer, right *object.Integer) error {
	var err error = nil
	switch op {
	case code.OpAdd:
		err = vm.push(&object.Integer{Value: left.Value + right.Value})

	case code.OpSub:
		err = vm.push(&object.Integer{Value: left.Value - right.Value})

	case code.OpMul:
		err = vm.push(&object.Integer{Value: left.Value * right.Value})

	case code.OpDiv:
		err = vm.push(&object.Integer{Value: left.Value / right.Value})
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return err
}

func (vm *Vm) execComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.IntegerObj || right.Type() == object.IntegerObj {
		return vm.execIntComparison(op, left, right)
	}

	switch op {
	case code.OpEq:
		return vm.push(boolToObject(left == right))
	case code.OpNe:
		return vm.push(boolToObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *Vm) execIntComparison(op code.Opcode, left object.Object, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case code.OpEq:
		return vm.push(boolToObject(leftVal == rightVal))
	case code.OpNe:
		return vm.push(boolToObject(leftVal != rightVal))
	case code.OpGt: // less than is converted to Gt in compiler
		return vm.push(boolToObject(leftVal > rightVal))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *Vm) execNegation() error {
	right := vm.pop()

	if right.Type() != object.IntegerObj {
		return fmt.Errorf("illegal operator - on type %s", right.Type())
	}

	if err := vm.push(&object.Integer{Value: -right.(*object.Integer).Value}); err != nil {
		return err
	}

	return nil
}

func (vm *Vm) execBoolNegation() error {
	right := vm.pop()

	switch right.(type) {
	case *object.Integer:
		rightVal := right.(*object.Integer).Value
		if rightVal == 0 {
			if err := vm.push(True); err != nil {
				return err
			}
		} else {
			if err := vm.push(False); err != nil {
				return err
			}
		}

	case *object.Boolean:
		switch right {
		case True:
			return vm.push(False)
		case False:
			return vm.push(True)
		default:
			return vm.push(False)
		}

	case *object.Null:
		return vm.push(True)

	default:
		return fmt.Errorf("illegal operator ! for type: %s", right.Type())
	}
	return nil
}

func (vm *Vm) buildArray(start int, end int) object.Object {
	elem := make([]object.Object, end-start)

	for i := start; i < end; i++ {
		elem[i-start] = vm.stack[i]
	}

	return &object.Array{Elements: elem}
}

// returns Dict, or error if the key is not hashable
func (vm *Vm) buildDict(startIndex int, endIndex int) (object.Object, error) {
	dictPairs := make(map[object.DictKey]object.DictPair)

	// key + val
	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		val := vm.stack[i+1]

		pair := object.DictPair{Key: key, Value: val}

		dictKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unhashable key %s", key.Type())
		}
		dictPairs[dictKey.DictKey()] = pair
	}

	return &object.Dict{Pairs: dictPairs}, nil
}
