// ==============================================================================================
// FILE: parser/parser_unit_test.go
// ==============================================================================================
// PURPOSE: Unit tests for individual parser components.
// ==============================================================================================

package parser

import (
	"testing"

	"eloquence/ast"
	"eloquence/lexer"
)

func newParser(input string) *Parser {
	l := lexer.New(input)
	return New(l)
}

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

func TestCallExpression(t *testing.T) {
	input := `show(x)`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement")
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not *ast.ExpressionStatement")
	}
	call, ok := expStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression")
	}
	if call.Function.String() != "show" {
		t.Errorf("function name not 'show'")
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
		t.Fatalf("expected 3 statements")
	}

	stmtA := program.Statements[0].(*ast.AssignmentStatement)
	prefixA, ok := stmtA.Value.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("stmtA.Value is not PrefixExpression")
	}
	if prefixA.Operator != "-" {
		t.Errorf("operator is not '-'")
	}

	stmtC := program.Statements[2].(*ast.AssignmentStatement)
	if _, ok := stmtC.Value.(*ast.PointerReferenceExpression); !ok {
		t.Errorf("stmtC.Value is not PointerReferenceExpression")
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
			t.Fatalf("stmt is not AssignmentStatement")
		}
		if _, ok := assign.Value.(*ast.InfixExpression); !ok {
			t.Errorf("assign.Value is not InfixExpression")
		}
	}
}

func TestFunctionAndCall(t *testing.T) {
	// FIXED SYNTAX: Added braces around function body
	input := `fn is takes (x, y) {
  return x adds y
}
result is fn(1, 2)`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}

	fnStmt := program.Statements[0].(*ast.AssignmentStatement)
	if _, ok := fnStmt.Value.(*ast.FunctionLiteral); !ok {
		t.Errorf("expected FunctionLiteral")
	}

	callStmt := program.Statements[1].(*ast.AssignmentStatement)
	if _, ok := callStmt.Value.(*ast.CallExpression); !ok {
		t.Errorf("expected CallExpression")
	}
}

func TestIfExpression(t *testing.T) {
	input := `result is if x less y {
  show(x)
} else {
  show(y)
}`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement")
	}
	assign := program.Statements[0].(*ast.AssignmentStatement)
	if _, ok := assign.Value.(*ast.IfExpression); !ok {
		t.Errorf("expected IfExpression")
	}
}

func TestLoopStatements(t *testing.T) {
	input := `for i in list {
  show(i)
}
while flag {
  flag is false
}`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements")
	}
	if _, ok := program.Statements[0].(*ast.RangeLoopStatement); !ok {
		t.Errorf("expected RangeLoopStatement")
	}
	if _, ok := program.Statements[1].(*ast.LoopStatement); !ok {
		t.Errorf("expected LoopStatement")
	}
}

func TestTryCatchFinally(t *testing.T) {
	input := `try {
  x is 5
} catch {
  show("error")
} finally {
  show("done")
}`

	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement")
	}
	if _, ok := program.Statements[0].(*ast.TryCatchStatement); !ok {
		t.Errorf("expected TryCatchStatement")
	}
}

func TestStructDefinition(t *testing.T) {
	input := `define Node as struct { val, next }`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement")
	}
	stmt := program.Statements[0].(*ast.StructDefinitionStatement)
	if stmt.Name.Value != "Node" {
		t.Errorf("expected struct name Node")
	}
}

func TestPointerAssignmentStatement(t *testing.T) {
	input := `pointing from ptr is 10`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement")
	}
	stmt, ok := program.Statements[0].(*ast.PointerAssignmentStatement)
	if !ok {
		t.Fatalf("expected PointerAssignmentStatement")
	}
	if stmt.Name.Value != "ptr" {
		t.Errorf("expected name 'ptr'")
	}
}

func TestStructInstantiation(t *testing.T) {
	input := `user is User { name: "John", age: 25 }`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement")
	}
	assign := program.Statements[0].(*ast.AssignmentStatement)
	structInst, ok := assign.Value.(*ast.StructInstantiationExpression)
	if !ok {
		t.Fatalf("expected StructInstantiationExpression")
	}
	if structInst.Name.Value != "User" {
		t.Errorf("expected struct name 'User'")
	}
	if len(structInst.Fields) != 2 {
		t.Errorf("expected 2 fields")
	}
}

func TestIncludeStatement(t *testing.T) {
	input := `include "math.eq"`
	p := newParser(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement")
	}
	inc, ok := program.Statements[0].(*ast.IncludeStatement)
	if !ok {
		t.Fatalf("expected IncludeStatement")
	}
	if inc.Path.String() != `"math.eq"` {
		t.Errorf("expected path string")
	}
}
