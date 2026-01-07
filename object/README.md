<!-- ======================================================== -->
<!-- Object Package README — Eloquence Programming Language -->
<!-- ======================================================== -->

<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-English--First%20Language-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Package-Object%20System-6fcf97?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Stage-Runtime%20Representation-111111?style=for-the-badge" />
</p>

---

# Object System  
## Eloquence Programming Language

The **object** package defines the **runtime representation** of all data values in Eloquence.  

Every value—from simple integers to functions, structs, and pointers—is an `Object` implementing a common interface.  

It also provides the **Environment** to handle variable storage, scoping, memory management, and closure state.

---

## Table of Contents

1. [Overview](#1-overview)  
2. [Folder Structure](#2-folder-structure)  
3. [Core Architecture](#3-core-architecture)  
4. [Data Types](#4-data-types)  
5. [Memory & Scoping (The Environment)](#5-memory--scoping-the-environment)  
6. [Hashing System](#6-hashing-system)  
7. [Testing Strategy](#7-testing-strategy)  
8. [Performance Benchmarks](#8-performance-benchmarks)  
9. [How to Run Tests](#9-how-to-run-tests)  

---

## 1. Overview

The Object System provides a unified type hierarchy for the interpreter.  

Key responsibilities:

- **Type Identity**: Every object knows its type constant (`INTEGER`, `BOOLEAN`, `FUNCTION`, etc.).  
- **Inspection**: Objects can serialize themselves for debugging, REPL output, and errors.  
- **State Management**: `Environment` tracks object lifecycles, scopes, and pointer references.

---

## 2. Folder Structure

```
object/
├── object.go
├── builtins.go
├── environment.go
├── object_unit_test.go
├── object_integration_test.go
├── object_sanity_test.go
├── object_benchmark_test.go
└── environment_unit_test.go
```

| File | Purpose |
|---|---|
| `object.go` | Definitions of `Object` interface & data structs (Integer, Function, etc.) |
| `builtins.go` | Standard library (`show`, `append`, `len`) |
| `environment.go` | Variable storage (`Get`/`Set`), scope extension, pointer resolution |
| `object_unit_test.go` | Verifies `Inspect()` output & type constants |
| `object_integration_test.go` | Tests complex interactions (Maps, Struct nesting) |
| `environment_unit_test.go` | Validates scoping, shadowing, closures |
| `object_benchmark_test.go` | Measures hashing & environment lookup performance |

---

## 3. Core Architecture

### Object Interface

```go
type Object interface {
    Type() ObjectType  // Returns internal type constant
    Inspect() string   // Returns string representation (e.g., "10", "true")
}
```

This allows the Evaluator to manipulate `[]Object` or `map[string]Object` without caring about the concrete Go type.

---

## 4. Data Types

### Primitive Types

- **Integer**: 64-bit signed integers  
- **Float**: 64-bit floating-point numbers  
- **Boolean**: `true` / `false`  
- **String**: UTF-8 string literals  
- **Char**: Single UTF-8 rune  
- **Null**: Represents absence of value (`none`)  

### Composite Types

- **Array**: Ordered list of Objects  
- **Map**: Hash map (Integer, String, Boolean keys)  
- **StructDefinition**: Blueprint of a struct  
- **StructInstance**: Concrete instance containing data fields  

### Internal Types

- **ReturnValue**: Wraps return values for block statements  
- **Error**: Runtime error objects (division by zero, type mismatch)  
- **Function**: First-class function with closure environment  
- **Pointer**: Reference to variable & scope for pass-by-reference  

### Built-in Functions

- **Builtin**: Wraps native Go functions callable from Eloquence (`show`, `append`, etc.)

---

## 5. Memory & Scoping (The Environment)

```go
type Environment struct {
    store map[string]Object // Local variables
    outer *Environment      // Parent scope
}
```

### Rules

1. **Global Scope**: Root environment at program start  
2. **Local Scope**: Created in functions or blocks  
3. **Lookup (`Get`)**: Checks current store, then recursively outer scopes  
4. **Shadowing**: New variable in inner scope preserves outer value  
5. **Mutation**: `Pointer` object references specific environment for updates  

---

## 6. Hashing System

To support Maps, hashable objects implement:

```go
type Hashable interface {
    HashKey() HashKey
}

type HashKey struct {
    Type  ObjectType
    Value uint64
}
```

### Algorithm

- **Integer/Boolean**: Direct value  
- **String**: FNV-1a non-cryptographic hash  

Ensures O(1) map access and minimal collisions.

---

## 7. Testing Strategy

| Test Suite | Focus | Pass Criteria |
|---|---|---|
| `object_unit_test` | Representation | `Inspect()` returns correct string |
| `environment_unit_test` | Scoping | Inner scopes read outer variables; shadowing correct |
| `object_integration_test` | Complexity | Structs, arrays, maps maintain integrity |
| `object_sanity_test` | Edge Cases | Deeply nested structures and empty collections don't crash |

---

## 8. Performance Benchmarks

```
BenchmarkEnvironment_Get_Deep-12    50000000    30 ns/op
BenchmarkHashKey_String-12          30000000    40 ns/op
```

**Observations:**

- Deep nested lookups remain nanosecond-scale  
- String hashing optimized for high-performance map access  

---

## 9. How to Run Tests

Run all object system tests:

```bash
go test -v ./object
```

Run benchmarks:

```bash
go test -bench=. ./object
```