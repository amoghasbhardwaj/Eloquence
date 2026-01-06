// ==============================================================================================
// FILE: evaluator/evaluator_unit_test.go
// ==============================================================================================
// PURPOSE: Unit tests for specific evaluation rules.
//          Validates simple logic, arithmetic, and basic statement execution.
// ==============================================================================================

package evaluator

import (
	"testing"

	"eloquence/lexer"
	"eloquence/object"
	"eloquence/parser"
)

// Helper to bootstrap the entire pipeline for a test case
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

// Helper to assert integer values
func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 adds 5 adds 5 adds 5 minus 10", 10},
		{"2 times 2 times 2 times 2 times 2", 32},
		{"-50 adds 100 adds -50", 0},
		{"5 times 2 adds 10", 20},
		{"5 adds 2 times 10", 25},
		{"(5 adds 10 times 2 adds 15 divides 3) times 2 adds -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 less 2", true},
		{"1 greater 2", false},
		{"1 less 1", false},
		{"1 greater 1", false},
		{"1 equals 1", true},
		{"1 not_equals 1", false},
		{"1 not_equals 2", true},
		{"true equals true", true},
		{"false equals false", true},
		{"true equals false", false},
		{"true not_equals false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false}, // In Eloquence, anything not null/false is truthy
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		boolObj, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Errorf("object is not Boolean. got=%T", evaluated)
			continue
		}
		if boolObj.Value != tt.expected {
			t.Errorf("input %s: expected %t, got %t", tt.input, tt.expected, boolObj.Value)
		}
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true 10 end", 10},
		{"if false 10 end", nil},
		{"if 1 10 end", 10},
		{"if 1 less 2 10 end", 10},
		{"if 1 greater 2 10 end", nil},
		{"if 1 greater 2 10 else 20 end", 20},
		{"if 1 less 2 10 else 20 end", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			if evaluated != NULL {
				t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
			}
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10", 10},
		{"return 10; 9;", 10},
		{"return 2 times 5; 9;", 10},
		{"9; return 2 times 5; 9;", 10},
		{
			`if 10 greater 1
				if 10 greater 1
					return 10
				end
				return 1
			end`, 10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 adds true", "type mismatch: INTEGER adds BOOLEAN"},
		{"5 adds true; 5;", "type mismatch: INTEGER adds BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true adds false", "unknown operator: BOOLEAN adds BOOLEAN"},
		{"5; true adds false; 5", "unknown operator: BOOLEAN adds BOOLEAN"},
		{"if 10 greater 1; true adds false; end", "unknown operator: BOOLEAN adds BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}
