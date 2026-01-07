// ==============================================================================================
// FILE: parser/parser_integration_test.go
// ==============================================================================================
// PURPOSE: Integration tests for the Parser.
//          Validates the parsing of complete, multi-part logical structures like
//          recursive functions and object-oriented structs.
// ==============================================================================================

package parser

import (
	"testing"

	"eloquence/ast"
	"eloquence/lexer"
)

func TestIntegration_FactorialFunction(t *testing.T) {
	input := `
    factorial is takes (n) {
        if n less_equal 1 {
            return 1
        } else {
            return n times factorial(n minus 1)
        }
    }
    
    result is factorial(5)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}

	// 1. Verify Function Definition
	stmt1, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt1 not AssignmentStatement, got %T", program.Statements[0])
	}
	if stmt1.Name.Value != "factorial" {
		t.Errorf("expected function name 'factorial', got %s", stmt1.Name.Value)
	}

	fnLit, ok := stmt1.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt1 value not FunctionLiteral, got %T", stmt1.Value)
	}
	if len(fnLit.Parameters) != 1 || fnLit.Parameters[0].Value != "n" {
		t.Errorf("expected 1 parameter 'n'")
	}

	// 2. Verify Call
	stmt2, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt2 not AssignmentStatement")
	}
	callExp, ok := stmt2.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt2 value not CallExpression")
	}
	if callExp.Function.String() != "factorial" {
		t.Errorf("expected call to 'factorial', got %s", callExp.Function.String())
	}
}

func TestIntegration_StructsAndLogic(t *testing.T) {
	input := `
    define User as struct { name, age }
    
    u is User { name: "Alice", age: 30 }
    
    if u.age greater 18 {
        show("Adult")
    }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	// 1. Struct Definition
	if _, ok := program.Statements[0].(*ast.StructDefinitionStatement); !ok {
		t.Errorf("expected StructDefinitionStatement at 0")
	}

	// 2. Struct Instantiation
	assign, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		t.Errorf("expected AssignmentStatement at 1")
	}
	if _, ok := assign.Value.(*ast.StructInstantiationExpression); !ok {
		t.Errorf("expected StructInstantiationExpression value")
	}

	// 3. Logic with Field Access
	ifStmt, ok := program.Statements[2].(*ast.ExpressionStatement).Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("expected IfExpression at 2, got %T", program.Statements[2])
	}

	// Check condition: (u.age greater 18)
	infix, ok := ifStmt.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("condition not infix")
	}
	if infix.Operator != "greater" {
		t.Errorf("expected operator greater")
	}
	_, isFieldAccess := infix.Left.(*ast.FieldAccessExpression)
	if !isFieldAccess {
		t.Errorf("left side of condition expected FieldAccessExpression")
	}
}
