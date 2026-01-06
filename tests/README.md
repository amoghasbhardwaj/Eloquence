# System Tests - Eloquence Programming Language

This documentation covers the system-level integration tests (`system_test.go`) and benchmarks (`main_benchmark_test.go`). Unlike unit tests found within specific packages (token, lexer, ast, parser, evaluator), these tests verify the entire compiler pipeline working in unison.

They serve as the "Turing Completeness" verification for the language, ensuring that all components interact correctly to solve actual computational problems.

## Table of Contents

1.  Overview
2.  Test Categories
3.  Key Test Scenarios
4.  Benchmark Performance
5.  How to Run

---

## Overview

System tests execute the full lifecycle of an Eloquence program:
1.  **Source Input**: Raw string representing code.
2.  **Lexing**: Converted to tokens.
3.  **Parsing**: Built into an AST.
4.  **Evaluation**: Executed against a live Environment.
5.  **Assertion**: The final Object result is compared against Go native types.

These tests detect regressions that might pass individual unit tests but fail when components interact (e.g., scoping issues between Parser and Evaluator).

---

## Test Categories

The system tests are divided into four specific architectural pillars:

### 1. Algorithmic Logic
Verifies recursion, arithmetic precedence, and conditional branching.
*   **Goal**: Ensure the language can solve mathematical problems.

### 2. Higher-Order Functions
Verifies first-class functions, closures, and passing functions as arguments.
*   **Goal**: Ensure functional programming capabilities (Map/Reduce patterns).

### 3. Data Structures
Verifies Struct definitions, instantiation, and recursive type references.
*   **Goal**: Ensure complex data modeling (e.g., Linked Lists, Trees).

### 4. Memory Management
Verifies the manual pointer system (`pointing to`, `pointing from`) and variable shadowing rules.
*   **Goal**: Ensure reference semantics work alongside value semantics.

---

## Key Test Scenarios

The following scenarios from `system_test.go` illustrate the language capabilities being verified.

### Fibonacci Sequence (Recursion)
Tests stack depth, function calls, and integer arithmetic.

    fib is takes(x)
        if x less 2
            return x
        end
        return fib(x minus 1) adds fib(x minus 2)
    end
    fib(10)

### Linked List Traversal (Structs)
Tests struct field access and handling of `none` (null) values.

    define Node as struct { val, next }
    head is Node { val: 10, next: none }
    
    // Traversal logic
    sumList is takes(node)
        if node equals none
            return 0
        end
        return node.val adds sumList(node.next)
    end

### Pointer Mutation
Tests that a pointer can modify a variable defined in an outer scope.

    globalVal is 100
    ptr is pointing to globalVal
    pointing from ptr is 999
    // globalVal is now 999

---

## Benchmark Performance

Benchmarks in `main_benchmark_test.go` stress-test the pipeline under heavy load.

| Benchmark | Description | Target Metric |
| :--- | :--- | :--- |
| **HeavyLoop** | Runs a loop 1000 times performing addition. | Iteration overhead. |
| **DeepRecursion** | Recursively calls a function 200 times. | Stack frame allocation speed. |
| **StringConcat** | Concatenates strings in a loop. | Heap allocation/GC pressure. |

---

## How to Run

### Run System Tests
To execute only the system-level verification:

    go test -v system_test.go

### Run Benchmarks
To measure the performance of the full interpreter pipeline:

    go test -bench=. main_benchmark_test.go main.go