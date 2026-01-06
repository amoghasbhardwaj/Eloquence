// ==============================================================================================
// FILE: token/token_unit_test.go
// ==============================================================================================
// PURPOSE: Validates the core logic of token mapping. Ensures that every reserved keyword
//          resolves to the correct internal constant.
// ==============================================================================================

package token

import (
	"testing"
)

// TestTokenConstants verifies that the LookupIdent function correctly maps
// string literals to their respective TokenType constants.
// This is critical because if this mapping fails, the Parser will misinterpret keywords as variables.
func TestTokenConstants(t *testing.T) {
	// Table-driven test setup
	tests := []struct {
		word     string    // Input string
		expected TokenType // Expected internal constant
	}{
		// 1. Check Logic Operators
		{"is", IS},
		{"adds", ADDS},
		{"subtracts", SUBTRACTS},

		// 2. Check Comparison Operators
		{"equals", EQUALS},
		{"greater", GREATER},
		{"less_equal", LESS_EQUAL},

		// 3. Check Control Flow
		{"if", IF},
		{"return", RETURN},
		{"end", END},

		// 4. Check Functions & IO
		{"takes", TAKES},
		{"show", SHOW},

		// 5. Check Exception Handling
		{"try", TRY},
		{"catch", CATCH},

		// 6. Check Literals
		{"true", BOOL},
		{"false", BOOL},
		{"none", NIL},

		// 7. Check Structs
		{"define", DEFINE},
		{"struct", STRUCT},

		// 8. Check Non-Keywords (Standard Identifiers)
		{"myVariable", IDENT},
		{"calculateSum", IDENT},
		{"x", IDENT},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			got := LookupIdent(tt.word)
			if got != tt.expected {
				t.Errorf("FAIL: LookupIdent(%q) returned %q, expected %q", tt.word, got, tt.expected)
			}
		})
	}
}

// TestTokenStructStructure verifies that the Token struct is defined correctly.
// While simple, this ensures no regressions occur if fields are renamed later.
func TestTokenStructStructure(t *testing.T) {
	// Create a dummy token
	tok := Token{
		Type:    IS,
		Literal: "is",
		Line:    1,
		Column:  5,
	}

	// Verify fields
	if tok.Type != IS {
		t.Errorf("FAIL: Token.Type mismatch. Got %q, want %q", tok.Type, IS)
	}
	if tok.Literal != "is" {
		t.Errorf("FAIL: Token.Literal mismatch. Got %q, want %q", tok.Literal, "is")
	}
	if tok.Line != 1 {
		t.Errorf("FAIL: Token.Line mismatch. Got %d, want 1", tok.Line)
	}
	if tok.Column != 5 {
		t.Errorf("FAIL: Token.Column mismatch. Got %d, want 5", tok.Column)
	}
}
