<!-- ===================================================== -->
<!-- Lexer Package README — Eloquence Programming Language -->
<!-- ===================================================== -->

<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-English--First%20Language-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Package-Lexer-f2994a?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Stage-Lexical%20Analysis-9b51e0?style=for-the-badge" />
</p>

---

# Lexer Package  
## Eloquence Programming Language

The **lexer** (lexical analyzer) is responsible for **reading Eloquence source code**.

It transforms a raw stream of characters into a structured stream of **Tokens**, as defined by the `token` package.  
This is the **first executable stage** of the Eloquence compiler pipeline.

The lexer understands:

- Whitespace and formatting  
- Comments  
- Identifiers and keywords  
- Numbers and strings  
- Multi-word English keywords  

---

## Table of Contents

- Overview
- Compiler Placement
- Scanning Strategy
- Handling Complex Tokens
- Unicode Support
- Visual Architecture
- Testing & Verification

---

## Overview

The Lexer behaves like a **deterministic state machine** that walks through the source code **one character at a time**.

Its primary responsibility is:

    Group raw characters into meaningful units (tokens)

### Example

Source input:

    x is 10

Token stream output:

    ┌──────────────┬─────────┐
    │ TOKEN TYPE   │ LITERAL │
    ├──────────────┼─────────┤
    │ IDENT        │ "x"     │
    │ IS           │ "is"    │
    │ INT          │ "10"    │
    └──────────────┴─────────┘

The lexer does **not** understand grammar or intent.  
It only answers:

    “What kind of thing is this text?”

---

## Compiler Placement

The lexer is the **gateway** to the compiler.

    ┌────────────┐
    │  Source    │
    │   Code     │
    └─────┬──────┘
          │
          ▼
    ┌────────────┐
    │   Lexer    │
    │ (scanner)  │
    └─────┬──────┘
          │ emits Token{}
          ▼
    ┌────────────┐
    │   Parser   │
    │ (syntax)   │
    └────────────┘

Without the lexer:

- The parser would see raw bytes
- Error reporting would be impossible
- English-first syntax could not exist

---

## Scanning Strategy

### Character Consumption Model

The lexer maintains **two pointers**:

    position      → current character
    readPosition  → next character (lookahead)

Visualized:

    source:  x  i  s     1  0
             ↑
          position
                ↑
           readPosition

This design enables:

- Single-character lookahead
- Safe multi-character token detection
- Zero backtracking

---

### Whitespace & Comments

Eloquence is **format-agnostic**.

The lexer automatically skips:

- Spaces
- Tabs
- Newlines
- Carriage returns

Comment handling:

    // single-line comment
    /* multi-line comment */

Comments are discarded and **never reach the parser**.

---

## Handling Complex Tokens

### Multi-Word Keywords

Eloquence introduces **English phrases** as keywords.

Examples:

    pointing to
    pointing from

Lexer strategy:

    1. Read the first identifier ("pointing")
    2. Peek ahead without consuming
    3. If next word is "to" or "from"
    4. Emit a single composite token

Resulting tokens:

    ┌────────────────┬──────────────────┐
    │ Phrase         │ Token Type       │
    ├────────────────┼──────────────────┤
    │ pointing to    │ POINTING_TO      │
    │ pointing from  │ POINTING_FROM    │
    └────────────────┴──────────────────┘

If no phrase match exists, the word is treated as a normal identifier.

---

### Numeric Literals

Numeric types are identified **during lexing**.

    10     → INT
    3.14   → FLOAT

Detection logic:

    - Read digits
    - If '.' is encountered
    - Peek ahead for digits
    - Promote to FLOAT

This removes ambiguity before parsing begins.

---

## Unicode Support

Eloquence is **human-first**, not ASCII-bound.

The lexer uses Go’s unicode/utf8 decoding to:

- Safely process multi-byte characters
- Support non-Latin identifiers
- Allow emojis and symbols inside strings

This guarantees correctness for UTF-8 encoded source files.

---

## Visual Architecture

The lexer operates in a strict execution loop:

    Source Code
        |
        v
    readChar()
        |
        v
    classify current character
        |
        |-- Symbol? --------> emit token
        |
        |-- Letter? --------> readIdentifier()
        |                       → LookupIdent()
        |
        |-- Digit? ---------> readNumber()
        |
        |-- Quote? ---------> readString()
        |
        v
    return Token to Parser

The loop continues until **EOF** is emitted.

---

## Testing & Verification

The lexer is tested in isolation to ensure correctness.

Test layers:

    - Unit Tests
      Validate individual token recognition
      File: lexer_unit_test.go

    - Integration Tests
      Validate mixed real-world programs
      File: lexer_integration_test.go

    - Sanity Tests
      Ensure stability on large or malformed input
      File: lexer_sanity_test.go

Run all tests:

    go test ./lexer -v

Run benchmarks:

    go test ./lexer -bench=.

---

## Summary

The lexer package:

- Converts raw text into structured tokens
- Enables English-first syntax
- Handles multi-word keywords cleanly
- Supports Unicode safely
- Feeds precise data into the parser

It is the **reader** of the Eloquence language —  
the point where *text becomes structure*.