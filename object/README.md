# Object System - Eloquence Programming Language

The **object** package defines the runtime representation of all data values within the Eloquence language. In Eloquence, every value—from simple integers to complex functions and structs—is an `Object` that implements a common interface.

This package also implements the **Environment**, which handles variable storage, scoping rules, and memory management during program execution.

## Table of Contents

1.  Overview
2.  Folder Structure
3.  Core Architecture
4.  Data Types
5.  Memory & Scoping (The Environment)
6.  Hashing System
7.  Testing Strategy
8.  Performance Benchmarks
9.  How to Run Tests

---

## Overview

The Object System provides a unified type hierarchy for the interpreter. It allows the `Evaluator` to treat distinct data types polymorphically.

**Key Responsibilities:**
*   **Type Identity**: Every object knows its own type (`INTEGER`, `BOOLEAN`, `FUNCTION`, etc.).
*   **Inspection**: Every object can serialize itself back into a readable string format for debugging and output.
*   **State Management**: The Environment struct manages the lifecycle of these objects across different scopes (global, local, closure).

---

## Folder Structure

*   **object.go**: Definitions of all object structs (`Integer`, `String`, `Array`, `Function`, etc.) and the `Object` interface.
*   **environment.go**: Implementation of the `Environment` struct, including `Get`, `Set`, and scope extension logic.
*   **object_unit_test.go**: Verifies `Inspect()` output and `Type()` identification for all object variants.
*   **object_integration_test.go**: Tests complex interactions like storing structs in environments and using objects as map keys.
*   **object_sanity_test.go**: Checks edge cases like empty collections and deeply nested environments to prevent stack overflows.
*   **object_benchmark_test.go**: Measures the performance of object creation, hashing, and environment lookups.
*   **environment_unit_test.go**: Specifically validates variable shadowing and scope traversal rules.

---

## Core Architecture

### The Object Interface

Every value in Eloquence must implement this interface:

    type Object interface {
        Type() ObjectType  // Returns the internal type constant (e.g., INTEGER_OBJ)
        Inspect() string   // Returns a string representation (e.g., "10", "true")
    }

---

## Data Types

### Primitive Types
*   **Integer**: 64-bit signed integers.
*   **Float**: 64-bit floating-point numbers.
*   **Boolean**: True/False values.
*   **String**: UTF-8 string literals.
*   **Char**: Single UTF-8 characters (runes).
*   **Null**: Represents the absence of a value (`none`).

### Composite Types
*   **Array**: An ordered list of objects.
*   **Map**: A hash map supporting Integers, Strings, and Booleans as keys.
*   **Struct**: Custom user-defined data structures (`StructDefinition` and `StructInstance`).

### Internal Types
*   **ReturnValue**: A wrapper used to transport values out of nested block statements.
*   **Error**: Represents runtime errors to stop execution.
*   **Function**: First-class functions supporting closures (capturing their definition environment).
*   **Pointer**: Holds a reference to a variable name and its specific environment scope.

---

## Memory & Scoping (The Environment)

The `Environment` is a hash map string-to-object store that supports lexical scoping via a linked-list structure.

    type Environment struct {
        store map[string]Object // Local variables
        outer *Environment      // Pointer to the parent scope
    }

### Scoping Rules
1.  **Global Scope**: The root environment.
2.  **Local Scope**: Created for function calls and blocks.
3.  **Shadowing**: Writing to a variable in an inner scope creates a new variable, protecting the outer scope (unless pointers are used).
4.  **Lookup**: `Get(name)` searches the current scope; if not found, it recursively checks the `outer` scope.

---

## Hashing System

To support Maps, objects that can be used as keys must implement the `Hashable` interface.

    type Hashable interface {
        HashKey() HashKey
    }

This uses the **FNV-1a** hash algorithm for strings and direct value mapping for integers/booleans to ensure O(1) map access speeds.

---

## Testing Strategy

| Test Suite | Focus | Pass Criteria |
| :--- | :--- | :--- |
| `object_unit_test` | Representation | `Inspect()` returns the exact string format expected by the REPL. |
| `environment_unit_test` | Scoping | Inner scopes can read outer variables but cannot overwrite them without explicit shadowing. |
| `object_integration_test` | Complexity | Structs can be stored in variables and retrieved with fields intact. |
| `object_benchmark_test` | Latency | Environment lookups remain fast even with 50+ layers of nesting. |

---

## Performance Benchmarks

    BenchmarkEnvironment_Get_Deep-12    50000000    30 ns/op
    BenchmarkHashKey_String-12          30000000    40 ns/op

*   **Deep Lookup**: Even with deep recursion, variable access remains nanosecond-scale.
*   **Hashing**: String hashing is optimized to allow fast map lookups.

---

## How to Run Tests

    # Run all object system tests
    go test -v ./object

    # Run benchmarks
    go test -bench=. ./object