package vm

import (
	"crabscript.rs/code"
	"crabscript.rs/compiler"
	"crabscript.rs/object"
	"fmt"
)

const stackSize = 2048

type Vm struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int // stack pointer, always at next free slot at top of stack
}

func New(bytecode *compiler.Bytecode) *Vm {
	return &Vm{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, stackSize),
		sp:           0,
	}
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

		switch op {
		case code.OpConstant:
			constIdx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIdx])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// push object onto stack
func (vm *Vm) push(o object.Object) error {
	if vm.sp >= stackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}
