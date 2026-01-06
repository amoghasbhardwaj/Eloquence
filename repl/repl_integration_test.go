// ==============================================================================================
// FILE: repl/repl_integration_test.go
// ==============================================================================================
// PURPOSE: Integration tests for the REPL.
//          Validates multi-line interactions involving complex types (structs, functions).
// ==============================================================================================

package repl

import (
	"strings"
	"testing"
)

func TestIntegration_ComplexSession(t *testing.T) {
	input := `
	define User as struct { name, age }
	u is User { name: "Amogh", age: 25 }
	
	age_checker is takes(person)
		if person.age greater 18
			return "Adult"
		else
			return "Minor"
		end
	end
	
	age_checker(u)
	.exit`

	output := runSession(input)

	// We expect "Adult" in the output
	if !strings.Contains(output, "Adult") {
		t.Errorf("Complex struct integration failed. Output:\n%s", output)
	}
}

func TestIntegration_PointersInRepl(t *testing.T) {
	input := `
	x is 100
	p is pointing to x
	pointing from p is 200
	x
	.exit`

	output := runSession(input)

	// Value of x should change to 200 via pointer mutation
	if !strings.Contains(output, "200") {
		t.Errorf("Pointer integration failed. Output:\n%s", output)
	}
}
