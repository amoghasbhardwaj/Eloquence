# AST - Eloquence Programming Language

The `ast` package serves as the **structural backbone** of the Eloquence compiler. It transforms a linear stream of tokens into a hierarchical tree representation. Each node in the tree represents a semantic component of the language—ranging from simple literal values to complex control flow structures.

By abstracting raw source code into this structured format, the AST allows the **Evaluator** and **Compiler** to navigate the logic of an Eloquence program with mathematical precision while retaining its English-first readability.

---

## Folder Structure

To ensure a clean separation of concerns and maintain high testability, the package is organized as follows:

* **ast.go**: The core definition file. It contains the node interfaces and the implementation of all Statement and Expression types.
* **ast_unit_test.go**: Granular tests focusing on the correctness of individual leaf and branch nodes.
* **ast_integration_test.go**: Validates the composition of nodes into full program structures.
* **ast_sanity_test.go**: Stress-tests the AST with deeply nested recursions to ensure stability.
* **ast_benchmark_test.go**: Measures the performance of string serialization and tree traversal.

---

## Core Architecture

The AST is built on a strictly typed interface system that ensures every node can be serialized back to a string for debugging or analyzed for its token content.

### Node Interfaces

| Interface | Purpose | Key Methods |
| :--- | :--- | :--- |
| **Node** | The base interface for every element in the tree. | `TokenLiteral()`, `String()` |
| **Statement** | Represents an action or instruction that does not return a value. | `statementNode()` |
| **Expression** | Represents a unit of code that evaluates to a specific value. | `expressionNode()` |

### Root Node: `Program`
The `Program` node is the entry point of the AST. It contains a slice of `Statements`, representing the entirety of the source code.

---

## Language Constructs

### Statements
Statements define the "skeleton" of the program logic.

| Node Type | Eloquence Syntax Example | Functional Role |
| :--- | :--- | :--- |
| `AssignmentStatement` | `x is 10` | Binds a value to an identifier. |
| `ShowStatement` | `show x` | Standard output instruction. |
| `LoopStatement` | `while x less 10 { ... }` | Iterative control flow. |
| `StructDefinition` | `define Node as struct { ... }` | Custom data type schema. |

### Expressions
Expressions are the "muscles" that perform calculations and return data.

| Node Type | Eloquence Syntax Example | Functional Role |
| :--- | :--- | :--- |
| `InfixExpression` | `x adds y` | Binary operations (Arithmetic/Logic). |
| `PointerReference` | `pointing to x` | Memory address referencing. |
| `FunctionLiteral` | `takes (x) { ... }` | Anonymous function definition. |
| `CallExpression` | `calculate(5, 10)` | Function invocation. |

---

## Visual Flow of the AST

The diagram below illustrates how a simple Eloquence statement is decomposed into an AST structure.

    INPUT: "x is 5 adds 10"
    
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

---

## Testing Matrix

Our testing strategy ensures that as the language grows, the structural integrity of the tree remains intact.

| Test File | Focus | Pass Criteria |
| :--- | :--- | :--- |
| `ast_unit_test.go` | Literals & Operators | Node `.String()` exactly matches the token literal. |
| `ast_integration_test.go` | Nested Blocks & Programs | Recursive structures (If-Else in Loops) serialize correctly. |
| `ast_sanity_test.go` | Deep Recursion | Tree handles 500+ levels of nesting without stack overflow. |
| `ast_benchmark_test.go` | String Throughput | `.String()` execution remains sub-microsecond for common nodes. |

---

## How to Run Tests

Navigate to the package directory and run the following commands:

    # Run all tests (Unit, Integration, Sanity)
    go test -v ./ast

    # Run specifically by category
    go test -v ast_unit_test.go ast.go

    # Run performance benchmarks
    go test -bench=. ./ast

---

## Summary

The **ast package** provides a robust, recursive, and human-readable representation of Eloquence source code. 
* **English-First**: Nodes are named and structured to reflect natural language.
* **Serialized**: Every node supports `.String()` for seamless debugging.
* **Verified**: Comprehensive testing ensures the tree remains stable regardless of program complexity.