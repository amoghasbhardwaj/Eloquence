// ==============================================================================================
// FILE: lexer/lexer_unit_test.go
// ==============================================================================================
// PURPOSE: Validates that the Lexer correctly identifies all token types and literals.
// ==============================================================================================

package lexer

import (
	"testing"

	"eloquence/token"
)

// TestNextToken checks that the lexer correctly produces tokens
// for all token types in the English-first language.
func TestNextToken(t *testing.T) {
	// --- SECTION 1: Identifiers, assignment, numbers, strings, booleans ---
	input1 := `
x is 10
y is 20
name is "Amogh"
flag is true
pi is 3.14
`
	expected1 := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// x is 10
		{token.IDENT, "x"},
		{token.IS, "is"},
		{token.INT, "10"},

		// y is 20
		{token.IDENT, "y"},
		{token.IS, "is"},
		{token.INT, "20"},

		// name is "Amogh"
		{token.IDENT, "name"},
		{token.IS, "is"},
		{token.STRING, "Amogh"},

		// flag is true
		{token.IDENT, "flag"},
		{token.IS, "is"},
		{token.BOOL, "true"},

		// pi is 3.14
		{token.IDENT, "pi"},
		{token.IS, "is"},
		{token.FLOAT, "3.14"},

		// EOF
		{token.EOF, ""},
	}
	runLexerTest(t, input1, expected1)

	// --- SECTION 2: Arithmetic operators ---
	input2 := `
a adds b
c subtracts d
e times f
g divides h
i modulo j
`
	expected2 := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "a"},
		{token.ADDS, "adds"},
		{token.IDENT, "b"},

		{token.IDENT, "c"},
		{token.SUBTRACTS, "subtracts"},
		{token.IDENT, "d"},

		{token.IDENT, "e"},
		{token.TIMES, "times"},
		{token.IDENT, "f"},

		{token.IDENT, "g"},
		{token.DIVIDES, "divides"},
		{token.IDENT, "h"},

		{token.IDENT, "i"},
		{token.MODULO, "modulo"},
		{token.IDENT, "j"},

		{token.EOF, ""},
	}
	runLexerTest(t, input2, expected2)

	// --- SECTION 3: Comparison operators ---
	input3 := `
x equals y
a not_equals b
c greater d
e less f
g greater_equal h
i less_equal j
`
	expected3 := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "x"},
		{token.EQUALS, "equals"},
		{token.IDENT, "y"},

		{token.IDENT, "a"},
		{token.NOT_EQUALS, "not_equals"},
		{token.IDENT, "b"},

		{token.IDENT, "c"},
		{token.GREATER, "greater"},
		{token.IDENT, "d"},

		{token.IDENT, "e"},
		{token.LESS, "less"},
		{token.IDENT, "f"},

		{token.IDENT, "g"},
		{token.GREATER_EQUAL, "greater_equal"},
		{token.IDENT, "h"},

		{token.IDENT, "i"},
		{token.LESS_EQUAL, "less_equal"},
		{token.IDENT, "j"},

		{token.EOF, ""},
	}
	runLexerTest(t, input3, expected3)

	// --- SECTION 4: Logical operators ---
	input4 := `
x and y
a or b
not flag
`
	expected4 := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "x"},
		{token.AND, "and"},
		{token.IDENT, "y"},

		{token.IDENT, "a"},
		{token.OR, "or"},
		{token.IDENT, "b"},

		{token.NOT, "not"},
		{token.IDENT, "flag"},

		{token.EOF, ""},
	}
	runLexerTest(t, input4, expected4)

	// --- SECTION 5: Control flow and output ---
	input5 := `
if x equals 10
show(x)
else
show(y)
end
return x
`
	expected5 := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.EQUALS, "equals"},
		{token.INT, "10"},

		// UPDATED: 'show' is now an IDENT (Function Call)
		{token.IDENT, "show"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},

		{token.ELSE, "else"},

		// UPDATED: 'show' is now an IDENT (Function Call)
		{token.IDENT, "show"},
		{token.LPAREN, "("},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},

		{token.END, "end"},

		{token.RETURN, "return"},
		{token.IDENT, "x"},

		{token.EOF, ""},
	}
	runLexerTest(t, input5, expected5)
}

// runLexerTest is a helper to iterate expected tokens and check against lexer output
func runLexerTest(t *testing.T, input string, expectedTokens []struct {
	expectedType    token.TokenType
	expectedLiteral string
},
) {
	lexer := New(input)

	for i, expected := range expectedTokens {
		actual := lexer.NextToken()

		if actual.Type != expected.expectedType {
			t.Fatalf(
				"tests[%d] - token type mismatch. expected=%q, got=%q",
				i, expected.expectedType, actual.Type,
			)
		}

		if actual.Literal != expected.expectedLiteral {
			t.Fatalf(
				"tests[%d] - token literal mismatch. expected=%q, got=%q",
				i, expected.expectedLiteral, actual.Literal,
			)
		}
	}
}
