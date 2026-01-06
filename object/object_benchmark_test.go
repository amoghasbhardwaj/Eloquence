// ==============================================================================================
// FILE: object/object_benchmark_test.go
// ==============================================================================================
// PURPOSE: Performance benchmarks for the Object system.
//          Measures hashing costs, environment access time, and object creation overhead.
// ==============================================================================================

package object

import (
	"fmt"
	"testing"
)

// BenchmarkHashKey_String measures the cost of hashing a string.
// Important for Map performance.
func BenchmarkHashKey_String(b *testing.B) {
	s := &String{Value: "some_long_identifier_name"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.HashKey()
	}
}

// BenchmarkHashKey_Integer measures the cost of hashing an integer.
func BenchmarkHashKey_Integer(b *testing.B) {
	num := &Integer{Value: 123456789}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		num.HashKey()
	}
}

// BenchmarkEnvironment_Get_Deep measures lookup time in a deeply nested scope.
func BenchmarkEnvironment_Get_Deep(b *testing.B) {
	// Setup a deep environment chain
	root := NewEnvironment()
	root.Set("target", &Integer{Value: 1})

	curr := root
	for i := 0; i < 50; i++ {
		curr = NewEnclosedEnvironment(curr)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		curr.Get("target")
	}
}

func BenchmarkObjectInspect_LargeArray(b *testing.B) {
	elements := make([]Object, 100)
	for i := 0; i < 100; i++ {
		elements[i] = &Integer{Value: int64(i)}
	}
	arr := &Array{Elements: elements}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arr.Inspect()
	}
}

func BenchmarkEnvironment_Set(b *testing.B) {
	env := NewEnvironment()
	val := &Integer{Value: 1}
	keys := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = fmt.Sprintf("var%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Cycle through keys to avoid simple overwrite optimization
		env.Set(keys[i%1000], val)
	}
}
