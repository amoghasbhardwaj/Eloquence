// ==============================================================================================
// FILE: ast/ast_integration_test.go
// ==============================================================================================
// PURPOSE: Integration tests for AST nodes.
//          Verifies that complex, nested structures (like functions and structs)
//          are assembled and stringified correctly.
// ==============================================================================================

package ast

import (
	"testing"

	"eloquence/token"
)

// TestFunctionAndCallIntegration verifies the structure of a function definition
// combined with a function call.
func TestFunctionAndCallIntegration(t *testing.T) {
	// Construct: takes (x) { return x }
	fn := &FunctionLiteral{
		Token:      token.Token{Type: token.TAKES, Literal: "takes"},
		Parameters: []*Identifier{{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"}},
		Body: &BlockStatement{
			Token:      token.Token{Type: token.LBRACE, Literal: "{"},
			Statements: []Statement{&ReturnStatement{Token: token.Token{Type: token.RETURN, Literal: "return"}, ReturnValue: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"}}},
		},
	}

	// Construct: <func>(5)
	call := &CallExpression{
		Token:     token.Token{Type: token.LPAREN, Literal: "("},
		Function:  fn,
		Arguments: []Expression{&IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5}},
	}

	expectedCall := "takes (x) return x(5)"
	if call.String() != expectedCall {
		t.Fatalf("expected %s, got %s", expectedCall, call.String())
	}
}

// TestProgramStringIntegration verifies that a Program node correctly concatenates
// multiple statements into a coherent source string.
func TestProgramStringIntegration(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&AssignmentStatement{
				Token: token.Token{Type: token.IS, Literal: "is"},
				Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
				Value: &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "10"}, Value: 10},
			},
			&ShowStatement{
				Token: token.Token{Type: token.SHOW, Literal: "show"},
				Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
			},
		},
	}

	// Eloquence statements don't force newlines in .String(), they just concatenate.
	expected := "x is 10show x"
	if prog.String() != expected {
		t.Fatalf("expected %s, got %s", expected, prog.String())
	}
}

// TestStructAndPointers verifies the AST representation for struct definitions.
func TestStructAndPointers(t *testing.T) {
	// Construct: define Node as struct { v }
	structDef := &StructDefinitionStatement{
		Token:      token.Token{Type: token.DEFINE, Literal: "define"},
		Name:       &Identifier{Token: token.Token{Type: token.IDENT, Literal: "Node"}, Value: "Node"},
		Attributes: []*Identifier{{Token: token.Token{Type: token.IDENT, Literal: "v"}, Value: "v"}},
	}

	expectedDef := "define Node as struct { v }"
	if structDef.String() != expectedDef {
		t.Fatalf("expected %s, got %s", expectedDef, structDef.String())
	}
}
