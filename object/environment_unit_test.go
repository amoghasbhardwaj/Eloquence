// ==============================================================================================
// FILE: object/environment_unit_test.go
// ==============================================================================================
// PURPOSE: Specific unit tests for the Environment struct.
//          Validates shadowing rules, scope traversal, and variable persistence.
// ==============================================================================================

package object

import "testing"

func TestEnvironment_GetSet(t *testing.T) {
	env := NewEnvironment()

	// 1. Test Retrieval of non-existent variable
	if _, ok := env.Get("x"); ok {
		t.Errorf("expected 'x' to not exist")
	}

	// 2. Test Set and Get
	val := &Integer{Value: 10}
	env.Set("x", val)

	result, ok := env.Get("x")
	if !ok {
		t.Fatalf("expected 'x' to exist")
	}
	if result != val {
		t.Errorf("expected got %v, want %v", result, val)
	}
}

func TestEnclosedEnvironments(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Integer{Value: 10})
	outer.Set("y", &Integer{Value: 5})

	inner := NewEnclosedEnvironment(outer)

	// 1. Test reading from outer scope
	val, ok := inner.Get("x")
	if !ok || val.(*Integer).Value != 10 {
		t.Errorf("failed to read from outer scope")
	}

	// 2. Test Shadowing (inner variable overrides outer)
	// 'x' is redefined in the inner scope
	inner.Set("x", &Integer{Value: 99})

	valInner, _ := inner.Get("x")
	if valInner.(*Integer).Value != 99 {
		t.Errorf("inner scope did not shadow outer scope")
	}

	valOuter, _ := outer.Get("x")
	if valOuter.(*Integer).Value != 10 {
		t.Errorf("outer scope was modified by inner set (shadowing failed)")
	}

	// 3. Test variable that only exists in outer
	yVal, ok := inner.Get("y")
	if !ok || yVal.(*Integer).Value != 5 {
		t.Errorf("failed to traverse up to outer scope")
	}
}
