// ----------------------------------------------------------------------------
// FILE: lexer/lexer_integration_test.go
// ----------------------------------------------------------------------------
package lexer

import (
	"testing"

	"eloquence/token"
)

// TestIntegrationLexer tests the lexer's ability to tokenize a complex input
// simulating a struct instantiation. This verifies the interaction between
// identifiers, special syntax characters (brace, colon), and literals.
func TestIntegrationLexer(t *testing.T) {
	input := `node is Node { value: 10 }`
	expected := []struct {
		typ     token.TokenType
		literal string
	}{
		{token.IDENT, "node"},
		{token.IS, "is"},
		{token.IDENT, "Node"},
		{token.LBRACE, "{"},
		{token.IDENT, "value"},
		{token.COLON, ":"},
		{token.INT, "10"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, e := range expected {
		tok := l.NextToken()
		if tok.Type != e.typ || tok.Literal != e.literal {
			t.Fatalf("[%d] got %q %q, want %q %q", i, tok.Type, tok.Literal, e.typ, e.literal)
		}
	}
}
