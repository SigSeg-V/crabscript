package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{
			OpConstant,
			[]int{65534},
			[]byte{byte(OpConstant), 255, 254},
		},
		{
			OpAdd,
			[]int{},
			[]byte{byte(OpAdd)},
		},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction is wrong length, got %d want %d", len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != b {
				t.Errorf("wrong byte at pos %d, want %d got %d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 1),
		Make(OpConstant, 65534),
	}

	expected := `0000 OpAdd
0001 OpConstant 1
0004 OpConstant 65534
`
	concatted := Instructions{}
	for _, inst := range instructions {
		concatted = append(concatted, inst...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions are formatted incorrectly.\nwant %q\ngot %q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}
