# Parser - Eloquence Programming Language

The `parser` package is the syntactic analysis engine of the Eloquence compiler. It takes the linear stream of tokens produced by the `lexer` and organizes them into a hierarchical structure called the **Abstract Syntax Tree (AST)**.

This implementation uses a **recursive descent parser** combined with a **Pratt Parser** (Top-Down Operator Precedence) for efficient and flexible expression handling. This allows Eloquence to support complex features like operator precedence, infix/prefix notation, and nested structures seamlessly.

---

## Folder Structure

*   **parser.go**: The core logic of the parser, including the Pratt parsing engine, statement parsing functions, and precedence definitions.
*   **parser_unit_test.go**: Granular tests verifying that individual language constructs (assignments, functions, loops, etc.) are parsed correctly into their corresponding AST nodes.
*   **parser_integration_test.go**: Validates the parsing of complete, multi-statement programs that combine various language features (e.g., recursive functions, struct manipulation).
*   **parser_sanity_test.go**: Ensures the parser handles edge cases like empty inputs, comments, and syntax errors gracefully without crashing.
*   **parser_benchmark_test.go**: Measures the performance of the parser on large inputs and deeply nested structures to ensure scalability.

---

## Core Architecture

### The Parser Struct

The `Parser` struct maintains the state of the parsing process:

| Field | Description |
| :--- | :--- |
| `l` | Pointer to the `lexer.Lexer` instance used to fetch tokens. |
| `curToken` | The current token being examined. |
| `peekToken` | The next token in the stream (lookahead). |
| `errors` | A list of syntax errors encountered during parsing. |
| `prefixParseFns` | Map of functions for parsing prefix expressions (e.g., `-5`, `!true`). |
| `infixParseFns` | Map of functions for parsing infix expressions (e.g., `5 adds 5`). |

### Pratt Parsing (Expression Parsing)

Eloquence uses Pratt parsing to handle expressions. This technique associates parsing functions with token types (prefix or infix) and assigns precedence levels to operators.

*   **Prefix Functions**: Handle tokens that appear at the beginning of an expression (e.g., identifiers, literals, unary operators).
*   **Infix Functions**: Handle tokens that appear between two expressions (e.g., `adds`, `equals`, function calls `()`, index access `[]`).
*   **Precedence**: Ensures that operations like multiplication (`times`) bind more tightly than addition (`adds`), respecting mathematical rules.

---

## Language Constructs Parsed

The parser handles the full range of Eloquence syntax:

### Statements
*   **Assignments**: `x is 10`, `pointing from ptr is 5`
*   **Return**: `return result`
*   **Output**: `show "hello"`
*   **Loops**: `for i less 10 ... end`, `while true ... end`
*   **Control Flow**: `try ... catch ... finally ... end`
*   **Struct Definitions**: `define User as struct { name, age }`

### Expressions
*   **Literals**: Integers, Floats, Strings, Booleans, Arrays, Maps, Functions (`takes ... end`).
*   **Binary Operations**: `adds`, `subtracts`, `times`, `divides`, `equals`, etc.
*   **Prefix Operations**: `not`, `-`, `pointing to`, `pointing from`.
*   **Grouping**: Parentheses `(...)` for controlling precedence.
*   **Complex Access**: Array indexing `arr[0]`, Field access `obj.prop`, Function calls `fn(arg)`.
*   **Control Expressions**: `if ... else ... end`.

---

## Visual Flow

    [ Token Stream ]  -->  [ Parser ]  -->  [ Abstract Syntax Tree (AST) ]

    Token Stream:
    IDENT("x"), IS, INT("5"), ADDS, INT("10")

    Parser Processing:
    1. Encounter "x" (IDENT) -> Start AssignmentStatement
    2. Expect "is" -> Consumed
    3. Parse Expression "5 adds 10":
       - Parse "5" (IntegerLiteral)
       - Lookahead "adds" (Infix Operator)
       - Recursively parse "10" (IntegerLiteral)
       - Construct InfixExpression(5, "adds", 10)
    4. Construct AssignmentStatement("x", InfixExpression(...))

    Resulting AST:
    Program
    └── AssignmentStatement
        ├── Name: Identifier("x")
        └── Value: InfixExpression
            ├── Left: IntegerLiteral(5)
            ├── Operator: "adds"
            └── Right: IntegerLiteral(10)

---

## Testing Strategy

The parser is rigorously tested to ensure it correctly interprets valid code and rejects invalid code with helpful errors.

| Test File | Focus | Pass Criteria |
| :--- | :--- | :--- |
| `parser_unit_test.go` | Individual Constructs | Each AST node type is correctly instantiated with proper values. |
| `parser_integration_test.go` | Complex Logic | Recursive functions and struct interactions produce correct deep AST structures. |
| `parser_sanity_test.go` | Error Handling | Parser reports errors for missing tokens or invalid syntax instead of panicking. |
| `parser_benchmark_test.go` | Performance | Parsing remains efficient (linear time) for large inputs. |

---

## How to Run Tests

    # Run all parser tests
    go test -v ./parser

    # Run benchmarks
    go test -bench=. ./parser

---

## Summary

The **parser package** is the bridge between raw text and executable logic. By leveraging Pratt parsing, it robustly handles the English-first syntax of Eloquence, transforming it into a structured AST ready for evaluation. Its architecture supports extensibility, allowing new language features to be added with minimal friction.