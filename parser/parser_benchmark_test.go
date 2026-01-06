// ==============================================================================================
// FILE: parser/parser_benchmark_test.go
// ==============================================================================================
// PURPOSE: Performance benchmarks for the Parser.
//          Measures parsing throughput for simple assignments, large programs, and
//          deeply nested expressions to ensure the parser scales linearly.
// ==============================================================================================

package parser

import (
	"fmt"
	"strings"
	"testing"

	"eloquence/lexer"
)

// BenchmarkParser_SimpleAssignment measures the cost of parsing a single basic statement.
// Usage: go test -bench=BenchmarkParser_SimpleAssignment ./parser
func BenchmarkParser_SimpleAssignment(b *testing.B) {
	input := "x is 5"
	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkParser_LargeProgram measures parsing speed for a 1000-line file.
// Usage: go test -bench=BenchmarkParser_LargeProgram ./parser
func BenchmarkParser_LargeProgram(b *testing.B) {
	// Generate a program with 1000 lines of assignments
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString(fmt.Sprintf("var%d is %d\n", i, i))
	}
	input := sb.String()

	b.ResetTimer() // Don't count setup time

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkParser_DeeplyNestedMath measures recursive parsing depth efficiency.
// Usage: go test -bench=BenchmarkParser_DeeplyNestedMath ./parser
func BenchmarkParser_DeeplyNestedMath(b *testing.B) {
	// Generate: result is 1 adds 1 adds 1 adds ... (Left recursive structure)
	var sb strings.Builder
	sb.WriteString("result is 1")
	for i := 0; i < 100; i++ {
		sb.WriteString(" adds 1")
	}
	input := sb.String()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		p.ParseProgram()
	}
}
