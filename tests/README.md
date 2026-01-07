# System Tests  
## Eloquence Programming Language

This documentation covers the **system-level integration tests** (`system_test.go`) and **benchmarks** (`main_benchmark_test.go`).

Unlike unit tests within specific packages, these tests verify the **entire compiler pipeline** in unison. They act as a "Turing Completeness" verification, ensuring all components interact correctly to solve real computational problems.

---

## Table of Contents

1. [Overview](#1-overview)  
2. [Test Categories](#2-test-categories)  
3. [Key Test Scenarios](#3-key-test-scenarios)  
4. [Benchmark Performance](#4-benchmark-performance)  
5. [How to Run](#5-how-to-run)  

---

## 1. Overview

System tests execute the **full lifecycle of an Eloquence program**:

```
Source Code → Lexer → Parser → AST → Evaluator → Object Result → Assertion
```

Steps:

1. **Source Input**: Raw Eloquence code string.  
2. **Lexing**: Converted into a token stream.  
3. **Parsing**: Tokens are structured into an AST.  
4. **Evaluation**: AST is executed in a live Environment.  
5. **Assertion**: The resulting Object is compared against Go native types.  

These tests catch regressions that may pass unit tests but fail during component interaction, such as **scope resolution or pointer behavior**.

---

## 2. Test Categories

System tests are organized into four architectural pillars:

| Category | Focus | Goal |
|----------|-------|------|
| Algorithmic Logic | Recursion, arithmetic, conditional branching | Ensure Eloquence solves math problems correctly |
| Higher-Order Functions | First-class functions, closures, function arguments | Support functional programming patterns (Map/Reduce) |
| Data Structures | Structs, instantiation, recursion | Enable complex user-defined data modeling |
| Memory Management | Pointers (`pointing to/from`), variable shadowing | Ensure reference semantics work correctly alongside value semantics |

---

## 3. Key Test Scenarios

### 3.1 Fibonacci Sequence (Recursion)

Tests stack depth, function calls, and arithmetic.

```eloquence
fib is takes(x) {
    if x less 2 {
        return x
    }
    return fib(x minus 1) adds fib(x minus 2)
}
fib(10)  // Expected: 55
```

---

### 3.2 Linked List Traversal (Structs)

Tests struct field access and handling of `none` values.

```eloquence
define Node as struct { val, next }

head is Node { val: 10, next: none }

// Traversal logic
sumList is takes(node) {
    if node equals none {
        return 0
    }
    return node.val adds sumList(node.next)
}
sumList(head)  // Expected: 10
```

---

### 3.3 Pointer Mutation

Tests that pointers can modify variables in an outer scope.

```eloquence
globalVal is 100
ptr is pointing to globalVal
pointing from ptr is 999
// globalVal is now 999
```

---

## 4. Benchmark Performance

Benchmarks in `main_benchmark_test.go` stress-test the **full pipeline**:

| Benchmark | Description | Target Metric |
|-----------|------------|---------------|
| HeavyLoop | Loop 1000 iterations performing addition | Iteration overhead, environment lookup speed |
| DeepRecursion | Recursively calls a function 200 times | Stack frame allocation speed, memory safety |
| StringConcat | Concatenates strings in a loop | Heap allocation efficiency, GC pressure |

---

## 5. How to Run

### Run System Tests

```bash
go test -v system_test.go main.go
```

### Run Benchmarks

```bash
go test -bench=. main_benchmark_test.go main.go
```

---

### Summary

The **System Tests** ensure:

- The **entire compiler pipeline** works correctly from source to evaluation.  
- Recursive functions, closures, structs, and pointers behave correctly.  
- Performance under loops, recursion, and memory-intensive operations meets design expectations.  
- The language is robust, human-readable, and ready for real-world use.