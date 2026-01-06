# REPL - Eloquence Programming Language

The **REPL** (Read-Eval-Print Loop) is the interactive shell for the Eloquence language. It provides an immediate feedback loop for developers to test snippets, debug logic, and explore the language features without the overhead of file compilation.

This package orchestrates the entire compiler pipeline: Lexer -> Parser -> Evaluator -> Output.

## Table of Contents

1.  Overview
2.  Folder Structure
3.  Core Features
4.  Commands & Control
5.  Architecture
6.  Testing Strategy
7.  Performance Benchmarks
8.  How to Run

---

## Overview

The Eloquence REPL is designed to be developer-friendly. It maintains a persistent session (Environment) so variables defined in one line are available in the next. It also includes built-in debugging tools to inspect the internal state of the compiler.

**Key Capabilities:**
*   **Persistent State**: Variables (`x is 10`) survive between inputs.
*   **Debugging**: View the raw Tokens and AST structure on command.
*   **Graceful Recovery**: Syntax errors print helpful messages without crashing the shell.

---

## Folder Structure

*   **repl.go**: The main entry point. Handles input scanning, command processing, and pipeline execution.
*   **repl_unit_test.go**: Verifies basic REPL functionality (math, variable persistence, commands).
*   **repl_integration_test.go**: Simulates complex user sessions involving functions, structs, and pointers.
*   **repl_sanity_test.go**: Ensures the REPL handles empty lines, garbage input, and parse errors gracefully.
*   **repl_benchmark_test.go**: Measures startup latency and calculation throughput.

---

## Core Features

### 1. Interactive Execution
Type any valid Eloquence code, and the result is printed immediately.
*   `>> 10 adds 5` -> `15`
*   `>> show "Hello"` -> `Hello`

### 2. Output Formatting
Results are color-coded for readability:
*   **Green**: Booleans (`true`), Strings (`"hello"`)
*   **Yellow**: Integers (`42`), Floats (`3.14`)
*   **Red**: Errors
*   **Blue**: Arrays, Maps
*   **Purple**: Functions

---

## Commands & Control

The REPL supports meta-commands prefixed with a dot `.`:

| Command | Description |
| :--- | :--- |
| `.exit` | Terminates the session. |
| `.clear` | Wipes the current memory environment (resets all variables). |
| `.debug` | Toggles verbose mode. Prints the Token stream and AST tree for every input. |
| `.help` | Displays the help menu. |

---

## Architecture

The `Start` function runs an infinite loop that:
1.  **Reads** a line from `stdin`.
2.  **Lexes** the line into tokens.
3.  **Parses** tokens into an AST.
    *   *If Parsing fails*: Prints errors and restarts loop.
4.  **Evaluates** the AST against the persistent `Environment`.
5.  **Prints** the resulting Object.

---

## Testing Strategy

Since the REPL is the primary interface for users, it must be robust.

| Test Suite | Focus | Pass Criteria |
| :--- | :--- | :--- |
| `repl_unit_test` | Basic Commands | `.clear` actually removes variables; math operations return correct strings. |
| `repl_integration_test` | Complex State | Structs defined in line 1 can be instantiated in line 2 and accessed in line 3. |
| `repl_sanity_test` | Robustness | Hitting Enter on an empty line doesn't crash; syntax errors are caught. |
| `repl_benchmark_test` | Latency | The loop overhead is minimal (< 1ms per interaction). |

---

## Performance Benchmarks

    BenchmarkREPL_StartupAndExit-12    500000    2500 ns/op
    BenchmarkREPL_Calculation-12       300000    4000 ns/op

*   **Low Latency**: Designed to feel instant to the user.
*   **Efficient**: Reuses the same environment pointer to minimize allocation.

---

## How to Run

To start the REPL manually:

    go run main.go

To run the test suite:

    go test -v ./repl