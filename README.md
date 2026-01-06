# 🖋️ Eloquence Programming Language

![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)
![Version](https://img.shields.io/badge/version-0.1.0-blue.svg?style=flat-square)
![Architecture](https://img.shields.io/badge/architecture-Tree--Walking%20Interpreter-orange.svg?style=flat-square)
![Design](https://img.shields.io/badge/design-English--First-purple.svg?style=flat-square)

> **"Code is for humans to read, and only incidentally for machines to execute."**

**Eloquence** is a high-level, interpreted programming language engineered to eliminate "Symbolic Friction." It replaces the dense punctuation of C-style languages (`{`, `&&`, `!=`, `->`) with fluid, natural language phrasing (`begin`, `and`, `not equals`, `pointing to`).

Designed and Architected by **Amogh S Bharadwaj**.

---

## 📑 Table of Contents

1.  [Philosophical Foundation](#-philosophical-foundation)
2.  [The Problem & Solution](#-the-problem--solution)
3.  [Visual Architecture](#-visual-architecture--pipeline)
4.  [Technical Deep Dive: Pratt Parsing](#-technical-deep-dive-why-pratt-parsing)
5.  [Project Structure & Module Guide](#-project-structure--module-responsibility)
6.  [Comprehensive Syntax Guide](#-comprehensive-syntax-guide)
7.  [The Interactive REPL](#-the-interactive-repl)
8.  [Testing Strategy](#-testing--verification-matrix)
9.  [Installation & Usage](#-installation--usage)

---

## 📖 Philosophical Foundation

### Etymology
> **el·o·quence** (/ˈeləkwəns/)
> *noun*
> 1. Fluent or persuasive speaking or writing.
> 2. The quality of delivering a clear, strong message.

In software engineering, Eloquence represents the bridge between the **Logic of the Mind** and the **Logic of the Machine**. It asserts that a codebase should read with the cadence and clarity of a novel, ensuring that intent is never obscured by syntax.

---

## ⚠️ The Problem & Solution

### The Friction: Cognitive Load
Modern programming suffers from **Symbolic Density**. When a developer reads `if (x != null && !y)`, their brain performs a two-step process:
1.  **Decode**: Translate symbols (`!=`, `&&`, `!`) into concepts.
2.  **Comprehend**: Understand the logic.

### The Solution: Semantic Fluency
Eloquence removes step #1. It shifts the focus from *decoding* to *reading*.

#### Comparison

    // Traditional (High Cognitive Load)
    if (x != null && (a > b || c <= d)) {
        return func(x);
    }

    // Eloquence (High Readability)
    if x not_equals none and (a greater b or c less_equal d)
        return func(x)
    end

**Key Architectural Decisions:**
*   **Linearity**: No semicolons. Line breaks signify the natural conclusion of a thought.
*   **Explicit Scoping**: No curly braces. Blocks are delimited by keywords (`if` ... `end`).
*   **Semantic Operators**: `is` replaces `=`, `adds` replaces `+`, `none` replaces `null`.

---

## 🏗️ Visual Architecture & Pipeline

Eloquence is engineered in **Go (Golang)**. The interpreter follows a linear, four-stage "Refinement Pipeline."

    [ SOURCE CODE ]  "x is 10 adds 5"
           |
           v
    +-------------+
    |   LEXER     |  <-- The Scanner
    +-------------+      Breaks raw text into atomic units (Tokens).
           |             It handles multi-word operators like "pointing to".
           v
    +-------------+
    |   PARSER    |  <-- The Architect (Pratt Approach)
    +-------------+      Organizes tokens into an Abstract Syntax Tree (AST).
           |             Enforces grammar and operator precedence.
           v
    +-------------+
    |  EVALUATOR  |  <-- The Engine
    +-------------+      Traverses the AST recursively.
           |             Manages Scope, Environment, and Object Creation.
           v
       [ RESULT ]    (Integer Object: 15)

---

## 🧠 Technical Deep-Dive: Why Pratt Parsing?

Eloquence implements a **Top-Down Operator Precedence (Pratt) Parser**.

### The Ambiguity Challenge
In an English-first language, operators are words, not distinct symbols.
Expression: `result is 5 adds 10 times 2`
*   Naive Left-to-Right: `(5 + 10) * 2 = 30` (Incorrect)
*   Mathematical Truth: `5 + (10 * 2) = 25` (Correct)

### The Pratt Solution
Every token in Eloquence is assigned a **Binding Power** (Precedence). The parser uses these values to decide which operands "stick" to which operator.

| Keyword | Binding Power | Role | Equivalent |
| :--- | :--- | :--- | :--- |
| `is` | 10 | Assignment | `=` |
| `equals` | 30 | Comparison | `==` |
| `adds` | 40 | Summation | `+` |
| `times` | 50 | Product | `*` |
| `pointing to` | 60 | Prefix | `&` |

Because `times` (50) has a higher binding power than `adds` (40), the parser groups `10 times 2` together first, ensuring mathematical correctness without complex grammar files.

---

## 📂 Project Structure & Module Responsibility

The codebase is modular, separating concerns into distinct packages for maintainability.

### 💎 `token/` (The Vocabulary)
Defines the language's dictionary. It maps string literals to byte constants.
*   **token.go**: Defines `TokenType` (e.g., `IDENT`, `INT`) and the `LookupIdent` function, which performs O(1) checks to see if a word is a variable or a keyword.

### 🔍 `lexer/` (The Scanner)
Performs lexical analysis.
*   **lexer.go**: Contains the state machine. It reads characters, skips whitespace/comments, and groups characters into Tokens. It includes specific look-ahead logic to detect compound keywords like `pointing to`.

### 🌳 `ast/` (The Blueprint)
Defines the hierarchical structure of the code.
*   **ast.go**: Defines the `Node`, `Statement`, and `Expression` interfaces. It includes struct definitions for every language construct (e.g., `IfExpression`, `FunctionLiteral`, `StructDefinition`).

### ⚙️ `parser/` (The Grammar)
The core logic engine for syntax analysis.
*   **parser.go**: Implements the Pratt Parser. It registers "Prefix" and "Infix" parsing functions for every token type and recursively builds the AST.

### 📦 `object/` (The Runtime)
Defines the memory model and type system.
*   **object.go**: Defines the `Object` interface. Every value (Integer, Boolean, Struct) implements `Type()` and `Inspect()`.
*   **environment.go**: Implements the **Symbol Table**. It handles variable storage, lexical scoping (nested environments), and variable shadowing.

### 🚀 `evaluator/` (The Executioner)
The runtime interpreter.
*   **evaluator.go**: Contains the giant `switch` statement that dispatches logic based on AST node types. It performs arithmetic, executes loops, and manages control flow (return/error).

### 🖥️ `repl/` (The Interface)
The interactive shell.
*   **repl.go**: Orchestrates the loop: Read Input -> Lex -> Parse -> Evaluate -> Print.

---

## ⌨️ Comprehensive Syntax Guide

### 1. Variables & Primitives
Eloquence uses dynamic typing with semantic keywords.

    // Assignment
    age is 25
    pi is 3.14159
    username is "Amogh"
    
    // Booleans
    is_active is true
    
    // Nullability ('none' replaces 'nil')
    data is none

### 2. Mathematics & Logic
    sum is 10 adds 5
    diff is 10 subtracts 2
    prod is 10 times (2 adds 3)
    
    if age greater 18 and username not_equals "Guest"
        show "Welcome, " adds username
    end

### 3. Collections
    // Arrays
    scores is [95, 88, 100]
    first is scores[0]
    
    // HashMaps (Key-Value)
    config is {
        "version": "1.0",
        "debug": false
    }

### 4. Functions & Closures
Functions are first-class citizens.

    // Definition
    add is takes(x, y)
        return x adds y
    end
    
    // Call
    result is add(10, 20)

### 5. Data Structures (Structs)
Object-oriented data modeling.

    define User as struct { name, email }
    
    // Instantiation
    me is User { name: "Amogh", email: "dev@eloquence.io" }
    
    // Field Access
    show me.name

### 6. Memory Management (Pointers)
Low-level control with high-level syntax.

    target is 500
    
    // Create Reference
    ref is pointing to target
    
    // Dereference and Mutation
    pointing from ref is 1000
    
    show target // Output: 1000

---

## 💻 The Interactive REPL

The **Read-Eval-Print Loop** allows instant experimentation. It includes a persistent environment and debugging tools.

### How to Run
    go run main.go

### Meta-Commands
The REPL supports special commands prefixed with `.`:

| Command | Function | Description |
| :--- | :--- | :--- |
| **.help** | Help Menu | Displays the list of available commands. |
| **.clear** | Reset Memory | Wipes the current Environment (deletes all variables) without exiting. |
| **.exit** | Quit | Terminates the session. |
| **.debug** | Toggle Verbose | **Advanced**: Switches on/off the "Compiler Internals" view. |

### Debug Mode Example
When `.debug` is enabled, the REPL reveals the pipeline:

    >> .debug
    Debug mode ENABLED
    
    >> x is 10
    
    ┌── [ TOKENS ] ──────────────────────┐
    │ IDENT           : x                │
    │ IS              : is               │
    │ INT             : 10               │
    └────────────────────────────────────┘
    ┌── [ AST TREE ] ────────────────────┐
    x is 10
    └────────────────────────────────────┘
    10

---

## 🧪 Testing & Verification Matrix

Eloquence is verified by a rigorous 4-tier testing strategy.

| Tier | File | Objective |
| :--- | :--- | :--- |
| **Unit** | `*_unit_test.go` | **Granularity**: Validates individual components (e.g., ensuring `5 adds 5` evaluates to `10`). |
| **Integration**| `*_integration_test.go` | **Interaction**: Ensures the Parser correctly builds ASTs for complex nested structures like closures. |
| **Sanity** | `*_sanity_test.go` | **Resilience**: Ensures the compiler handles empty inputs, comments, and garbage without panicking. |
| **System** | `tests/system_test.go` | **Completeness**: Verifies Turing Completeness by solving algorithms (Fibonacci, Linked Lists). |

---

## 🚀 Installation & Usage

### Prerequisites
*   Go 1.20 or higher.

### 1. Installation
Clone the repository:

    git clone https://github.com/amogh/eloquence.git
    cd eloquence

### 2. Running the Tests
Verify the integrity of the compiler:

    go test ./... -v

### 3. Running a Script
Create a file named `hello.eq`:

    show "Hello, Eloquence!"

Run it:

    go run main.go hello.eq

---

*Eloquence: The art of fluent logic.*