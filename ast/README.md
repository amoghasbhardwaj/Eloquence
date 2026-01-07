<!-- ===================================================== -->
<!-- AST Package README — Eloquence Programming Language -->
<!-- ===================================================== -->

<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-English--First%20Language-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Package-AST-bb6bd9?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Stage-Structural%20Representation-111111?style=for-the-badge" />
</p>

---

# AST Package  
## Eloquence Programming Language

The **AST** (Abstract Syntax Tree) package is the **structural backbone** of the Eloquence compiler.

It transforms a **linear token stream** into a **hierarchical tree**, where each node represents a semantic unit—from literals to full control flow structures.

By abstracting raw code into a tree:

- The **Evaluator** can compute values efficiently  
- The **Compiler** can reason about program logic  
- English-first syntax is preserved in a structured format

---

## Table of Contents

- [Overview](#overview)  
- [Folder Structure](#folder-structure)  
- [Core Architecture](#core-architecture)  
- [Language Constructs](#language-constructs)  
  - [Statements](#statements)  
  - [Expressions](#expressions)  
- [Visual Flow of the AST](#visual-flow-of-the-ast)  
- [Testing & Verification](#testing--verification)  
- [Performance Benchmarks](#performance-benchmarks)  
- [How to Run Tests](#how-to-run-tests)  

---

## Overview

If **Tokens** are "words," the AST represents **sentences and paragraphs**.

It enforces **syntax structure**, for example:

- `if` statements must have a **Condition** and a **Consequence Block**  
- `Functions` must include parameters and a body  
- Operators are associated with operands

### Example

Source code:

```eloquence
x is 5 adds 10
```

AST (conceptual):

```
      AssignmentStatement
      ├── Name: "x"
      └── Value: InfixExpression
          ├── Left: 5
          ├── Operator: "adds"
          └── Right: 10
```

---

## Folder Structure

```
ast/
├── ast.go
├── ast_unit_test.go
├── ast_integration_test.go
├── ast_sanity_test.go
└── ast_benchmark_test.go
```

| File | Purpose |
|------|--------|
| `ast.go` | Node interfaces and struct definitions |
| `ast_unit_test.go` | Basic literals and node tests |
| `ast_integration_test.go` | Nested structures like Functions and Structs |
| `ast_sanity_test.go` | Deep recursion stress tests |
| `ast_benchmark_test.go` | String reconstruction speed |

---

## Core Architecture

The AST is **strictly typed** with three main interfaces.

| Interface | Role | Key Methods |
|-----------|------|-------------|
| **Node** | Base of all AST nodes | `TokenLiteral()`, `String()` |
| **Statement** | Nodes performing actions | `statementNode()` |
| **Expression** | Nodes evaluating to values | `expressionNode()` |

### Root Node: `Program`

The `Program` node is the entry point:

- Contains all top-level `Statements`  
- Represents the **full source code**  

---

## Language Constructs

### Statements

Statements define **program structure**.

| Node Type | Syntax Example | Purpose |
|-----------|----------------|--------|
| AssignmentStatement | `x is 10` | Bind value to a variable |
| ReturnStatement | `return 10` | Exit function with a value |
| LoopStatement | `while x < 10 { ... }` | Iterative control flow |
| StructDefinition | `define Node as struct` | User-defined data type |
| ExpressionStatement | `show(x)` | Wraps standalone expressions |

### Expressions

Expressions are **computational units**.

| Node Type | Syntax Example | Purpose |
|-----------|----------------|--------|
| InfixExpression | `x adds y` | Binary arithmetic or logical operations |
| PointerReference | `pointing to x` | Memory address reference |
| FunctionLiteral | `takes (x) { ... }` | Anonymous function definition |
| CallExpression | `calculate(5,10)` | Function invocation |

---

## Visual Flow of the AST

**Input:** `x is 5 adds 10`

```
      +-----------------------+
      |  AssignmentStatement  |
      +-----------+-----------+
                  |
       +----------+----------+
       |                     |
 [ Identifier ]      [ InfixExpression ]
   Literal: "x"              |
                    +--------+--------+
                    |        |        |
                [Left]    [Op]    [Right]
                  5      "adds"     10
```

---

## Testing & Verification

The AST is heavily tested to ensure correctness and stability.

| Test File | Focus | Pass Criteria |
|-----------|-------|---------------|
| `ast_unit_test` | Literals | `.String()` matches token literal |
| `ast_integration_test` | Nesting | Recursive structures serialize correctly |
| `ast_sanity_test` | Stability | Handles 100+ levels of nesting without stack overflow |

**Run all tests:**

```bash
go test -v ./ast
```

---

## Performance Benchmarks

Sample results:

```text
BenchmarkInfixExpressionString-12    40000000    30 ns/op
```

**Observation:**  
String reconstruction is highly optimized for fast debugging and error reporting.

---

## How to Run Tests

Run all AST tests:

```bash
go test -v ./ast
```

Run benchmarks:

```bash
go test -bench=. ./ast
```

---

## Summary

The AST package:

- Converts linear token streams into hierarchical trees  
- Enforces English-first semantic structure  
- Supports statements and expressions  
- Handles deep nesting and recursion efficiently  
- Provides a foundation for the Evaluator and Compiler

It is the **structural core** of the Eloquence language —  
the point where *words become meaningful logic*.