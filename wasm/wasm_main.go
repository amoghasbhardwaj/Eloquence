// ==============================================================================================
// FILE: wasm/wasm_main.go
// BUILD: GOOS=js GOARCH=wasm go build -o main.wasm wasm/wasm_main.go
// ==============================================================================================
package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"eloquence/ast"
	"eloquence/evaluator"
	"eloquence/lexer"
	"eloquence/object"
	"eloquence/parser"
)

// We use this buffer to capture output from "show()" calls
var outputBuffer strings.Builder

func main() {
	// Create a channel to keep the Go WASM running
	c := make(chan struct{}, 0)

	// Override Builtins for the Web Environment
	overrideBuiltinsForWeb()

	// Expose the function to JavaScript
	js.Global().Set("runEloquence", js.FuncOf(runCode))

	fmt.Println("Eloquence WASM Engine Loaded.")
	<-c
}

// runCode is the bridge between JS and Go
func runCode(this js.Value, p []js.Value) interface{} {
	code := p[0].String()

	// Reset output buffer for this run
	outputBuffer.Reset()

	// 1. Setup Environment
	env := object.NewEnvironment()

	// 2. Setup Parser Hook (Disable Include for Web)
	evaluator.ParserFunc = func(input string) *ast.Program {
		l := lexer.New(input)
		return parser.New(l).ParseProgram()
	}

	// 3. Lexing & Parsing
	l := lexer.New(code)
	pObj := parser.New(l)
	program := pObj.ParseProgram()

	// Handle Parser Errors
	if len(pObj.Errors()) > 0 {
		var errs []interface{}
		for _, msg := range pObj.Errors() {
			errs = append(errs, "PARSER ERROR: "+msg)
		}
		return map[string]interface{}{
			"error": errs,
		}
	}

	// 4. Evaluation
	result := evaluator.Eval(program, env)

	// 5. Prepare Result
	finalResult := ""
	if result != nil && result.Type() != object.NULL_OBJ {
		finalResult = result.Inspect()
	}

	// Handle Runtime Errors
	if result != nil && result.Type() == object.ERROR_OBJ {
		return map[string]interface{}{
			"error": []interface{}{result.Inspect()},
		}
	}

	return map[string]interface{}{
		"logs":   outputBuffer.String(), // Captured "show()" output
		"result": finalResult,           // The return value of the script
	}
}

// overrideBuiltinsForWeb modifies the 'show' and 'ask' commands to work in browser
func overrideBuiltinsForWeb() {
	// Find and replace "show"
	for i, b := range object.Builtins {
		if b.Name == "show" {
			object.Builtins[i].Builtin = &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					var parts []string
					for _, arg := range args {
						parts = append(parts, arg.Inspect())
					}
					// Write to buffer instead of os.Stdout
					outputBuffer.WriteString(strings.Join(parts, " ") + "\n")
					return &object.Null{}
				},
			}
		}
		// Find and replace "ask" (Input)
		// We cannot pause execution in WASM easily, so we return a placeholder.
		if b.Name == "ask" {
			object.Builtins[i].Builtin = &object.Builtin{
				Fn: func(args ...object.Object) object.Object {
					outputBuffer.WriteString("[Input not supported in Web Demo]\n")
					return &object.String{Value: "mock_input"}
				},
			}
		}
	}
}
