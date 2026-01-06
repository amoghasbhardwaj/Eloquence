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
	// if x equals 10 show x end
	programWords := []string{
		"x", "is", "10",
		"if", "x", "equals", "10",
		"show", "x",
		"end",
	}

	// Expected types for the sequence above
	// Note: "10" is conceptually an INT, but LookupIdent treats anything not in the map as IDENT.
	// The Lexer handles INT/FLOAT logic distinct from LookupIdent.
	// So "10" here results in IDENT from *LookupIdent*, which is correct behavior for this specific function.
	expectedTypes := []TokenType{
		IDENT, IS, IDENT,
		IF, IDENT, EQUALS, IDENT,
		SHOW, IDENT,
		END,
	}

	for i, word := range programWords {
		got := LookupIdent(word)
		if got != expectedTypes[i] {
			t.Errorf("FAIL: Word index %d (%q). Got %q, expected %q", i, word, got, expectedTypes[i])
		}
	}
}
