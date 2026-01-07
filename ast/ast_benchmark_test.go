// ==============================================================================================
// FILE: ast/ast_benchmark_test.go
// ==============================================================================================
// PURPOSE: Performance benchmarks for the Abstract Syntax Tree (AST).
//          These tests measure the efficiency of the .String() methods, which involves
//          recursive tree traversal and string concatenation.
//          High performance here is important for logging, debugging, and code formatting tools.
// ==============================================================================================

package ast

import (
	"testing"

	"eloquence/token"
)

// BenchmarkInfixExpressionString measures the allocation and speed cost of
// converting a binary expression (e.g., "100 adds 200") back to its string representation.
// Usage: go test -bench=BenchmarkInfixExpressionString ./ast
func BenchmarkInfixExpressionString(b *testing.B) {
	// Setup a static expression tree: (100 adds 200)
	left := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "100"}, Value: 100}
	right := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "200"}, Value: 200}
	expr := &InfixExpression{
		Token:    token.Token{Type: token.ADDS, Literal: "adds"},
		Left:     left,
		Operator: "adds",
		Right:    right,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We discard the result, we only care about the CPU cycles used.
		_ = expr.String()
	}
}

// BenchmarkLargeProgramString measures the performance of the root Program node
// when iterating over a large slice of statements. This simulates the overhead
// of printing a moderately sized source file.
// Usage: go test -bench=BenchmarkLargeProgramString ./ast
func BenchmarkLargeProgramString(b *testing.B) {
	// Construct a program with 1000 "show(1)" statements
	count := 1000
	prog := &Program{Statements: make([]Statement, count)}

	// Create a reusable statement: show(1)
	// This corresponds to an ExpressionStatement containing a CallExpression
	stmt := &ExpressionStatement{
		Token: token.Token{Type: token.IDENT, Literal: "show"},
		Expression: &CallExpression{
			Token: token.Token{Type: token.LPAREN, Literal: "("},
			Function: &Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "show"},
				Value: "show",
			},
			Arguments: []Expression{
				&IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
			},
		},
	}

	// Fill the program
	for i := 0; i < count; i++ {
		prog.Statements[i] = stmt
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = prog.String()
	}
}
