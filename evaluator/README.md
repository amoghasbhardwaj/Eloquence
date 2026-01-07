<!-- =========================================================== -->
<!-- Evaluator Package README — Eloquence Programming Language -->
<!-- =========================================================== -->

<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-English--First%20Language-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Package-Evaluator-6fcf97?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Stage-Execution%20Engine-111111?style=for-the-badge" />
</p>

---

# Evaluator Package  
## Eloquence Programming Language

The **Evaluator** is the **runtime engine** of Eloquence.  

It takes the **Abstract Syntax Tree (AST)** produced by the Parser and executes it. Its responsibilities include:

- Variable storage and management  
- Arithmetic and logical computation  
- Control flow execution (`if`, `while`, loops)  
- Function calls and closures  
- Integration with the **Object System**

The Evaluator is essentially **where code comes alive**.

---

## Table of Contents

1. [Overview](#1-overview)  
2. [Folder Structure](#2-folder-structure)  
3. [Execution Model](#3-execution-model-tree-walking)  
4. [Object System Integration](#4-object-system-integration)  
5. [Environment & Scope](#5-environment--scope)  
6. [Built-in Functions](#6-built-in-functions)  
7. [Error Handling](#7-error-handling)  
8. [Testing Strategy](#8-testing-strategy)  
9. [Running Tests](#9-running-tests)  

---

## 1. Overview

The Evaluator **traverses the AST** recursively and executes statements and expressions.

**Example:**

Source code:

```
x is 5 adds 10
```

Evaluation steps:

```
1. Eval visits AssignmentStatement
2. Eval recursively on Left node → Identifier "x"
3. Eval recursively on Right node → InfixExpression
4. Eval Left → 5
5. Eval Right → 10
6. Apply operator "adds" → 15
7. Store result in Environment under "x"
```

---

## 2. Folder Structure

```
evaluator/
├── evaluator.go
├── evaluator_test.go
└── evaluator_integration_test.go
```

| File | Purpose |
|------|---------|
| `evaluator.go` | Main evaluation logic, tree-walking interpreter |
| `evaluator_test.go` | Unit tests for arithmetic, logic, and helper functions |
| `evaluator_integration_test.go` | Integration tests for recursion, closures, structs, pointers |

---

## 3. Execution Model (Tree Walking)

The evaluator implements a **Tree-Walking Interpreter**:

```
AST Node
   │
   ▼
Eval(Node)
   │
   ├─ If Statement → Eval Condition → Eval Consequence
   │
   ├─ Expression → Eval Left → Eval Right → Apply Operator
   │
   └─ Function → Create Enclosed Environment → Eval Body
```

**Example Evaluation:**

```
5 adds 5
```

```
Eval(InfixExpression)
├─ Eval(Left) → 5
├─ Eval(Right) → 5
└─ Apply operator "adds" → returns 10
```

---

## 4. Object System Integration

All runtime values are **objects** implementing the `object.Object` interface.

Supported object types:

| Primitive | Composite | Special |
|-----------|-----------|---------|
| Integer   | Array     | Function |
| Float     | Map       | Pointer |
| Boolean   | StructInstance | ReturnValue |
| String    | StructDefinition | Error |
| Null      |           | Builtin |

This allows **dynamic typing** and **polymorphic evaluation**.

---

## 5. Environment & Scope

The **Environment** is a nested hash map representing variable scopes:

```
type Environment struct {
    store map[string]Object
    outer *Environment
}
```

### Key Mechanics:

- **Lexical Scope:** Functions create enclosed environments
- **Closures:** Functions capture the environment at definition time
- **Pointers:** Enable direct mutation across scopes

**Illustration:**

```
Global Env
 ├─ x: 10
 └─ func add(y)
      └─ Enclosed Env
           ├─ y: 5
           └─ Pointer to x
```

---

## 6. Built-in Functions

| Function | Purpose |
|----------|---------|
| `show(...)` | Prints to console |
| `count(x)` | Returns length of array/string |
| `append(arr, val)` | Adds an element to array |
| `upper(s)` / `lower(s)` | String case conversion |
| `split(s, sep)` / `join(arr, sep)` | String-array utilities |

---

## 7. Error Handling

Runtime errors propagate using the `object.Error` type:

- Division by zero  
- Type mismatches  
- Undefined identifiers  

**Example:**

```
5 adds "hello"
→ Error: type mismatch: INTEGER + STRING
```

Execution stops gracefully, allowing the program to report **user-friendly messages**.

---

## 8. Testing Strategy

| Test Suite | Focus | Pass Criteria |
|-----------|-------|---------------|
| `evaluator_test.go` | Unit logic | Arithmetic, boolean, if/else evaluation |
| `evaluator_integration_test.go` | Integration | Recursion, closures, pointers, structs |

**Unit Test Example:**

```go
func TestEvalIntegerExpression(t *testing.T) {
    input := "5 adds 5"
    evaluated := Eval(input, NewEnvironment())
    if evaluated.Value != 10 {
        t.Errorf("expected 10, got %v", evaluated.Value)
    }
}
```

---

## 9. Running Tests

Run all tests:

```bash
go test -v ./evaluator
```

Run integration tests:

```bash
go test -v ./evaluator -run Integration
```

---

### Summary

The Evaluator:

- Traverses AST nodes recursively  
- Integrates deeply with the **Object System**  
- Manages nested **Environments and closures**  
- Executes **English-first syntax** correctly  
- Supports **runtime errors gracefully**  
- Powers the **core execution engine** of Eloquence