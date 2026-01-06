// ==============================================================================================
// FILE: evaluator/evaluator_sanity_test.go
// ==============================================================================================
// PURPOSE: Sanity checks for the runtime.
//          Ensures that invalid programs fail gracefully and empty programs
//          return expected nil/null results.
// ==============================================================================================

package evaluator

import (
	"testing"

	"eloquence/object"
)

func TestSanity_EmptyProgram(t *testing.T) {
	input := ""
	evaluated := testEval(input)
	if evaluated != nil {
		t.Errorf("empty program expected nil result, got %T", evaluated)
	}
}

func TestSanity_DanglingPointer(t *testing.T) {
	// This tests a pointer that refers to an identifier that doesn't exist
	// Although lexically valid, it should fail at runtime
	input := `
	ptr is pointing to missing
	pointing from ptr`

	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("expected error for dangling pointer, got %T", evaluated)
	}
	if errObj.Message != "identifier not found: missing" {
		t.Errorf("unexpected error message: %s", errObj.Message)
	}
}

func TestSanity_UnknownStructField(t *testing.T) {
	input := `
	define Box as struct { item }
	b is Box { item: 1 }
	b.missing`

	evaluated := testEval(input)
	_, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("expected error for missing field, got %T", evaluated)
	}
}
