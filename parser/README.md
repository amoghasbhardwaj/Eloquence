<!-- ======================================================== -->
<!-- Parser Package README — Eloquence Programming Language -->
<!-- ======================================================== -->

<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-English--First%20Language-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Package-Parser-2d9cdb?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Stage-Semantic%20Analysis-111111?style=for-the-badge" />
</p>

---

# Parser Package  
## Eloquence Programming Language

The **Parser** package is the **semantic brain** of the Eloquence compiler.  

It takes a **stream of tokens** from the Lexer and builds the **Abstract Syntax Tree (AST)**. The parser understands grammar, operator precedence, and ambiguity resolution. It ensures that **English-first syntax** is translated into **structured logic**.

---

## Table of Contents

- [Overview](#overview)  
- [Folder Structure](#folder-structure)  
- [Parsing Strategy (Pratt Parsing)](#parsing-strategy-pratt-parsing)  
- [Precedence Handling](#precedence-handling)  
- [Lookahead & Ambiguity Resolution](#lookahead--ambiguity-resolution)  
- [Error Recovery](#error-recovery)  
- [Visual Flow of Parsing](#visual-flow-of-parsing)  
- [Testing & Verification](#testing--verification)  
- [Performance Benchmarks](#performance-benchmarks)  
- [How to Run Tests](#how-to-run-tests)  

---

## Overview

The Parser:

- Converts linear token streams into structured AST nodes  
- Validates grammar (e.g., blocks, loops, function calls)  
- Handles operator precedence automatically  
- Resolves ambiguities in English-first syntax  

**Example Input Tokens:**

```
[TOKEN_IDENT "x"] [TOKEN_IS "is"] [TOKEN_INT "5"] [TOKEN_ADDS "adds"] [TOKEN_INT "10"]
```

**Parsed AST:**

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
parser/
├── parser.go
├── parser_unit_test.go
├── parser_integration_test.go
├── parser_sanity_test.go
└── parser_benchmark_test.go
```

| File | Purpose |
|------|--------|
| `parser.go` | Main parsing logic, Pratt tables, parsing functions |
| `parser_unit_test.go` | Individual grammar rule validation |
| `parser_integration_test.go` | Nested blocks, recursion, struct parsing |
| `parser_sanity_test.go` | Robustness and edge case testing |
| `parser_benchmark_test.go` | Performance measurements |

---

## Parsing Strategy (Pratt Parsing)

Eloquence uses a **Pratt parser** (top-down operator precedence).

- **Prefix Functions:** Handle tokens that **start** an expression  
  e.g., `IDENT`, `INT`, `IF`, `!`  
- **Infix Functions:** Handle tokens that **connect** expressions  
  e.g., `adds`, `minus`, `times`, `(`  

**Benefits:**

- Modular design
- Easy to add new operators
- Natural handling of precedence

---

## Precedence Handling

Precedence ensures correct order of operations **without parentheses**:

```
LOWEST
EQUALS      (==)
LESSGREATER (>, <)
SUM         (+, -)
PRODUCT     (*, /)
PREFIX      (-X, !X)
CALL        (func())
INDEX       (arr[0])
```

**Example:**

```
5 adds 10 times 2
```

Parsed as:

```
5 adds (10 times 2)
```

No parentheses needed.

---

## Lookahead & Ambiguity Resolution

### Struct vs. Block

Braces `{}` serve both:

- **Blocks**: `if`, `while`  
- **Struct Definitions**: `User { ... }`  

**Parser Solution:** **3-Token Lookahead**

```
Current: IDENT ("User")
Next:    {
Peek:    : or }
```

- If peek indicates fields → parse as **Struct**  
- Else → parse as **Block**

This preserves **English-first syntax** without forcing new symbols.

---

## Error Recovery

- Errors are **collected**, not thrown immediately  
- Multiple syntax issues reported in a single pass  

```go
errors []string
```

- Example: missing closing brace or invalid operator

---

## Visual Flow of Parsing

```
Token Stream → Parser → AST Nodes

[ IDENT "x" ][ IS ][ INT "5" ][ ADDS ][ INT "10" ]
      │
      ▼
Parser applies Pratt functions:
      │
      ▼
AST Generated:
AssignmentStatement
├── Name: "x"
└── Value: InfixExpression
    ├── Left: 5
    ├── Operator: "adds"
    └── Right: 10
```

---

## Testing & Verification

| Test File | Focus | Pass Criteria |
|-----------|-------|---------------|
| `parser_unit_test` | Grammar | Each statement parses to correct AST node |
| `parser_integration_test` | Nested logic | Recursion, loops, structs parse correctly |
| `parser_sanity_test` | Robustness | Invalid code produces clean errors |

**Run all tests:**

```bash
go test -v ./parser
```

---

## Performance Benchmarks

```text
BenchmarkParser_SimpleAssignment-12    2000000    600 ns/op
BenchmarkParser_LargeProgram-12        5000       200000 ns/op
```

**Observation:** Parser can handle **thousands of lines** of code in milliseconds.

---

## How to Run Tests

Run all parser tests:

```bash
go test -v ./parser
```

Run benchmarks:

```bash
go test -bench=. ./parser
```

---

## Summary

The Parser:

- Builds AST from token streams  
- Resolves ambiguity in English-first syntax  
- Handles operator precedence automatically  
- Collects multiple syntax errors per pass  
- Feeds precise AST to the Evaluator and Compiler

It is the **semantic core** of Eloquence —  
where *linear words become structured logic*.