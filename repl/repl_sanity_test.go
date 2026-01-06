// ==============================================================================================
// FILE: repl/repl_sanity_test.go
// ==============================================================================================
// PURPOSE: Sanity checks for the REPL.
//          Ensures robust handling of edge cases like empty lines and bad commands.
// ==============================================================================================

package repl

import (
	"strings"
	"testing"
)

func TestSanity_EmptyLines(t *testing.T) {
	input := "\n\n\n\n10\n.exit"
	output := runSession(input)
	if !strings.Contains(output, "10") {
		t.Error("REPL choked on empty lines")
	}
}

func TestSanity_ParseErrors(t *testing.T) {
	input := "if x less\n.exit"
	output := runSession(input)
	if !strings.Contains(output, "Parser Errors") {
		t.Error("REPL did not report parser errors gracefully")
	}
}

func TestSanity_UnknownCommand(t *testing.T) {
	input := ".foobar\n.exit"
	output := runSession(input)
	if !strings.Contains(output, "Unknown command") {
		t.Error("REPL did not catch unknown command")
	}
}
