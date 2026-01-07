// ==============================================================================================
// FILE: token/token_sanity_test.go
// ==============================================================================================
// PURPOSE: A high-level check to ensure the token system holds up under a simulated program flow.
//          It mimics the sequence of words a lexer might produce.
// ==============================================================================================

package token

import "testing"

// TestSanityFullProgram simulates a small Eloquence program broken into words
// and verifies that looking them up doesn't cause panics or unexpected behavior.
func TestSanityFullProgram(t *testing.T) {
	// Program representation:
	// x is 10
	// if x equals 10 { show(x) }
	programWords := []string{
		"x", "is", "10",
		"if", "x", "equals", "10", "{",
		"show", "(", "x", ")",
		"}",
	}

	// Expected types for the sequence above
	expectedTypes := []TokenType{
		IDENT, IS, IDENT, // 10 is IDENT via LookupIdent (lexer handles numbers separately)
		IF, IDENT, EQUALS, IDENT, LBRACE,
		IDENT, LPAREN, IDENT, RPAREN, // 'show' is now an IDENT
		RBRACE,
	}

	for i, word := range programWords {
		got := LookupIdent(word)
		// Special handling: punctuation (braces/parens) usually skips LookupIdent in Lexer,
		// but if passed here, they default to IDENT.
		// For this sanity test, we assume direct token mapping for keywords, IDENT for others.

		expected := expectedTypes[i]

		// Map special chars manually for this specific test structure if needed,
		// or rely on the fact that LookupIdent returns IDENT for symbols not in keyword map.
		if word == "{" || word == "}" || word == "(" || word == ")" {
			if got != IDENT {
				t.Errorf("Symbols should be IDENT in raw lookup. Got %q", got)
			}
			continue
		}

		if got != expected {
			t.Errorf("FAIL: Word index %d (%q). Got %q, expected %q", i, word, got, expected)
		}
	}
}
