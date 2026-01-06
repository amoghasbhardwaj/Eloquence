// ==============================================================================================
// FILE: repl/repl_benchmark_test.go
// ==============================================================================================
// PURPOSE: Performance benchmarks for the REPL loop.
//          Measures startup overhead and input processing latency.
// ==============================================================================================

package repl

import (
	"bytes"
	"strings"
	"testing"
)

// BenchmarkREPL_StartupAndExit measures the cost of initializing the REPL environment.
func BenchmarkREPL_StartupAndExit(b *testing.B) {
	input := ".exit"
	for i := 0; i < b.N; i++ {
		in := strings.NewReader(input)
		var out bytes.Buffer
		Start(in, &out)
	}
}

// BenchmarkREPL_Calculation measures throughput for a simple calculation cycle.
func BenchmarkREPL_Calculation(b *testing.B) {
	input := "10 times 10 adds 5\n.exit"
	for i := 0; i < b.N; i++ {
		in := strings.NewReader(input)
		var out bytes.Buffer
		Start(in, &out)
	}
}
