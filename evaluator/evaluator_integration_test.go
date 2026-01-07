// ==============================================================================================
// FILE: evaluator/evaluator_integration_test.go
// ==============================================================================================
// PURPOSE: Integration tests for the Evaluator.
//          Validates complex, multi-statement logic like recursion, closures, and structs.
// ==============================================================================================

package evaluator

import (
	"testing"
)

func TestIntegration_FunctionApplication(t *testing.T) {
	input := `
	identity is takes(x) { x }
	identity(5)`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func TestIntegration_Closures(t *testing.T) {
	input := `
	newAdder is takes(x) {
		return takes(y) { x adds y }
	}
	addTwo is newAdder(2)
	addTwo(2)`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}

func TestIntegration_RecursiveFactorial(t *testing.T) {
	input := `
	factorial is takes(n) {
		if n equals 0 {
			return 1
		}
		return n times factorial(n minus 1)
	}
	factorial(5)`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 120)
}

func TestIntegration_Structs(t *testing.T) {
	input := `
	define Box as struct { width, height }
	b is Box { width: 10, height: 20 }
	b.width times b.height`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 200)
}

func TestIntegration_Pointers(t *testing.T) {
	input := `
	val is 50
	ptr is pointing to val
	pointing from ptr is 100
	val`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 100)
}

func TestIntegration_MapAndArray(t *testing.T) {
	input := `
	arr is [1, 2, 3]
	dict is { "first": arr[0] }
	dict["first"]`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 1)
}
