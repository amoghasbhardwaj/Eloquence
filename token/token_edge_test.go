// ==============================================================================================
// FILE: token/token_edge_test.go
// ==============================================================================================
// PURPOSE: Tests boundary conditions and unusual inputs to ensure the Token system is robust against
//          malformed or unexpected strings.
// ==============================================================================================

package token

import "testing"

// TestLookupIdentEdgeCases checks empty strings, case sensitivity, and multi-word handling.
func TestLookupIdentEdgeCases(t *testing.T) {
	tests := []struct {
		input string
		want  TokenType
	}{
		// Edge Case 1: Empty String
		// Should default to IDENT, though the lexer usually catches this before calling LookupIdent.
		{"", IDENT},

		// Edge Case 2: Numeric identifiers
		// "123abc" is typically handled by the lexer, but if passed to Lookup, it should be an IDENT.
		{"123abc", IDENT},

		// Edge Case 3: Case Sensitivity
		// Eloquence is case-sensitive. "TRUE" is an identifier, "true" is a boolean literal.
		{"TRUE", IDENT},
		{"If", IDENT},
		{"Include", IDENT},

		// Edge Case 4: Multi-word keywords
		// Ensure that if the lexer successfully groups "pointing from", it maps correctly.
		{"pointing from", POINTING_FROM},
		{"pointing to", POINTING_TO},

		// Edge Case 5: Partial matches for multi-word
		// "pointing" alone or "pointing To" (wrong case) should be identifiers.
		{"pointing", IDENT},
		{"pointing To", IDENT},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := LookupIdent(tt.input)
			if got != tt.want {
				t.Errorf("FAIL: LookupIdent(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}
