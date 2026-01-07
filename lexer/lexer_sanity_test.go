// ==============================================================================================
// FILE: lexer/lexer_sanity_test.go
// ==============================================================================================
// PURPOSE: Performs a basic sanity check on the lexer.
//          It ensures that processing a standard string does not cause panic
//          and terminates gracefully at EOF.
// ==============================================================================================

package lexer

import (
	"testing"

	"eloquence/token"
)

func TestSanityLexer(t *testing.T) {
	// A simple program using typical constructs
	input := "x is 10 if x equals 10 { show(x) } end"

	l := New(input)

	// Loop until EOF is hit
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// No specific assertions here; the test passes if this loop finishes without panic.
	}
}
