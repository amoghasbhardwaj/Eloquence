// ==============================================================================================
// FILE: object/object_unit_test.go
// ==============================================================================================
// PURPOSE: Unit tests for Object methods.
//          Verifies that Inspect() produces correct string representations and
//          Type() returns the correct constants.
// ==============================================================================================

package object

import (
	"testing"
)

func TestObjectInspect(t *testing.T) {
	tests := []struct {
		obj      Object
		expected string
	}{
		// Primitives
		{&Integer{Value: 10}, "10"},
		{&Float{Value: 3.14}, "3.14"},
		{&Boolean{Value: true}, "true"},
		{&Boolean{Value: false}, "false"},
		{&String{Value: "hello"}, "hello"},
		{&Char{Value: 'a'}, "a"},
		{&Null{}, "none"},

		// Internal
		{&ReturnValue{Value: &Integer{Value: 5}}, "5"},
		{&Error{Message: "something went wrong"}, "ERROR: something went wrong"},

		// Complex
		{&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}}, "[1, 2]"},
		{&Function{}, "takes(...) { ... }"},
		{&StructDefinition{Name: "User"}, "struct User"},
		{&Pointer{Name: "ptr"}, "pointing to ptr"},
	}

	for _, tt := range tests {
		if tt.obj.Inspect() != tt.expected {
			t.Errorf("Inspect() wrong. expected=%q, got=%q", tt.expected, tt.obj.Inspect())
		}
	}
}

func TestObjectType(t *testing.T) {
	tests := []struct {
		obj          Object
		expectedType ObjectType
	}{
		{&Integer{Value: 5}, INTEGER_OBJ},
		{&Boolean{Value: true}, BOOLEAN_OBJ},
		{&String{Value: "x"}, STRING_OBJ},
		{&Null{}, NULL_OBJ},
		{&Array{}, ARRAY_OBJ},
		{&Map{}, MAP_OBJ},
		{&StructInstance{}, STRUCT_INST_OBJ},
	}

	for _, tt := range tests {
		if tt.obj.Type() != tt.expectedType {
			t.Errorf("Type() wrong. expected=%q, got=%q", tt.expectedType, tt.obj.Type())
		}
	}
}

func TestHashKeys(t *testing.T) {
	// 1. Verify two identical values produce the same hash key
	int1 := &Integer{Value: 5}
	int2 := &Integer{Value: 5}
	if int1.HashKey() != int2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	bool1 := &Boolean{Value: true}
	bool2 := &Boolean{Value: true}
	if bool1.HashKey() != bool2.HashKey() {
		t.Errorf("booleans with same value have different hash keys")
	}

	str1 := &String{Value: "hello"}
	str2 := &String{Value: "hello"}
	if str1.HashKey() != str2.HashKey() {
		t.Errorf("strings with same value have different hash keys")
	}

	// 2. Verify distinct values produce different hash keys
	int3 := &Integer{Value: 10}
	if int1.HashKey() == int3.HashKey() {
		t.Errorf("different integers have same hash key")
	}
}
