# 🖋️ Eloquence Programming Language

![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)
![Go Version](https://img.shields.io/badge/go-1.20+-blue.svg?style=flat-square)
![Architecture](https://img.shields.io/badge/architecture-Tree--Walking%20Interpreter-orange.svg?style=flat-square)
![Design](https://img.shields.io/badge/design-English--First-purple.svg?style=flat-square)

> **"Code is for humans to read, and only incidentally for machines to execute."**

**Eloquence** is a high-level, interpreted programming language engineered to eliminate "Symbolic Friction." It replaces dense punctuation of C-style languages (`{`, `&&`, `!=`, `->`) with fluid, natural language phrasing (`begin`, `and`, `not equals`, `pointing to`).  

Designed and Architected by **Amogh S Bharadwaj**.

---

## 📑 Table of Contents

1. [Philosophical Foundation](#-philosophical-foundation)  
2. [Problem & Solution](#-problem--solution)  
3. [Visual Architecture & Pipeline](#-visual-architecture--pipeline)  
4. [Technical Deep Dive: Pratt Parsing](#-technical-deep-dive-pratt-parsing)  
5. [Project Structure & Module Guide](#-project-structure--module-responsibility)  
6. [Runtime Data & Memory Flow](#-runtime-data--memory-flow)  
7. [Comprehensive Syntax Guide](#-comprehensive-syntax-guide)  
8. [Interactive REPL](#-interactive-repl)  
9. [Testing Strategy](#-testing-strategy)  
10. [WebAssembly & Browser Integration](#-webassembly--browser-integration)  
11. [Installation & Usage](#-installation--usage)  
12. [License & Contribution](#-license--contribution)

---

## 📖 Philosophical Foundation

### Etymology

> **el·o·quence** (/ˈeləkwəns/)  
> *noun*  
> 1. Fluent or persuasive speaking or writing.  
> 2. The quality of delivering a clear, strong message.

Eloquence bridges **Logic of the Mind** ↔ **Logic of the Machine**, ensuring code reads like a novel while preserving execution correctness.

---

## ⚠️ Problem & Solution

### Cognitive Load in Traditional Languages

    if (x != null && !y) { doSomething(); }

requires decoding (`!=`, `&&`, `!`) before comprehension.

### Eloquence Semantic Fluency

    if x not_equals none and (a greater b or c less_equal d)
        return func(x)
    end

**Design Principles:**

* Linearity (no semicolons)  
* Explicit scoping (`end`)  
* Semantic operators (`is`, `adds`, `none`)  

---

## 🏗️ Visual Architecture & Pipeline

### Full Compiler Flow

graph TD
    A[Source Code String] --> B[Lexer / Scanner]
    B --> C[Parser / AST Builder]
    C --> D[Evaluator / Interpreter]
    D --> E[Environment / Memory Store]
    D --> F[Output / show()]
    D --> G[Built-in Functions]

**Description:**

* Lexer: Tokenizes semantic phrases (e.g., `pointing to`)  
* Parser: Constructs AST using Pratt Parsing  
* Evaluator: Tree-walking interpreter executing nodes  
* Environment: Handles scope, closures, pointers, shadowing  
* Output: Captured by `show()` or returned to JS in WASM  

---

### Lexer → Token Flow

graph LR
    input["x is 10 adds 5"] --> lexer[Lexer]
    lexer --> token1[IDENT:x]
    lexer --> token2[IS:is]
    lexer --> token3[INT:10]
    lexer --> token4[ADDS:adds]
    lexer --> token5[INT:5]

* Lexer uses **lookahead** to detect multi-word tokens (`pointing to`, `greater_equal`).  

---

### Parser → AST Flow

graph TD
    tokens[Tokens Array] --> parser[Parser]
    parser --> ast1[AssignmentStatement]
    ast1 --> name["Name: x"]
    ast1 --> value[InfixExpression]
    value --> left[5]
    value --> operator[adds]
    value --> right[10]

* Uses **Pratt parsing** for precedence and correct grouping  
* Handles **block vs struct ambiguity** with 3-token lookahead  

---

### Evaluator → Environment Flow

graph TD
    AST[AST Node] --> eval[Evaluator]
    eval --> env[Environment]
    env --> store[Variable Store / Scope]
    eval --> ptr[Pointer Object]
    ptr --> store
    eval --> output[show() / return value]

* **Closures** capture the environment at definition  
* **Pointers** reference variables across scopes  

---

## 🧠 Technical Deep Dive: Pratt Parsing

Expression:  

    result is 5 adds 10 times 2

Correct evaluation uses **binding power**:

| Keyword         | Binding Power | Role          | Equivalent |
|-----------------|---------------|---------------|------------|
| `is`            | 10            | Assignment    | =          |
| `equals`        | 30            | Comparison    | ==         |
| `adds`          | 40            | Summation     | +          |
| `times`         | 50            | Product       | *          |
| `pointing to`   | 60            | Reference     | &          |

---

## 📂 Project Structure & Module Responsibility

eloquence/

    ast/        # AST Node definitions
    evaluator/  # Runtime evaluation
    lexer/      # Lexical analysis
    object/     # Data types & environment
    parser/     # Pratt parser & precedence
    repl/       # Interactive shell
    token/      # Token constants & keywords
    wasm/       # WebAssembly runtime
    tests/      # System tests
    main.go     # CLI entry point
    main.wasm   # WebAssembly binary

---

## 🗃️ Runtime Data & Memory Flow

### Object System Overview

graph TD
    Integer --> Object[Object Interface]
    Boolean --> Object
    String --> Object
    Array --> Object
    Map --> Object
    StructInstance --> Object
    Function --> Object
    Pointer --> Object

* **Object Interface**: Type() and Inspect()  
* **Environment**: Hash map with `outer` pointer for lexical scoping  
* **Shadowing**: Local writes do not overwrite outer variables  
* **Mutation**: Pointers allow cross-scope updates  

---

### Memory & Environment Flow

graph TD
    Global[Global Environment] --> Local1[Function Call Env]
    Local1 --> Local2[Block Env]
    Local2 --> Lookup[Variable Lookup]
    Pointer --> Lookup
    Eval --> Output[show() / return value]

---

## ⌨️ Comprehensive Syntax Guide

### [📘 Full SYNTAX.md](SYNTAX.md)

### Quick Cheat Sheet

| Concept        | Syntax Example |
|----------------|----------------|
| Assignment     | `x is 10` |
| Math           | `x adds 5 times 2` |
| Logic          | `if x greater 10 ... end` |
| Functions      | `f is takes(x) ... end` |
| Structs        | `define User as struct { name }` |
| Pointers       | `ptr is pointing to x` |
| Output         | `show("Hello World")` |
| Arrays         | `list is [1,2,3]` |
| Maps           | `config is { "key": "value" }` |

---

## 💻 The Interactive REPL

Start:

    go run main.go

**Meta-commands:**

| Command   | Function | Description |
|-----------|---------|-------------|
| .help     | Help Menu | Lists commands |
| .clear    | Reset Memory | Clears Environment |
| .exit     | Quit | Exit REPL |
| .debug    | Toggle Verbose | Shows pipeline details |

---

## 🧪 Testing Strategy

| Tier        | File                     | Purpose |
|-------------|--------------------------|---------|
| Unit        | `*_unit_test.go`         | Component-level correctness |
| Integration | `*_integration_test.go`  | Module interactions |
| Sanity      | `*_sanity_test.go`       | Edge case handling |
| System      | `tests/system_test.go`   | Full pipeline / algorithm verification |

Run:

    go test ./... -v
    go test -bench=. ./...

---

## 🌐 WebAssembly & Browser Integration

### WASM Runtime Flow

graph TD
    JS[Browser] --> runEloquence[JS Function]
    runEloquence --> WASM[WebAssembly Binary]
    WASM --> Evaluator
    Evaluator --> Object/Env
    Evaluator --> Output[logs]
    JS --> Console[Display Logs]

Build WASM:

    GOOS=js GOARCH=wasm go build -o main.wasm wasm/wasm_main.go

HTML Integration:

    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
        const output = runEloquence('show("Hello World")');
        console.log(output.logs);
    </script>

---

## 🚀 Installation & Usage

1. **Prerequisites:** Go 1.20+  
2. **Clone & Build:**

        git clone https://github.com/amogh/eloquence.git
        cd eloquence
        go build -o eloquence main.go

3. **REPL Mode:** `./eloquence`  
4. **Run Script:** `./eloquence script.eq`

---

## 📜 License & Contribution

* **License:** MIT ([LICENSE](LICENSE))  
* **Contributing Guidelines:** [CONTRIBUTING.md](CONTRIBUTING.md)  
* **Code of Conduct:** [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

---

<p align="center">
<b>Eloquence</b><br>
<i>The art of fluent logic.</i>
</p>