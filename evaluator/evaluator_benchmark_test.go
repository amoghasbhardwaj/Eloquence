// ==============================================================================================
// FILE: evaluator/evaluator_benchmark_test.go
// ==============================================================================================
// PURPOSE: Performance benchmarks for the runtime.
//          Measures the speed of interpretation for CPU-intensive tasks like
//          deep recursion and large loops.
// ==============================================================================================

package evaluator

import (
	"strings"
	"testing"
)

// BenchmarkEvaluator_Fibonacci measures recursion overhead (stack frames, env creation).
// Usage: go test -bench=BenchmarkEvaluator_Fibonacci ./evaluator
func BenchmarkEvaluator_Fibonacci(b *testing.B) {
	input := `
	fib is takes(x)
		if x equals 0
			return 0
		end
		if x equals 1
			return 1
		end
		return fib(x subtracts 1) adds fib(x subtracts 2)
	end
	fib(10)`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testEval(input)
	}
}

// BenchmarkEvaluator_LargeArraySum measures loop overhead and variable lookups.
// Usage: go test -bench=BenchmarkEvaluator_LargeArraySum ./evaluator
func BenchmarkEvaluator_LargeArraySum(b *testing.B) {
	// Create a program that sums a large array
	var sb strings.Builder
	sb.WriteString("arr is [")
	for i := 0; i < 100; i++ {
		sb.WriteString("1")
		if i < 99 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]\n")
	sb.WriteString(`
	sum is 0
	i is 0
	len is 100
	for i less len
		sum is sum adds arr[i]
		i is i adds 1
	end
	sum`)
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testEval(input)
	}
}
