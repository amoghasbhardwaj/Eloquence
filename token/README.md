# Token Package - Eloquence Programming Language

The token package constitutes the atomic vocabulary of the Eloquence programming language. It serves as the architectural foundation that bridges raw text processing (Lexing) and semantic understanding (Parsing).

This package defines the canonical set of symbols, keywords, and operators that transform Eloquence from a concept into a computable syntax, strictly adhering to the "English-First" design philosophy.

## Table of Contents

1.  Conceptual Overview
2.  Architecture and Implementation
3.  The English-First Mapping
4.  Data Structures
5.  Performance Strategy
6.  Visual Architecture
7.  Testing and Quality Assurance

---

## Conceptual Overview

### What is a Token?
In compiler design, if source code is a paragraph of text, a Token is a single grammatical unit: a word, a punctuation mark, or a number.

The computer does not natively understand the string "x adds 5". The Token Package defines the types that allow the Lexer to transform that string into a stream of meaningful data objects:
1.  IDENT "x"
2.  OP_ADDS
3.  INT "5"

### Role in Eloquence
Unlike C-style languages that rely on dense symbols (like {, &&, !=), Eloquence relies on natural language. The Token package is responsible for standardizing these English phrases into machine-readable constants (e.g., mapping "pointing to" to POINTING_TO).

---

## Architecture and Implementation

### 1. Type Definition Strategy
We utilize a String-Typed Enum approach for TokenType.

    type TokenType string

*   Why String? While integers are marginally faster in some specific compiler implementations, using strings provides superior debuggability. When a parser crashes with an error saying "Expected IDENT, got EOF", it is immediately readable by the developer without needing to look up integer codes in a reference table.

### 2. O(1) Keyword Resolution
The LookupIdent function is the critical logic gate of this package. It uses a Hash Map Lookup strategy rather than a switch statement or linear search.

*   Mechanism: Go's native map implementation.
*   Complexity: O(1) (Constant Time).
*   Purpose: To instantly differentiate between a user's variable (e.g., myVar) and a reserved language keyword (e.g., return) regardless of how many keywords are added to the language in the future.

---

## The English-First Mapping

Eloquence replaces dense symbolic notation with readable English. This table illustrates the translation layer defined in token.go.

| Category       | Standard C/Go Syntax | Eloquence Keyword | Internal Token Type |
| :---           | :---                 | :---              | :---                |
| Assignment     | =                    | is                | IS                  |
| Arithmetic     | +                    | adds              | ADDS                |
| Arithmetic     | *                    | times             | TIMES               |
| Comparison     | ==                   | equals            | EQUALS              |
| Comparison     | >=                   | greater_equal     | GREATER_EQUAL       |
| Pointer Ref    | &                    | pointing to       | POINTING_TO         |
| Pointer Deref  | *                    | pointing from     | POINTING_FROM       |
| Definition     | func                 | takes             | TAKES               |
| Nullability    | nil / null           | none              | NIL                 |

---

## Data Structures

### The Token Struct
This struct is the data transfer object (DTO) passed from the Lexer to the Parser. It is designed to be lightweight to minimize garbage collection overhead during large compilations.

    type Token struct {
        Type    TokenType // The semantic category (e.g., IS, IDENT, INT)
        Literal string    // The actual text captured (e.g., "is", "myVar", "10")
        Line    int       // Line number (Critical for error reporting)
        Column  int       // Column number (Critical for IDE linting)
    }

---

## Performance Strategy

The LookupIdent function is hot-path code, meaning it is executed for every single word found in a source file.

1.  Map Initialization: The keywords map is initialized once at startup time.
2.  No Allocation: The lookup performs a read-only operation on the map; it does not allocate new memory during validation.
3.  Short-Circuiting: If a word is not found in the map, it immediately defaults to being an Identifier (user variable).

---

## Visual Architecture

The following flow illustrates the lifecycle of a word within the Eloquence front-end:

    RAW SOURCE CODE ("x is 10")
          |
          v
      [ LEXER ] 
    (Breaks text into words: "x", "is", "10")
          |
          v
    [ token.LookupIdent ]
          |
    +-----+----------------------+
    |                            |
    v                            v
    RESERVED KEYWORD             USER IDENTIFIER
    Input: "is"                  Input: "x"
    Output: token.IS             Output: token.IDENT
          |                            |
          +----------+-----------------+
                     |
                     v
                [ PARSER ]
          (Builds Abstract Syntax Tree)

---

## Testing and Quality Assurance

The package employs a rigorous testing matrix to ensure stability and performance across all operating conditions.

### 1. Unit Tests (token_unit_test.go)
Validates the basic 1-to-1 mapping of keywords to constants.
*   Goal: Ensure "is" returns IS and "while" returns WHILE.

### 2. Edge Case Tests (token_edge_test.go)
Tests boundary conditions and unusual inputs.
*   Scenarios: Empty strings, numeric-prefixed identifiers (e.g., "123abc"), case sensitivity checks ("TRUE" vs "true"), and partial matches for multi-word keywords.

### 3. Sanity Tests (token_sanity_test.go)
A high-level check using a simulated program stream.
*   Goal: Ensure that a sequence of tokens representing a real program does not trigger panics or state corruption.

### 4. Integration Tests (token_integration_test.go)
Groups tests by functional category (Math, Logic, Pointers, Structures) to ensure complete feature coverage.

### 5. Benchmarks (token_benchmark_test.go)
Measures the latency of the lookup function.
*   Goal: Ensure lookups remain in the nanosecond range to support fast compilation.

### Running the Tests

To execute the full suite with verbose output:

    go test ./token -v

To run the performance benchmarks:

    go test ./token -bench=.