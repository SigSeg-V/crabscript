package compiler

import (
	"fmt"
	"testing"

	"crabscript.rs/ast"
	"crabscript.rs/code"
	"crabscript.rs/lexer"
	"crabscript.rs/object"
	"crabscript.rs/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1; 2", // 2 distinct expressions
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConst, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpNeg),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestBooleanExpresions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 > 2", expectedConstants: []interface{}{1, 2}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpGt),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 < 2", expectedConstants: []interface{}{2, 1}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpGt),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 == 2", expectedConstants: []interface{}{1, 2}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpEq),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 != 2", expectedConstants: []interface{}{1, 2}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpNe),
				code.Make(code.OpPop),
			},
		},
		{
			input: "true == false", expectedConstants: []interface{}{}, expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEq),
				code.Make(code.OpPop),
			},
		},
		{
			input: "true != false", expectedConstants: []interface{}{}, expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNe),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestIndexExperssions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpConst, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConst, 3),
				code.Make(code.OpConst, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIdx),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []interface{}{1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpDict, 2),
				code.Make(code.OpConst, 2),
				code.Make(code.OpConst, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIdx),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestLetStatementScope(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
  let num = 69;
  fn() { num };
  `,
			expectedConstants: []interface{}{69,
				[]code.Instructions{
					code.Make(code.OpGetGbl, 0),
					code.Make(code.OpRetVal),
				},
			},

			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpSetGbl, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
      fn() {
        let num = 69;
        num
      };
      `,
			expectedConstants: []interface{}{69,
				[]code.Instructions{
					code.Make(code.OpConst, 0),
					code.Make(code.OpSetLcl, 0),
					code.Make(code.OpGetLcl, 0),
					code.Make(code.OpRetVal),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
        fn() {
          let a = 55;
          let b = 77;
          a + b
        };
        `,
			expectedConstants: []interface{}{
				55,
				77,
				[]code.Instructions{
					code.Make(code.OpConst, 0),
					code.Make(code.OpSetLcl, 0),
					code.Make(code.OpConst, 1),
					code.Make(code.OpSetLcl, 1),
					code.Make(code.OpGetLcl, 0),
					code.Make(code.OpGetLcl, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpRetVal),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 2),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)

		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got %d, want %d",
			len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}

		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}

		case []code.Instructions: // fn block of instructions
			fn, ok := actual[i].(*object.CompFn)
			if !ok {
				return fmt.Errorf("constant %d - not a fn, type %T", i, actual[i])
			}
			if err := testInstructions(constant, fn.Instructions); err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s", i, err)
			}
		}
	}
	return nil
}

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex incorrect, got %d want %d", compiler.scopeIndex, 0)
	}
	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex incorrect, got %d want %d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.OpSub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong, got %d want %d", len(compiler.scopes[compiler.scopeIndex].instructions), 1)
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpSub {
		t.Errorf("lastinstruction wrong, got %d want %d", last.Opcode, code.OpSub)
	}

	// make sure we enclose scope correctly
	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose scope")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex incorrect, got %d want %d", compiler.scopeIndex, 0)
	}

	// make sure we enclose scope correctly
	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not leave scope")
	}

	compiler.emit(code.OpAdd)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong, got %d want %d", len(compiler.scopes[compiler.scopeIndex].instructions), 2)
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpAdd {
		t.Errorf("lastinstruction wrong, got %d want %d", last.Opcode, code.OpSub)
	}

	prev := compiler.scopes[compiler.scopeIndex].previousInstruction
	if prev.Opcode != code.OpMul {
		t.Errorf("previnstruction wrong, got %d want %d", last.Opcode, code.OpSub)
	}
}

func testStringObject(expected string, o object.Object) interface{} {
	result, ok := o.(*object.String)
	if !ok {
		return fmt.Errorf("object not String, got %T (%+v)", o, o)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got (%q) want (%q)", result.Value, expected)
	}
	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got %T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %d, want %d", result.Value, expected)
	}
	return nil
}

func testInstructions(
	expected []code.Instructions, actual code.Instructions,
) error {
	concatted := concatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant %q\ngot %q",
			concatted, actual)
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q", i, concatted, actual)
		}
	}
	return nil
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `if (true) { 10 } else { 20 }; 3333;`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0
				code.Make(code.OpTrue),
				// 1
				code.Make(code.OpJmpNt, 10),
				// 4
				code.Make(code.OpConst, 0),
				// 7
				code.Make(code.OpJmp, 13),
				// 10
				code.Make(code.OpConst, 1),
				// 13
				code.Make(code.OpPop),
				// 14
				code.Make(code.OpConst, 2),
				// 17
				code.Make(code.OpPop),
			},
		},
		{
			input: `
if (true) { 10 }; 3333;
`,
			expectedConstants: []interface{}{10, 3333}, expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJmpNt, 10),
				// 0004
				code.Make(code.OpConst, 0),
				// 0007
				code.Make(code.OpJmp, 11),
				// 0010
				code.Make(code.OpNull),
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConst, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
let one = 1;
let two = 2;
`,
			expectedConstants: []interface{}{1, 2}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpSetGbl, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpSetGbl, 1),
			}},
		{
			input: `
let one = 1;
one;
`,
			expectedConstants: []interface{}{1}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpSetGbl, 0),
				code.Make(code.OpGetGbl, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
let one = 1;
let two = one;
two;
`,
			expectedConstants: []interface{}{1}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpSetGbl, 0),
				code.Make(code.OpGetGbl, 0),
				code.Make(code.OpSetGbl, 1),
				code.Make(code.OpGetGbl, 1),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `"monkey"`, expectedConstants: []interface{}{"monkey"}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `"mon" + "key"`, expectedConstants: []interface{}{"mon", "key"}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "[]", expectedConstants: []interface{}{}, expectedInstructions: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: "[1, 2, 3]", expectedConstants: []interface{}{1, 2, 3}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpConst, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input: "[1 + 2, 3 - 4, 5 * 6]", expectedConstants: []interface{}{1, 2, 3, 4, 5, 6}, expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConst, 2),
				code.Make(code.OpConst, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConst, 4),
				code.Make(code.OpConst, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestDictLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "{}",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpDict, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpConst, 2),
				code.Make(code.OpConst, 3),
				code.Make(code.OpConst, 4),
				code.Make(code.OpConst, 5),
				code.Make(code.OpDict, 6),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpConst, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConst, 3),
				code.Make(code.OpConst, 4),
				code.Make(code.OpConst, 5),
				code.Make(code.OpMul),
				code.Make(code.OpDict, 4),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFns(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn(){ return 5 + 10 }`,
			expectedConstants: []interface{}{5, 10, // definitions in the fn
				[]code.Instructions{ // actual fn as a const
					code.Make(code.OpConst, 0),
					code.Make(code.OpConst, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpRetVal),
				}},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn(){ 5 + 10 }`, // implicit returns also work
			expectedConstants: []interface{}{5, 10, // definitions in the fn
				[]code.Instructions{ // actual fn as a const
					code.Make(code.OpConst, 0),
					code.Make(code.OpConst, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpRetVal),
				}},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn(){ 1; 2 }`, // testing imp return with multiple statements
			expectedConstants: []interface{}{1, 2,
				[]code.Instructions{
					code.Make(code.OpConst, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConst, 1),
					code.Make(code.OpRetVal),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 2),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFnsNoRetVal(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `fn(){}`,
			expectedConstants: []interface{}{[]code.Instructions{code.Make(code.OpRet)}},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestFnCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { 24 }();`, expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConst, 0), // The literal "24"
					code.Make(code.OpRetVal),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 1), // The compiled function code.Make(code.OpCall),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
				let noArg = fn() {24};
				noArg();
		`,
			expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConst, 0), // The literal "24"
					code.Make(code.OpRetVal),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 1),  // The compiled function code.Make(code.OpSetGlobal, 0), code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGbl, 0), // putting fn onto the const pool
				code.Make(code.OpGetGbl, 0), // getting from the const pool
				code.Make(code.OpCall, 0),   // calling fn
				code.Make(code.OpPop),
			},
		},
		{
			input: `
let oneArg = fn(hi) { hi };
oneArg(69);
`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLcl, 0),
					code.Make(code.OpRetVal),
				},
				69,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpSetGbl, 0),
				code.Make(code.OpGetGbl, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
let polyArgs = fn(a,b,c,d) {a; b; c; d };
polyArgs(1, 2, 3, 4);
`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLcl, 0),
					code.Make(code.OpPop),
					code.Make(code.OpGetLcl, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetLcl, 2),
					code.Make(code.OpPop),
					code.Make(code.OpGetLcl, 3),
					code.Make(code.OpRetVal),
				},
				1,
				2,
				3,
				4,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConst, 0),
				code.Make(code.OpSetGbl, 0),
				code.Make(code.OpGetGbl, 0),
				code.Make(code.OpConst, 1),
				code.Make(code.OpConst, 2),
				code.Make(code.OpConst, 3),
				code.Make(code.OpConst, 4),
				code.Make(code.OpCall, 4),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)

	return p.ParseProgram()
}
