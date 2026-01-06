# Lexer - Eloquence Programming Language

The **lexer** package constitutes the initial phase of the Eloquence compiler pipeline. Its primary responsibility is **lexical analysis**: the transformation of raw source code strings into a structured stream of meaningful tokens.

As an **English-first programming language**, the Eloquence lexer is specifically engineered to treat natural language phrases like `"is"`, `"adds"`, and `"pointing to"` as functional operators, effectively bridging the gap between human readability and machine logic.

---

## Table of Contents

1. Overview
2. Folder Structure
3. Core Architecture
4. Primary Functions
5. Lexical Logic Flow
6. Testing Strategy
7. Benchmarking & Performance
8. How to Run Tests

---

## Overview

The Lexer scans source code character by character using a state-machine approach. Unlike traditional lexers that rely heavily on dense symbols (e.g., `+`, `-`), the Eloquence Lexer prioritizes word-based identification to support its unique, readable syntax.

**Key Technical Features:**
*   **Multi-word Tokenization**: Intelligently combines discrete words like `"pointing"` and `"to"` into a single semantic `POINTER` token.
*   **UTF-8 Support**: Utilizes Go's `rune` type to ensure full compatibility with global character sets.
*   **Lookahead Logic**: Implements a `peekChar` mechanism to resolve ambiguities (like distinguishing `/` from `//`) without advancing the cursor prematurely.

---

## Folder Structure

The package is organized for high modularity and strict separation of concerns:

*   **lexer.go**: The core engine containing the `Lexer` struct and state-management logic.
*   **lexer_unit_test.go**: Exhaustive unit tests for every token type and operator.
*   **lexer_sanity_test.go**: Quick smoke tests for end-to-end code snippet validation.
*   **lexer_integration_test.go**: Validates interaction between the lexer and the token lookup tables.
*   **lexer_benchmark_test.go**: Performance profiling for the tokenization process.

---

## Core Architecture

### The Lexer Struct

The `Lexer` maintains the state of the scanner as it moves through the input string.

| Field | Description |
| :--- | :--- |
| `input string` | The raw source code to be tokenized. |
| `position int` | Current index in `input` (points to `ch`). |
| `readPosition int` | Current read index (points to the next character). |
| `ch rune` | The current character being examined. |
| `line int` | Current line number (for error reporting). |
| `column int` | Current column number (for error reporting). |

---

## Primary Functions

| Function | Description |
| :--- | :--- |
| `New(input string)` | Factory function that initializes the Lexer state. |
| `readChar()` | Advances the cursor and handles EOF (End Of File) logic. |
| `peekChar()` | Non-destructive lookahead used for multi-character symbols and negative numbers. |
| `NextToken()` | The "heart" of the package; evaluates the current state to return the next `token.Token`. |
| `readIdentifier()` | Collects alphanumeric sequences and performs keyword mapping (including multi-word keywords). |
| `readNumberToken()` | Differentiates between `INT` and `FLOAT` types based on decimal points. |
| `readString()` | Captures literal text within double quotes, handling escape sequences. |

---

## Lexical Logic Flow

    RAW SOURCE CODE
          |
          v
    [ Skip Whitespace ]
          |
          v
    [ Check Comments (// or /*) ]
          |
          v
    [ Character Analysis ]
    |
    +-- Letters? ----> [ readIdentifier ] ----> Keyword Lookup / Multi-word Check
    |
    +-- Digits?  ----> [ readNumberToken ] ---> INT / FLOAT
    |
    +-- Quotes?  ----> [ readString ] --------> STRING
    |
    +-- Symbols? ----> [ newToken ] ----------> Symbol Mapping
    |
    +-- 0 / Null ----> [ EOF ]
          |
          v
    [ Emit Token to Parser ]

---

## Testing Strategy

To ensure industrial-grade reliability, the lexer undergoes a multi-tiered testing process.

### Test Matrix

| Test Category | File | Description |
| :--- | :--- | :--- |
| **Unit Tests** | `lexer_unit_test.go` | Table-driven tests verifying every operator (Arithmetic, Comparison, Logical). |
| **Sanity Tests** | `lexer_sanity_test.go` | Validates that the lexer can process a full sentence without panic. |
| **Integration** | `lexer_integration_test.go` | Validates the Lexer's output against the expected Token definitions. |
| **Benchmarks** | `lexer_benchmark_test.go` | Measures the speed of tokenization for performance regression testing. |

---

## Benchmarking & Performance

Performance is a critical metric for compiler tools. We measure the overhead of tokenization to ensure the lexer remains efficient under heavy workloads.

    BenchmarkLexerNextToken-12    2000000    650 ns/op

*   **Low Memory Footprint**: Minimal allocations per token.
*   **Constant Time Lookup**: Keyword resolution uses a hash-map with O(1) complexity.

---

## How to Run Tests

Maintain the integrity of the codebase by running the suite before every commit:

    # Run all tests with verbose output
    go test ./lexer -v

    # Run benchmarks to check for performance regressions
    go test -bench=. ./lexer