package vm

import (
	"fmt"
	"testing"

	"crabscript.rs/ast"
	"crabscript.rs/compiler"
	"crabscript.rs/lexer"
	"crabscript.rs/object"
	"crabscript.rs/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"2 * 2", 4},
		{"8 / 2", 4},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!(if (false) { 5; })", true},
	}
	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}
	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}
	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

func TestDictLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.DictKey]int64{},
		},
		{
			"{1: 2, 2: 3}", map[object.DictKey]int64{
				(&object.Integer{Value: 1}).DictKey(): 2,
				(&object.Integer{Value: 2}).DictKey(): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.DictKey]int64{
				(&object.Integer{Value: 2}).DictKey(): 4,
				(&object.Integer{Value: 6}).DictKey(): 16,
			},
		},
	}

	runVmTests(t, tests)
}

func TestIndexLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVmTests(t, tests)
}

func TestCallFnNoArgs(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
let ninePlusTen = fn() { 9 + 10; };
ninePlusTen();
`,
			expected: 19,
		},
		{
			input: `
let one = fn() { 1; };
let two = fn() { 2; };
one() + two()
`,
			expected: 3,
		},
		{
			input: `
let a = fn() { 1 };
let b = fn() { a() + 1 };
let c = fn() { b() + 1 };
c();
`,
			expected: 3,
		},
	}

	runVmTests(t, tests)
}

func TestFnWithReturn(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
let earlyExit = fn() { return 69; 420; }
earlyExit();
`,
			expected: 69,
		},
		{

			input: `
let earlyExit = fn() { return 69; return 420; }
earlyExit();
`,
			expected: 69,
		},
	}

	runVmTests(t, tests)
}

func TestFnNoReturnVal(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
let noReturn = fn() {};
noReturn();
`,
			expected: Null,
		},
		{
			input: `
let noReturn = fn() {};
let noReturnTwo = fn() { noReturn(); };
noReturn();
noReturnTwo();
`,
			expected: Null,
		},
	}
	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %S", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObj(t, tt.expected, stackElem)
	}
}

func TestAnonymousFns(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
let returnsOne = fn() { 1; };
let uselessFactory = fn() { returnsOne; };
uselessFactory()();
    `,
			expected: 1,
		},
	}
	runVmTests(t, tests)
}

func TestCallFnWithBinding(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let one = fn() { let one = 1; one };
		one();
		`,
			expected: 1,
		},
		{
			input: `
		let oneAndTwo = fn() { let one = 1; let two = 2; one + two };
		oneAndTwo();
		`,
			expected: 3,
		},
		{
			input: `
		let oneAndTwo = fn() { let one = 1; let two = 2; one + two};
		let threeAndFour = fn() { let three = 3; let four = 4; three + four};
		oneAndTwo() + threeAndFour();
		`,
			expected: 10,
		},
		{
			input: `
		let fooOne = fn() { let foo = 5; foo; };
		let fooTwo = fn() { let foo = 6; foo; };
		fooOne() + fooTwo();
		`,
			expected: 11,
		},
		{
			input: `
		let globalVar = 69;
		let subOne = fn() {
			let num = 1;
			globalVar - 1
		}
		
		let subTwo = fn() {
			let num = 2;
			globalVar - 2
		}
		
		subOne() + subTwo();
		`,
			expected: 135,
		},
	}
	runVmTests(t, tests)
}

func testExpectedObj(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}

	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}

	case []int:
		arr, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object is not of the type Array, got %T (%+v)", arr, arr.Type())
		}

		if len(arr.Elements) != len(expected) {
			t.Errorf("Unexpected length, want %d, got %d", len(expected), len(arr.Elements))
		}

		for i, e := range expected {
			if err := testIntegerObject(int64(e), arr.Elements[i]); err != nil {
				t.Errorf("mismatching element (%v) at position %d, want %v", arr.Elements[i].(*object.Integer).Value, i, e)
			}
		}
	case map[object.DictKey]int64:
		hash, ok := actual.(*object.Dict)
		if !ok {
			t.Errorf("object is not a dict, got %T (%+v)", actual, actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("dict has wrong number of elements, got %d want %d", len(hash.Pairs), len(expected))
			return
		}

		for eKey, eVal := range expected {
			pair, ok := hash.Pairs[eKey]
			if !ok {
				t.Errorf("key %v does not exist in dict", eKey)
			}

			if err := testIntegerObject(eVal, pair.Value); err != nil {
				t.Errorf("value got same, got %v want %v", pair.Value, eVal)
			}
		}
	}
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not a String. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got (%q) want (%q)", result.Value, expected)
	}
	return nil
}

func testBooleanObject(expected bool, actual object.Object) interface{} {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
	}
	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)

	return p.ParseProgram()
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
