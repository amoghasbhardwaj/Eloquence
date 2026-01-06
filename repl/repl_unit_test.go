// ==============================================================================================
// FILE: repl/repl_unit_test.go
// ==============================================================================================
// PURPOSE: Unit tests for basic REPL functionality.
//          Verifies that commands work and simple calculations produce output.
// ==============================================================================================

package repl

import (
	"bytes"
	"strings"
	"testing"
)

// Helper to simulate a REPL session
func runSession(input string) string {
	in := strings.NewReader(input)
	var out bytes.Buffer
	Start(in, &out)
	return out.String()
}

func TestREPL_Math(t *testing.T) {
	input := "10 adds 20\n.exit"
	output := runSession(input)

	if !strings.Contains(output, "30") {
		t.Errorf("REPL failed simple math. Output:\n%s", output)
	}
}

func TestREPL_VariablePersistence(t *testing.T) {
	// Ensure variables defined in one line persist to the next
	input := `
	x is 50
	x adds 10
	.exit`
	output := runSession(input)

	if !strings.Contains(output, "60") {
		t.Errorf("REPL failed variable persistence. Output:\n%s", output)
	}
}

func TestREPL_Commands(t *testing.T) {
	// Test .debug toggle and .clear
	input := `
	.debug
	x is 10
	.clear
	x
	.exit`
	output := runSession(input)

	// Check for debug sections
	if !strings.Contains(output, "[ TOKENS ]") {
		t.Error("Debug mode did not print tokens")
	}
	if !strings.Contains(output, "[ AST TREE ]") {
		t.Error("Debug mode did not print AST")
	}

	// Check for environment clear (x should be gone)
	if !strings.Contains(output, "identifier not found: x") {
		t.Error("Environment was not cleared correctly")
	}
}
