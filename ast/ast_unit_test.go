// ==============================================================================================
// FILE: ast/ast_unit_test.go
// ==============================================================================================
// PURPOSE: Unit tests for individual AST nodes.
//          Verifies that literals and statements stringify themselves correctly.
// ==============================================================================================

package ast

import (
	"testing"

	"eloquence/token"
)

// ----------------------------------------------------------------------------
// LITERALS
// ----------------------------------------------------------------------------

func TestIntegerLiteral(t *testing.T) {
	node := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "42"}, Value: 42}
	if node.String() != "42" {
		t.Fatalf("expected 42, got %s", node.String())
	}
}

func TestStringLiteral(t *testing.T) {
	node := &StringLiteral{Token: token.Token{Type: token.STRING, Literal: "hello"}, Value: "hello"}
	// String() must wrap the value in quotes to represent source code
	expected := `"hello"`
	if node.String() != expected {
		t.Fatalf("expected %s, got %s", expected, node.String())
	}
}

// ----------------------------------------------------------------------------
// EXPRESSIONS
// ----------------------------------------------------------------------------

func TestInfixExpression(t *testing.T) {
	// Testing: 5 adds 3
	node := &InfixExpression{
		Token:    token.Token{Type: token.ADDS, Literal: "adds"},
		Left:     &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
		Operator: "adds",
		Right:    &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "3"}, Value: 3},
	}
	expected := "(5 adds 3)"
	if node.String() != expected {
		t.Fatalf("expected %s, got %s", expected, node.String())
	}
}

func TestArrayLiteral(t *testing.T) {
	// Testing: [1, 2]
	node := &ArrayLiteral{
		Token: token.Token{Type: token.LBRACKET, Literal: "["},
		Elements: []Expression{
			&IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
			&IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
		},
	}
	expected := "[1, 2]"
	if node.String() != expected {
		t.Fatalf("expected %s, got %s", expected, node.String())
	}
}

// ----------------------------------------------------------------------------
// STATEMENTS
// ----------------------------------------------------------------------------

func TestAssignmentStatement(t *testing.T) {
	// Testing: x is 5
	node := &AssignmentStatement{
		Token: token.Token{Type: token.IDENT, Literal: "x"},
		Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
		Value: &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
	}
	expected := "x is 5"
	if node.String() != expected {
		t.Fatalf("expected %s, got %s", expected, node.String())
	}
}

func TestReturnStatement(t *testing.T) {
	// Testing: return 10
	node := &ReturnStatement{
		Token:       token.Token{Type: token.RETURN, Literal: "return"},
		ReturnValue: &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "10"}, Value: 10},
	}
	expected := "return 10"
	if node.String() != expected {
		t.Fatalf("expected %s, got %s", expected, node.String())
	}
}

// Replaced TestShowStatement with TestCallExpression since 'show' is now just a function call
func TestCallExpression(t *testing.T) {
	// Testing: show(msg)
	node := &CallExpression{
		Token:    token.Token{Type: token.LPAREN, Literal: "("},
		Function: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "show"}, Value: "show"},
		Arguments: []Expression{
			&Identifier{Token: token.Token{Type: token.IDENT, Literal: "msg"}, Value: "msg"},
		},
	}

	expected := "show(msg)"
	if node.String() != expected {
		t.Fatalf("expected %s, got %s", expected, node.String())
	}
}
