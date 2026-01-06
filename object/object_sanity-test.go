// ==============================================================================================
// FILE: object/object_sanity_test.go
// ==============================================================================================
// PURPOSE: Sanity checks for the Object system.
//          Verifies that empty collections behave correctly and deep recursion doesn't crash.
// ==============================================================================================

package object

import "testing"

func TestSanity_EmptyCollections(t *testing.T) {
	// Empty Array
	arr := &Array{Elements: []Object{}}
	if arr.Inspect() != "[]" {
		t.Errorf("empty array inspect failed")
	}

	// Empty Map
	m := &Map{Pairs: map[HashKey]HashPair{}}
	if m.Inspect() != "{}" {
		t.Errorf("empty map inspect failed")
	}
}

func TestSanity_NestedEnvironments(t *testing.T) {
	// Create a chain of 100 environments to ensure no stack overflow on simple lookup
	root := NewEnvironment()
	root.Set("target", &Boolean{Value: true})

	current := root
	for i := 0; i < 100; i++ {
		current = NewEnclosedEnvironment(current)
	}

	val, ok := current.Get("target")
	if !ok {
		t.Fatalf("deep nested lookup failed")
	}
	if val.Inspect() != "true" {
		t.Errorf("deep nested value corrupted")
	}
}
