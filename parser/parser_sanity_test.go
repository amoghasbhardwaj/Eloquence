// ==============================================================================================
// FILE: parser/parser_sanity_test.go
// ==============================================================================================
// PURPOSE: Sanity checks for the Parser.
// ==============================================================================================

package parser

import (
	"testing"

	"eloquence/lexer"
)

func TestSanity_EmptyInput(t *testing.T) {
	input := "   \n  \t  "
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser reported errors on empty input: %v", p.Errors())
	}
	if len(program.Statements) != 0 {
		t.Errorf("expected 0 statements for empty input")
	}
}

func TestSanity_CommentsOnly(t *testing.T) {
	input := `
    /* This is a comment */
    /* Another one */
    // Single line
    `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser errors on comments: %v", p.Errors())
	}
	if len(program.Statements) != 0 {
		t.Errorf("expected 0 statements for comments")
	}
}

func TestSanity_GracefulErrorHandling(t *testing.T) {
	// Missing value after 'is'
	input := `x is`
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Errorf("expected parser errors for incomplete assignment")
	}
}

func TestSanity_UnterminatedBlock(t *testing.T) {
	// Missing '}' - Expects parser error now
	input := `if x less 5 {
        show(x)`

	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Errorf("expected parser errors for unterminated block, got none")
	} else {
		// Optional: verify error message content
		expectedMsg := "unterminated block: expected '}', got EOF"
		found := false
		for _, err := range p.Errors() {
			if err == expectedMsg {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected error %q, got %v", expectedMsg, p.Errors())
		}
	}
}