<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-English--First%20Language-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Package-Token-eb5757?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Stage-Lexical%20Analysis-111111?style=for-the-badge" />
</p>

# Token Package
## Eloquence Programming Language

The **token package** defines the *lexical building blocks* of the Eloquence programming language.

It is responsible for transforming raw source text into **typed symbols** that the compiler can reason about. This package sits at the boundary between *plain text* and *language semantics*, enabling the Lexer and Parser to communicate using a precise, standardized vocabulary.

---

## Table of Contents

- [Conceptual Overview](#conceptual-overview)
- [Position in the Compiler Pipeline](#position-in-the-compiler-pipeline)
- [Folder Structure](#folder-structure)
- [Core Architecture](#core-architecture)
  - [Token Types](#token-types)
  - [Keyword Resolution](#keyword-resolution)
- [English-First Syntax Mapping](#english-first-syntax-mapping)
- [Token Data Model](#token-data-model)
- [Performance Considerations](#performance-considerations)
- [Testing Strategy](#testing-strategy)
- [Benchmarks](#benchmarks)
- [Running Tests](#running-tests)

---

## Conceptual Overview

In human terms, source code is a sentence.  
In compiler terms, source code is just **characters**.

The token package defines how those characters become **meaningful units**.

Example source:

    x adds 5

After tokenization:

    IDENT("x")
    ADDS
    INT("5")

Each word or symbol becomes a **Token** with:
- a semantic type
- a literal value
- a precise source location

This is the **first semantic step** in the Eloquence compiler.

---

## Position in the Compiler Pipeline

The token package is not used directly by users. It exists to support the compiler pipeline.

Pipeline overview:

    Source Code
        ↓
      Lexer
        ↓
    Token Stream  ←── token package defines this
        ↓
      Parser
        ↓
        AST
        ↓
     Evaluator

Without tokens, the parser would have no structured input.

---

## Folder Structure

    token/
    ├── token.go
    ├── token_unit_test.go
    ├── token_integration_test.go
    ├── token_sanity_test.go
    ├── token_edge_test.go
    └── token_benchmark_test.go

### File Responsibilities

| File | Responsibility |
|----|----|
| token.go | Token types, constants, keyword map |
| token_unit_test.go | Verifies keyword → token mappings |
| token_integration_test.go | Groups tokens by language feature |
| token_sanity_test.go | Full program token streams |
| token_edge_test.go | Case sensitivity and boundary inputs |
| token_benchmark_test.go | Performance of LookupIdent |

---

## Core Architecture

### Token Types

Token categories are represented using a string-based enum:

    type TokenType string

This choice prioritizes **clarity over micro-optimizations**.

Advantages:
- Readable parser error messages
- Self-documenting debug output
- No integer → string mapping required

Example error:

    expected IDENT, got RETURN

This is immediately understandable during development.

---

### Keyword Resolution

The central logic is keyword lookup:

    func LookupIdent(ident string) TokenType

Resolution strategy:

1. Check identifier against a keyword map
2. If found, return keyword token
3. Otherwise, default to IDENT

| Property | Value |
|------|------|
| Data structure | Go map |
| Time complexity | O(1) |
| Allocations | None |
| Execution frequency | Every identifier |

Design note:  
Built-ins like `show` are **not** keywords. They remain identifiers to allow user overrides and future extensibility.

---

## English-First Syntax Mapping

Eloquence replaces symbolic syntax with readable English phrases.  
The token package defines how these phrases map to internal representations.

### Arithmetic & Assignment

| Concept | Traditional | Eloquence | Token |
|------|------|------|------|
| Assignment | = | is | IS |
| Addition | + | adds | ADDS |
| Subtraction | - | minus | MINUS |
| Multiplication | * | times | TIMES |

### Comparison & Logic

| Concept | Eloquence | Token |
|------|------|------|
| Equality | equals | EQUALS |
| Inequality | not_equals | NOT_EQUALS |
| Logical AND | and | AND |
| Logical OR | or | OR |
| Negation | not | NOT |

### Control & Structure

| Feature | Eloquence | Token |
|------|------|------|
| Function | takes | TAKES |
| Import | include | INCLUDE |
| Loop | for … in | FOR / IN |
| Pointer Ref | pointing to | POINTING_TO |
| Pointer Deref | pointing from | POINTING_FROM |

This mapping layer is where **human-readable syntax becomes machine-readable logic**.

---

## Token Data Model

The Token struct is the data object passed from the Lexer to the Parser.

    type Token struct {
        Type    TokenType
        Literal string
        Line    int
        Column  int
    }

### Design Goals

- Minimal memory footprint
- Accurate error diagnostics
- IDE and tooling compatibility

Line and column tracking enables:
- Precise runtime errors
- Static analysis
- Editor integration

---

## Performance Considerations

Tokenization is on the compiler hot path.

Optimizations used:
- Keyword map initialized once
- Read-only access during execution
- Immediate fallback to IDENT

There is **no regex**, **no reflection**, and **no allocation** during keyword lookup.

---

## Testing Strategy

The token package is tested at multiple levels.

| Test Suite | Purpose |
|----|----|
| token_unit_test | Keyword correctness |
| token_edge_test | Boundary and malformed input |
| token_integration_test | Feature completeness |
| token_sanity_test | Stability on full programs |
| token_benchmark_test | Performance regression guard |

Tests ensure correctness without coupling to the Lexer or Parser.

---

## Benchmarks

Sample benchmark output:

    BenchmarkLookupIdent
    100000000 iterations
    ~12 ns/op

Interpretation:
- Keyword resolution is effectively free
- Lexer performance is limited by IO, not tokens

---

## Running Tests

Run all token tests:

    go test -v ./token

Run performance benchmarks:

    go test -bench=. ./token