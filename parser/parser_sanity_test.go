// ==============================================================================================
// FILE: parser/parser_sanity_test.go
// ==============================================================================================
// PURPOSE: Sanity checks for the Parser.
//          Ensures the parser handles empty files, comments, and invalid syntax
//          gracefully (by reporting errors) rather than crashing.
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
		t.Errorf("expected 0 statements for empty input, got %d", len(program.Statements))
	}
}

func TestSanity_CommentsOnly(t *testing.T) {
	input := `
    /* This is a comment */
    /* Another one */
    `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser errors on comments: %v", p.Errors())
	}
	if len(program.Statements) != 0 {
		t.Errorf("expected 0 statements for comments, got %d", len(program.Statements))
	}
}

func TestSanity_GracefulErrorHandling(t *testing.T) {
	// Missing value after 'is'
	input := `x is`
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Errorf("expected parser errors for incomplete assignment, got none")
	}
}

func TestSanity_UnterminatedBlock(t *testing.T) {
	// Missing 'end'
	input := `if x less 5
        show x`

	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Errorf("expected parser errors for unterminated block, got none")
	}
}
