// ==============================================================================================
// FILE: parser/parser_unit_test.go
// ==============================================================================================
// PURPOSE: Unit tests for individual parser components.
//          Verifies that specific grammar rules (assignments, math, logic) are parsed
//          correctly into isolated AST nodes.
// ==============================================================================================

package parser

import (
	"testing"

	"eloquence/ast"
	"eloquence/lexer"
)

// Helper: Initializes a parser from an input string.
func newParser(input string) *Parser {
	l := lexer.New(input)
	return New(l)
}

// Helper: Fails the test if the parser encountered errors.
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestAssignmentStatements(t *testing.T) {
	input := `x is 5
y is 10
flag is true
pi is 3.14
name is "Amogh"`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 5 {
		t.Fatalf("expected 5 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedName string
	}{
		{"x"}, {"y"}, {"flag"}, {"pi"}, {"name"},
	}

	for i, stmt := range program.Statements {
		assignStmt, ok := stmt.(*ast.AssignmentStatement)
		if !ok {
			t.Fatalf("test[%d] - statement is not *ast.AssignmentStatement. got=%T", i, stmt)
		}
		if assignStmt.Name.Value != tests[i].expectedName {
			t.Errorf("test[%d] - expected name %s, got %s", i, tests[i].expectedName, assignStmt.Name.Value)
		}
	}
}

func TestShowStatement(t *testing.T) {
	input := `show x`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	showStmt, ok := program.Statements[0].(*ast.ShowStatement)
	if !ok {
		t.Fatalf("statement is not *ast.ShowStatement. got=%T", program.Statements[0])
	}
	if showStmt.Value.String() != "x" {
		t.Errorf("showStmt.Value.String() not 'x'. got=%s", showStmt.Value.String())
	}
}

func TestPrefixExpressions(t *testing.T) {
	input := `a is -5
b is !true
c is pointing to x`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	// Test case 'a is -5'
	stmtA := program.Statements[0].(*ast.AssignmentStatement)
	prefixA, ok := stmtA.Value.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("stmtA.Value is not PrefixExpression. got=%T", stmtA.Value)
	}
	if prefixA.Operator != "-" {
		t.Errorf("operator is not '-'. got=%s", prefixA.Operator)
	}

	// Test case 'c is pointing to x'
	stmtC := program.Statements[2].(*ast.AssignmentStatement)
	if _, ok := stmtC.Value.(*ast.PointerReferenceExpression); !ok {
		t.Errorf("stmtC.Value is not PointerReferenceExpression. got=%T", stmtC.Value)
	}
}

func TestInfixExpressions(t *testing.T) {
	input := `x is a adds b
y is c less d
z is e equals f`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	for _, stmt := range program.Statements {
		assign, ok := stmt.(*ast.AssignmentStatement)
		if !ok {
			t.Fatalf("stmt is not AssignmentStatement. got=%T", stmt)
		}
		if _, ok := assign.Value.(*ast.InfixExpression); !ok {
			t.Errorf("assign.Value is not InfixExpression. got=%T", assign.Value)
		}
	}
}

func TestFunctionAndCall(t *testing.T) {
	input := `fn is takes (x, y)
  return x adds y
end
result is fn(1, 2)`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}

	fnStmt := program.Statements[0].(*ast.AssignmentStatement)
	if _, ok := fnStmt.Value.(*ast.FunctionLiteral); !ok {
		t.Errorf("expected FunctionLiteral, got=%T", fnStmt.Value)
	}

	callStmt := program.Statements[1].(*ast.AssignmentStatement)
	if _, ok := callStmt.Value.(*ast.CallExpression); !ok {
		t.Errorf("expected CallExpression, got=%T", callStmt.Value)
	}
}

func TestIfExpression(t *testing.T) {
	input := `result is if x less y
  show x
else
  show y
end`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	assign := program.Statements[0].(*ast.AssignmentStatement)
	if _, ok := assign.Value.(*ast.IfExpression); !ok {
		t.Errorf("expected IfExpression, got=%T", assign.Value)
	}
}

func TestLoopStatements(t *testing.T) {
	input := `for i less 10
  show i
end
while flag
  flag is false
end`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.LoopStatement); !ok {
			t.Errorf("expected LoopStatement, got %T", stmt)
		}
	}
}

func TestTryCatchFinally(t *testing.T) {
	input := `try
  x is 5
catch
  show "error"
finally
  show "done"
end`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*ast.TryCatchStatement); !ok {
		t.Errorf("expected TryCatchStatement, got=%T", program.Statements[0])
	}
}

func TestStructDefinition(t *testing.T) {
	input := `define Node as struct { val, next }`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	stmt := program.Statements[0].(*ast.StructDefinitionStatement)
	if stmt.Name.Value != "Node" {
		t.Errorf("expected struct name Node, got %s", stmt.Name.Value)
	}
}

func TestPointerAssignmentStatement(t *testing.T) {
	input := `pointing from ptr is 10`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.PointerAssignmentStatement)
	if !ok {
		t.Fatalf("expected PointerAssignmentStatement, got %T", program.Statements[0])
	}
	if stmt.Name.Value != "ptr" {
		t.Errorf("expected name 'ptr', got %s", stmt.Name.Value)
	}
}

func TestStructInstantiation(t *testing.T) {
	input := `user is User { name: "John", age: 25 }`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	assign := program.Statements[0].(*ast.AssignmentStatement)
	structInst, ok := assign.Value.(*ast.StructInstantiationExpression)
	if !ok {
		t.Fatalf("expected StructInstantiationExpression, got %T", assign.Value)
	}
	if structInst.Name.Value != "User" {
		t.Errorf("expected struct name 'User', got %s", structInst.Name.Value)
	}
	if len(structInst.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(structInst.Fields))
	}
}

func TestFieldAccess(t *testing.T) {
	input := `x is user.name`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	assign := program.Statements[0].(*ast.AssignmentStatement)
	fieldAccess, ok := assign.Value.(*ast.FieldAccessExpression)
	if !ok {
		t.Fatalf("expected FieldAccessExpression, got %T", assign.Value)
	}
	if fieldAccess.Field.Value != "name" {
		t.Errorf("expected field name 'name', got %s", fieldAccess.Field.Value)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x is a adds b times c", "x is (a adds (b times c))"},
		{"x is a times b adds c", "x is ((a times b) adds c)"},
		{"x is -a times b", "x is ((- a) times b)"},
		{"x is !a equals b", "x is ((! a) equals b)"},
	}

	for _, tt := range tests {
		p := newParser(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(program.Statements))
		}
		actual := program.Statements[0].String()
		if actual != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, actual)
		}
	}
}
