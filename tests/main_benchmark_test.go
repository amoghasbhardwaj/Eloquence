// ==============================================================================================
// FILE: main_benchmark_test.go
// ==============================================================================================
// PURPOSE: System-wide benchmarks.
// ==============================================================================================

package main

import (
	"strings"
	"testing"
)

func BenchmarkSystem_HeavyLoop(b *testing.B) {
	input := `
	sum is 0
	counter is 0
	limit is 1000
	
	while counter less limit {
		sum is sum adds 1
		counter is counter adds 1
	}
	sum`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runCode(input)
	}
}

func BenchmarkSystem_DeepRecursion(b *testing.B) {
	input := `
	dive is takes(n) {
		if n equals 0 {
			return 0
		}
		return dive(n minus 1)
	}
	dive(200)`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runCode(input)
	}
}

func BenchmarkSystem_StringConcatenation(b *testing.B) {
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
