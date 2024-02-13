package code

import "encoding/binary"

// generates bytecode from opcode
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// create buffer for operation
	instLen := 1
	for _, w := range def.OperandWidths {
		instLen += w
	}
	instruction := make([]byte, instLen)
	instruction[0] = byte(op)

	// add operands to the op buffer
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		// add found op
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))

		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}
	return instruction
}
