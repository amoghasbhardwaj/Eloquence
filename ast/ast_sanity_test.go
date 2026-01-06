// ==============================================================================================
// FILE: ast/ast_sanity_test.go
// ==============================================================================================
// PURPOSE: Sanity checks for the AST package.
//          Tests extreme cases like empty programs and deep nesting to ensure
//          no panics or stack overflows occur during stringification.
// ==============================================================================================

package ast

import (
	"testing"

	"eloquence/token"
)

// TestDeeplyNestedExpressions creates a highly recursive expression
// (not not not ... 1) to ensure the AST doesn't crash on deep traversal.
func TestDeeplyNestedExpressions(t *testing.T) {
	depth := 100
	var expr Expression = &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1}

	// Wrap it 100 times
	for i := 0; i < depth; i++ {
		expr = &PrefixExpression{
			Token:    token.Token{Type: token.NOT, Literal: "not"},
			Operator: "not",
			Right:    expr,
		}
	}

	// Just ensure it generates *something* without panicking
	if expr.String() == "" {
		t.Fatal("nested expression produced empty string")
	}
}

// TestEmptyProgramSanity verifies that an empty AST produces an empty string
// rather than a nil pointer dereference.
func TestEmptyProgramSanity(t *testing.T) {
	prog := &Program{Statements: []Statement{}}
	if prog.String() != "" {
		t.Fatalf("expected empty string for empty program, got %s", prog.String())
	}
}
