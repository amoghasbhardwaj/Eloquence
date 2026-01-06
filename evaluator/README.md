# Evaluator - Eloquence Programming Language

The **evaluator** package is the runtime engine of the Eloquence compiler. It executes the logic defined by the Abstract Syntax Tree (AST) produced by the parser.

This implementation uses a **tree-walking interpreter** strategy. It recursively traverses the AST nodes, evaluating expressions and executing statements in the context of a mutable environment.

---

## Table of Contents

1.  Overview
2.  Folder Structure
3.  Core Architecture
4.  Language Features Supported
5.  Testing Strategy
6.  Performance Benchmarks
7.  How to Run Tests

---

## Overview

The Evaluator is the "brain" of the language. It transforms static syntax trees into dynamic behavior. Its responsibilities include:

*   **Expression Evaluation**: Computing the results of arithmetic, logic, and function calls.
*   **State Management**: Maintaining variable scope and binding values to identifiers via the `Environment`.
*   **Control Flow**: Handling loops (`for`/`while`), conditionals (`if`/`else`), and return statements.
*   **Object System Integration**: Wrapping all computed values in the `object.Object` interface (e.g., `Integer`, `Boolean`, `Function`).

---

## Folder Structure

The package is organized to separate core logic from testing suites:

*   **evaluator.go**: The core execution logic. Contains the `Eval` function and handlers for all AST node types.
*   **evaluator_unit_test.go**: Granular tests for specific expressions (math, logic) and basic statements.
*   **evaluator_integration_test.go**: End-to-end tests for complex features like closures, recursion, and struct manipulation.
*   **evaluator_sanity_test.go**: Validation for edge cases, empty programs, and runtime error reporting.
*   **evaluator_benchmark_test.go**: Performance benchmarks for heavy computations.

---

## Core Architecture

### The Eval Function

The entry point is the `Eval(node, env)` function. It uses a type switch to dispatch execution based on the AST node type.

    func Eval(node ast.Node, env *object.Environment) object.Object {
        switch node := node.(type) {
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.InfixExpression:
            // ... recurse left and right
        // ... other types
        }
    }

### The Environment

The evaluator passes an `env` object through every call. This environment acts as the memory heap/stack, storing variable bindings.

*   **Block Scope**: New environments are created for blocks (e.g., inside `if` statements) to prevent variable leakage (shadowing).
*   **Closure Support**: Functions capture the environment in which they were defined, allowing inner functions to access variables from their defining scope.
*   **Loop Scoping**: Unlike `if` blocks, loops in Eloquence share the parent scope to allow updating external counters seamlessly.

---

## Language Features Supported

### 1. Arithmetic & Logic
Standard operators (`adds`, `subtracts`, `times`, `divides`) and comparisons (`equals`, `greater`) are fully implemented for Integers, Floats, and Booleans.

### 2. Control Structures
*   **If/Else**: Evaluates conditions and executes the appropriate block.
*   **Loops**: Supports standard iterative flows.
    *   *Feature*: Includes interruption logic to handle `return` statements inside loops.

### 3. Functions & Closures
Functions are first-class citizens. They can be assigned to variables, passed as arguments, and returned from other functions.

### 4. Data Structures
*   **Arrays**: Literal definition `[1, 2]` and index access `arr[0]`.
*   **Maps**: Key-value pairs `{"a": 1}` and hash-based access.
*   **Structs**: Custom type definitions (`define User as struct`), instantiation, and dot-notation field access (`user.name`).

### 5. Pointers
Direct memory manipulation via:
*   `pointing to x`: Creates a reference.
*   `pointing from ptr`: Dereferences or mutates the value at the specific environment address.

---

## Testing Strategy

The evaluator is critical infrastructure and is tested extensively to ensure correctness and stability.

| Test Suite | Focus | Pass Criteria |
| :--- | :--- | :--- |
| `unit_test` | Basic Math/Logic | `5 adds 5` returns `10`; `true equals false` returns `false`. |
| `integration_test` | Complex Flows | Recursive factorials, closure counters, and struct logic work as expected. |
| `sanity_test` | Error Handling | Invalid operations (e.g., `5 adds true`) return readable error messages instead of panicking. |
| `benchmark_test` | Speed | Recursive Fibonacci and large loop operations execute within acceptable timeframes. |

---

## Performance Benchmarks

Benchmarks measure the overhead of the interpreter loop and memory allocation.

    BenchmarkEvaluator_Fibonacci-12        50000     30000 ns/op
    BenchmarkEvaluator_LargeArraySum-12    20000     60000 ns/op

*   **Recursion**: Heavily exercises stack frame creation and environment switching.
*   **Loops**: Tests variable lookup speed and repeated statement execution.

---

## How to Run Tests

To execute the full test suite with verbose output:

    go test -v ./evaluator

To run the performance benchmarks:

    go test -bench=. ./evaluator