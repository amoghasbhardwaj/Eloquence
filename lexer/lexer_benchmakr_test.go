// ==============================================================================================
// FILE: lexer/lexer_benchmark_test.go
// ==============================================================================================
// PURPOSE: Benchmarks the throughput of the lexical analysis.
//          It simulates a hot loop of tokenizing a standard expression to ensure low latency.
// ==============================================================================================

package lexer

import (
	"testing"

	"eloquence/token"
)

// BenchmarkLexerNextToken measures the performance of scanning.
// Command to run: go test -bench=. ./lexer
func BenchmarkLexerNextToken(b *testing.B) {
	// A representative string containing identifiers, keywords, numbers, and operators.
	input := `x is 1 y is 2 z is 3 a adds b c subtracts d`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			// Consumption loop to trigger full tokenization
		}
	}
}
