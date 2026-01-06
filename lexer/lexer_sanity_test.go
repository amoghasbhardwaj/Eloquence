// ----------------------------------------------------------------------------
// FILE: lexer/lexer_sanity_test.go
// ----------------------------------------------------------------------------
package lexer

import (
	"testing"

	"eloquence/token"
)

// TestSanityLexer performs a basic sanity check on the lexer.
// It ensures that processing a standard string does not cause panic
// and terminates gracefully at EOF.
func TestSanityLexer(t *testing.T) {
	input := "x is 10 if x equals 10 show x end"
	l := New(input)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// Just sanity check: no panic
	}
}
