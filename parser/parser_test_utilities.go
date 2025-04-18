package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

// testIdentifier handles basic testing of an identifier expression.
//
// Parameters:
//   - t: The testing instance.
//   - exp: The expression that is thought to be an Identifier.
//   - value: The expected value of the expression.
//
// Returns:
//   - bool: True on success, otherwise false.
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not '%s'. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

// testLiteralExpression performs tests on literal expressions.
//
// Parameters:
//   - t: The testing instance.
//   - exp: The expression that is to be tested.
//   - expected: The expected value of the expression.
//
// Returns:
//   - bool: True when success, otherwise false.
func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// testInfixExpression handles testing infix expressions.
//
// Parameters:
//   - t: The testing instance.
//   - exp: The expression to test.
//   - left: The expected value of the left hand side.
//   - operator: The expected operator.
//   - right: The expected value of the right hand side.
//
// Returns:
//   - bool: True when success, otherwise false.
func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

// testIntegerLiteral tests to make sure that the expression passed is an integer literal type
// and that it matches the expected value.
//
// Parameters:
//   - t: The testing instance.
//   - il: The integer literal expression.
//   - value: The expected value.
//
// Returns:
//   - bool: True when the expression is an integer literal and its value matches the input value, otherwise false.
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

// checkParserErrors evaluates any parser errors that occurred during parsing and logs them as test failures.
//
// Parameters:
//   - t: The testing instance.
//   - p: The parser to check for errors.
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

// testLetStatement tests a Let Statement for the possible failures.
//
// Parameters:
//   - t: The testing instance.
//   - s: The statement to test.
//   - name: The expected name of the bound identifier.
//
// Returns:
//   - bool: True on success, otherwise false.
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. got=%T", s)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
	}

	return true
}

// constructTestProgram creates a lexer, a parser, and parses a Program node from the given input.
//
// Parameters:
//   - t: The testing instance.
//   - input: The source code input string to use to generate the program.
//
// Returns:
//   - *ast.Program: The program generated by the parser.
func constructTestProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	return program
}

// testBooleanLiteral tests a boolean expression.
//
// Parameters:
//   - t: The testing instance.
//   - exp: The expression thought to be a boolean.
//   - value: The expected value of the expression.
//
// Returns:
//   - bool: True when success, otherwise false.
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value is not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral() not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}
