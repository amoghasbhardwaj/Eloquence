// ==============================================================================================
// FILE: token/token_benchmark_test.go
// ==============================================================================================
// PURPOSE: Benchmarks the LookupIdent function. Since this is called for every identifier in the
//          source code, it must be extremely fast (zero allocation if possible).
// ==============================================================================================

package token

import "testing"

// BenchmarkLookupIdent measures the performance of keyword lookups.
// Running: go test -bench=.
func BenchmarkLookupIdent(b *testing.B) {
	// A mix of keywords and identifiers typical in code
	words := []string{
		"if", "adds", "subtracts",
		"pointing to", "include", "in",
		"unknown_var", "myFunction", "return",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			// discard result to focus on execution time
			_ = LookupIdent(w)
		}
	}
}
