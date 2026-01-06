// ==============================================================================================
// FILE: main_benchmark_test.go
// ==============================================================================================
// PURPOSE: System-wide benchmarks.
//          Measures the performance of the entire compiler pipeline (parsing + evaluation)
//          under heavy load conditions.
// ==============================================================================================

package main

import (
	"strings"
	"testing"
)

// BenchmarkSystem_HeavyLoop measures the interpretation speed of iterative logic.
func BenchmarkSystem_HeavyLoop(b *testing.B) {
	input := `
	sum is 0
	counter is 0
	limit is 1000
	
	for counter less limit
		sum is sum adds 1
		counter is counter adds 1
	end
	sum`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runCode(input)
	}
}

// BenchmarkSystem_DeepRecursion measures the overhead of stack frame allocation
// and environment switching.
func BenchmarkSystem_DeepRecursion(b *testing.B) {
	input := `
	dive is takes(n)
		if n equals 0
			return 0
		end
		return dive(n minus 1)
	end
	dive(200)`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runCode(input)
	}
}

// BenchmarkSystem_StringConcatenation measures the memory allocation overhead
// for string operations in a loop.
func BenchmarkSystem_StringConcatenation(b *testing.B) {
	// Construct a script that builds a large string
	var sb strings.Builder
	sb.WriteString(`str is "" `)
	for i := 0; i < 100; i++ {
		sb.WriteString(`str is str adds "a" `)
	}
	sb.WriteString("str")
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runCode(input)
	}
}
