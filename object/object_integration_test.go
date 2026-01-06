// ==============================================================================================
// FILE: object/object_integration_test.go
// ==============================================================================================
// PURPOSE: Integration tests for the Object system.
//          Validates the interaction between distinct object types, such as storing
//          structs inside environments or using primitives as keys in maps.
// ==============================================================================================

package object

import "testing"

func TestIntegration_ComplexStructStorage(t *testing.T) {
	// Define a struct, instantiate it, and store it in an environment
	def := &StructDefinition{Name: "Person", Fields: []string{"name", "age"}}

	instance := &StructInstance{
		Definition: def,
		Fields: map[string]Object{
			"name": &String{Value: "Alice"},
			"age":  &Integer{Value: 30},
		},
	}

	env := NewEnvironment()
	env.Set("user", instance)

	// Retrieve
	obj, ok := env.Get("user")
	if !ok {
		t.Fatalf("failed to retrieve struct")
	}

	retrievedStruct, ok := obj.(*StructInstance)
	if !ok {
		t.Fatalf("object is not a StructInstance")
	}

	// Verify fields
	nameObj := retrievedStruct.Fields["name"]
	if nameObj.(*String).Value != "Alice" {
		t.Errorf("struct field 'name' corrupted")
	}
}

func TestIntegration_MapHashing(t *testing.T) {
	// Create a map object using HashKeys
	m := &Map{Pairs: make(map[HashKey]HashPair)}

	key1 := &String{Value: "key"}
	val1 := &Integer{Value: 100}

	hashKey := key1.HashKey()
	m.Pairs[hashKey] = HashPair{Key: key1, Value: val1}

	// Store in Env
	env := NewEnvironment()
	env.Set("myMap", m)

	// Retrieve and verify
	obj, _ := env.Get("myMap")
	retrievedMap := obj.(*Map)

	// Try to look up using a fresh string object with same value
	lookupKey := &String{Value: "key"}
	pair, exists := retrievedMap.Pairs[lookupKey.HashKey()]

	if !exists {
		t.Fatalf("map lookup failed using identical string key")
	}
	if pair.Value.(*Integer).Value != 100 {
		t.Errorf("map value incorrect")
	}
}
