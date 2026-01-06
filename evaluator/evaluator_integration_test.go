// ==============================================================================================
// FILE: evaluator/evaluator_integration_test.go
// ==============================================================================================
// PURPOSE: Integration tests for the runtime.
//          Validates advanced features where multiple language constructs interact,
//          such as functions using structs, pointers modifying variables, etc.
// ==============================================================================================

package evaluator

import (
	"testing"
)

func TestIntegration_FunctionApplication(t *testing.T) {
	input := `
	double is takes(x)
		x times 2
	end
	
	add is takes(x, y)
		x adds y
	end
	
	add(double(5), add(3, 2))`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestIntegration_Closures(t *testing.T) {
	// Tests lexical scoping: 'x' is captured by the inner function
	input := `
	newAdder is takes(x)
		takes(y)
			x adds y
		end
	end
	
	addTwo is newAdder(2)
	addTwo(2)`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}

func TestIntegration_Structs(t *testing.T) {
	input := `
	define User as struct { name, age }
	u is User { name: "Alice", age: 30 }
	u.age`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 30)
}

func TestIntegration_Pointers(t *testing.T) {
	// Tests that pointer modification updates the original variable
	input := `
	x is 10
	ptr is pointing to x
	pointing from ptr is 20
	x` // x should now be 20

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 20)
}

func TestIntegration_RecursiveFactorial(t *testing.T) {
	input := `
	fact is takes(n)
		if n equals 0
			return 1
		else
			return n times fact(n subtracts 1)
		end
	end
	fact(5)`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 120)
}

func TestIntegration_MapAndArray(t *testing.T) {
	input := `
	arr is [1, 2, 3]
	m is { "a": 10, "b": 20 }
	arr[0] adds m["b"]` // 1 + 20 = 21

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 21)
}
