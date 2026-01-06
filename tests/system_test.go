// ==============================================================================================
// FILE: system_test.go
// ==============================================================================================
// PURPOSE: System-level integration tests.
//          These tests verify that all components (Lexer -> Parser -> Evaluator) work together
//          to execute valid Eloquence logic. They act as "Turing Completeness" verifications.
// ==============================================================================================

package main

import (
	"testing"

	"eloquence/evaluator"
	"eloquence/lexer"
	"eloquence/object"
	"eloquence/parser"
)

// Helper: Executes a string of Eloquence code and returns the final result.
func runCode(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return evaluator.Eval(program, env)
}

// Helper: Asserts that an object is a specific Integer value.
func assertInteger(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Result is not Integer. Got=%T (%+v)", obj, obj)
		return
	}
	if result.Value != expected {
		t.Errorf("Wrong integer value. Expected=%d, Got=%d", expected, result.Value)
	}
}

// ----------------------------------------------------------------------------
// 1. ALGORITHM TESTS
// ----------------------------------------------------------------------------

func TestSystem_Fibonacci_Recursion(t *testing.T) {
	// Tests recursion, if/else logic, and arithmetic precedence.
	input := `
	fib is takes(x)
		if x less 2
			return x
		end
		return fib(x minus 1) adds fib(x minus 2)
	end
	fib(10)`

	result := runCode(input)
	assertInteger(t, result, 55)
}

func TestSystem_MapReduce_HigherOrderFunctions(t *testing.T) {
	// Tests first-class functions, array passing, and loop accumulators.
	// We simulate a 'map' function that applies a transformation to an array.
	input := `
	map is takes(arr, func)
		result is []
		
		// In a full stdlib, we would have len(arr). 
		// Here we simulate iteration for the test case.
		val1 is func(arr[0])
		val2 is func(arr[1])
		val3 is func(arr[2])
		
		// Simulate array construction since push() isn't primitive yet
		// We just return the last mapped value to verify the function ran.
		return val3
	end

	double is takes(x) 
		return x times 2 
	end

	arr is [10, 20, 30]
	map(arr, double)
	`

	result := runCode(input)
	assertInteger(t, result, 60) // 30 * 2
}

// ----------------------------------------------------------------------------
// 2. DATA STRUCTURE TESTS
// ----------------------------------------------------------------------------

func TestSystem_LinkedList(t *testing.T) {
	// Verifies that Structs can refer to themselves (recursive types) and be traversed.
	input := `
	define Node as struct { val, next }
	
	// Create List: 10 -> 20 -> 30 -> none
	node3 is Node { val: 30, next: none }
	node2 is Node { val: 20, next: node3 }
	head  is Node { val: 10, next: node2 }
	
	// Traverse recursively to sum values
	sumList is takes(node)
		if node equals none
			return 0
		end
		return node.val adds sumList(node.next)
	end
	
	sumList(head)`

	result := runCode(input)
	assertInteger(t, result, 60) // 10 + 20 + 30
}

// ----------------------------------------------------------------------------
// 3. MEMORY & SCOPE TESTS
// ----------------------------------------------------------------------------

func TestSystem_PointerMutation(t *testing.T) {
	// Tests reference semantics (pass-by-reference).
	input := `
	globalVal is 100
	
	mutate is takes()
		ptr is pointing to globalVal
		pointing from ptr is 999
	end
	
	mutate()
	globalVal`

	result := runCode(input)
	assertInteger(t, result, 999)
}

func TestSystem_ShadowingAndScope(t *testing.T) {
	// Tests that block-scoped variables in 'if' do not leak,
	// but can shadow outer variables temporarily.
	input := `
	x is 10
	if true
		x is 20       // This defines a NEW x in the local block scope
		x is x adds 1 // Local x becomes 21
	end
	x` // Outer x should remain 10

	result := runCode(input)
	assertInteger(t, result, 10)
}

// ----------------------------------------------------------------------------
// 4. EDGE CASE TESTS
// ----------------------------------------------------------------------------

func TestSystem_EdgeCase_DivisionByZero(t *testing.T) {
	input := `10 divides 0`
	result := runCode(input)

	if result.Type() != object.ERROR_OBJ {
		t.Fatalf("Expected error for division by zero, got %s", result.Type())
	}
}

func TestSystem_EdgeCase_DanglingPointer(t *testing.T) {
	input := `
	ptr is pointing to nothing
	pointing from ptr`

	result := runCode(input)
	if result.Type() != object.ERROR_OBJ {
		t.Fatalf("Expected error for dangling pointer")
	}
}
