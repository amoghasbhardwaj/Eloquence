// ==============================================================================================
// FILE: object/environment.go
// ==============================================================================================
// PACKAGE: object
// PURPOSE: Implements the memory environment (symbol table) for the interpreter.
//          It handles variable storage, lexical scoping chains, and shadowing logic.
// ==============================================================================================

package object

type Environment struct {
	store map[string]Object // Storage for the current scope
	outer *Environment      // Link to the enclosing (outer) scope
}

// NewEnvironment creates a fresh global environment.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a new local scope linked to an outer scope.
// This is used for functions and blocks to implement lexical scoping.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a value associated with a name.
// It searches the current scope first, then recursively checks outer scopes.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set stores a value in the CURRENT scope.
// If the variable exists in an outer scope, this creates a new "shadow" variable
// in the current scope, preserving the outer variable's original value.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Resolve finds the specific environment instance where a variable is defined.
// This is used by Pointers to bypass shadowing and modify variables in their original scope.
func (e *Environment) Resolve(name string) *Environment {
	if _, ok := e.store[name]; ok {
		return e
	}
	if e.outer != nil {
		return e.outer.Resolve(name)
	}
	return nil
}
