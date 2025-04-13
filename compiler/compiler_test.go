package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

// runCompilerTests is a helper test method to setup and run tests for any compilerTestCase
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

// parse parses the input string into a program node for the compiler.
//
// Parameters:
//   - input: The input source code for the test.
//
// Returns:
//   - *ast.Program: The output program node that will be used to generate bytecode.
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

// testInstructions compares the expected instructions to the actual compiled testInstructions
// to ensure they match.
//
// Parameters:
//   - expected: A list of byte arrays of expected instructions.
//   - actual: The actual compiled testInstructions
//
// Returns:
//   - error: An error if a compilation error is encountered, otherwise nil.
func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(concatted) != len(actual) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d. \nwant=%q\ngot=%q", i, concatted, actual)
		}
	}

	return nil
}

// concatInstructions flattens an array of arrays to a single array containing code instructions.
//
// Parameters:
//   - expected: The expected instructions as separate arrays
//
// Returns:
//   - code.Instructions: A single array of bytes flattened to match the compiler output.
func concatInstructions(expected []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range expected {
		out = append(out, ins...)
	}

	return out
}

// testConstants compares the actual and expected constant values.
func testConstants(t *testing.T, expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

// testIntegerObject ensures the expected value matches the object provided.
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}
