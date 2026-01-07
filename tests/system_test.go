// ==============================================================================================
// FILE: system_test.go
// ==============================================================================================
// PURPOSE: System-level integration tests.
//          These tests verify that all components (Lexer -> Parser -> Evaluator) work together
//          to execute valid Eloquence logic.
// ==============================================================================================

package main

import (
	"testing"

	"eloquence/evaluator"
	"eloquence/lexer"
	"eloquence/object"
	"eloquence/parser"
)

func runCode(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	// Fail immediately on parse errors to aid debugging
	if len(p.Errors()) > 0 {
		return &object.Error{Message: "PARSER ERROR: " + p.Errors()[0]}
	}

	env := object.NewEnvironment()
	return evaluator.Eval(program, env)
}

func assertInteger(t *testing.T, obj object.Object, expected int64) {
	if obj == nil {
		t.Fatalf("got nil object")
	}
	// Check for errors first
	if err, ok := obj.(*object.Error); ok {
		t.Fatalf("runtime error: %s", err.Message)
	}

	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Result is not Integer. Got=%T (%+v)", obj, obj)
		return
	}
	if result.Value != expected {
		t.Errorf("Wrong integer value. Expected=%d, Got=%d", expected, result.Value)
	}
}

func TestSystem_Fibonacci_Recursion(t *testing.T) {
	input := `
	fib is takes(x) {
		if x less 2 {
			return x
		}
		return fib(x minus 1) adds fib(x minus 2)
	}
	fib(10)`

	result := runCode(input)
	assertInteger(t, result, 55)
}

func TestSystem_MapReduce_HigherOrderFunctions(t *testing.T) {
	input := `
	map is takes(arr, func) {
		// Simulating iteration for the test case
		val1 is func(arr[0])
		val2 is func(arr[1])
		val3 is func(arr[2])
		return val3
	}

	double is takes(x) {
		return x times 2 
	}

	arr is [10, 20, 30]
	map(arr, double)
	`

	result := runCode(input)
	assertInteger(t, result, 60) // 30 * 2
}

func TestSystem_LinkedList(t *testing.T) {
	input := `
	define Node as struct { val, next }
	
	node3 is Node { val: 30, next: none }
	node2 is Node { val: 20, next: node3 }
	head  is Node { val: 10, next: node2 }
	
	sumList is takes(node) {
		if node equals none {
			return 0
		}
		return node.val adds sumList(node.next)
	}
	
	sumList(head)`

	result := runCode(input)
	assertInteger(t, result, 60) // 10 + 20 + 30
}

func TestSystem_PointerMutation(t *testing.T) {
	input := `
	globalVal is 100
	
	mutate is takes() {
		ptr is pointing to globalVal
		pointing from ptr is 999
	}
	
	mutate()
	globalVal`

	result := runCode(input)
	assertInteger(t, result, 999)
}

func TestSystem_ShadowingAndScope(t *testing.T) {
	input := `
	x is 10
	if true {
		x is 20       
		x is x adds 1 
	}
	x`

	result := runCode(input)
	assertInteger(t, result, 10)
}

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
		t.Fatalf("Expected error for dangling pointer, got %s", result.Type())
	}
}
