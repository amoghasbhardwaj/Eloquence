// ==============================================================================================
// FILE: object/object.go
// ==============================================================================================
// PACKAGE: object
// PURPOSE: Defines the type system for the Eloquence language.
//          It provides the structures for all runtime values (Integers, Functions, Structs, etc.)
//          and the interfaces required to interact with them.
// ==============================================================================================

package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"eloquence/ast"
)

// ObjectType is a string alias for identifying the type of an object at runtime.
type ObjectType string

const (
	// Primitive Types
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	CHAR_OBJ    = "CHAR"
	NULL_OBJ    = "NULL"

	// Internal Control Flow Types
	RETURN_VALUE_OBJ = "RETURN_VALUE" // Wraps a return value to bubble up through the AST
	ERROR_OBJ        = "ERROR"        // Wraps a runtime error message

	// Composite Types
	FUNCTION_OBJ = "FUNCTION"
	ARRAY_OBJ    = "ARRAY"
	MAP_OBJ      = "MAP"

	// Memory Management
	POINTER_OBJ = "POINTER"

	// User-Defined Types
	STRUCT_DEF_OBJ  = "STRUCT_DEFINITION" // The blueprint (class)
	STRUCT_INST_OBJ = "STRUCT_INSTANCE"   // The concrete object (instance)

	// Builtin Functions
	BUILTIN_OBJ = "BUILTIN" // Builtin functions
)

// Object is the base interface that every value in Eloquence must implement.
type Object interface {
	Type() ObjectType // Returns the type constant
	Inspect() string  // Returns a string representation for display
}

// ==============================================================================================
// PRIMITIVE OBJECTS
// ==============================================================================================

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Char struct {
	Value rune
}

func (c *Char) Type() ObjectType { return CHAR_OBJ }
func (c *Char) Inspect() string  { return string(c.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "none" }

// ==============================================================================================
// INTERNAL WRAPPERS
// ==============================================================================================

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// ==============================================================================================
// COMPLEX OBJECTS
// ==============================================================================================

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // Closure: Holds the environment at definition time
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	return "takes(...) { ... }"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	parts := []string{}
	for _, el := range a.Elements {
		parts = append(parts, el.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(parts, ", "))
	out.WriteString("]")
	return out.String()
}

// ==============================================================================================
// MAP & HASHING SYSTEM
// ==============================================================================================

// HashKey is a distinct key for identifying objects in maps.
// It combines the object type and a unique 64-bit hash.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// HashPair connects the original key object with its stored value.
type HashPair struct {
	Key   Object
	Value Object
}

// Hashable interface allows an object to be used as a map key.
type Hashable interface {
	HashKey() HashKey
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: INTEGER_OBJ, Value: uint64(i.Value)}
}

func (b *Boolean) HashKey() HashKey {
	var v uint64
	if b.Value {
		v = 1
	}
	return HashKey{Type: BOOLEAN_OBJ, Value: v}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value)) // FNV-1a hash algorithm
	return HashKey{Type: STRING_OBJ, Value: h.Sum64()}
}

type Map struct {
	Pairs map[HashKey]HashPair
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(),
			pair.Value.Inspect(),
		))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// ==============================================================================================
// POINTERS
// ==============================================================================================

type Pointer struct {
	Name string       // The name of the variable being pointed to
	Env  *Environment // The specific scope where that variable lives
}

func (p *Pointer) Type() ObjectType { return POINTER_OBJ }
func (p *Pointer) Inspect() string  { return "pointing to " + p.Name }

// ==============================================================================================
// STRUCTS
// ==============================================================================================

type StructDefinition struct {
	Name   string
	Fields []string
}

func (sd *StructDefinition) Type() ObjectType { return STRUCT_DEF_OBJ }
func (sd *StructDefinition) Inspect() string {
	return "struct " + sd.Name
}

type StructInstance struct {
	Definition *StructDefinition
	Fields     map[string]Object
}

func (si *StructInstance) Type() ObjectType { return STRUCT_INST_OBJ }
func (si *StructInstance) Inspect() string {
	var out bytes.Buffer
	parts := []string{}
	for k, v := range si.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v.Inspect()))
	}
	out.WriteString(si.Definition.Name)
	out.WriteString("{")
	out.WriteString(strings.Join(parts, ", "))
	out.WriteString("}")
	return out.String()
}

// ==============================================================================================
// BUILTIN FUNCTIONS
// ==============================================================================================

type Builtin struct {
	Fn func(args ...Object) Object
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
