# Contributing to Eloquence

First off, thank you for considering contributing to **Eloquence**! 

Eloquence is an ambitious project to reshape how we think about code readability. By contributing, you are helping to build a language that bridges the gap between human thought and machine logic.

We welcome contributions of all forms: bug fixes, new features, documentation improvements, and performance optimizations.

---

## 🧭 The Eloquence Philosophy

Before writing code, please understand the core design principle: **"English-First."**

If you are proposing a new feature or syntax, ask yourself:
*   *Does this read like a sentence?*
*   *Am I introducing a symbol where a word would suffice?*

**Example:**
*   ❌ **Avoid:** `array.push(x)` (Too functional/method-heavy)
*   ✅ **Prefer:** `append x to array` (Reads naturally)

---

## 🛠️ Getting Started

### Prerequisites
*   **Go 1.23+** installed on your machine.
*   A basic understanding of how interpreters work (Lexing -> Parsing -> Evaluation) is helpful, but not required.

### Setting up the Environment

1.  **Fork and Clone** the repository:
    ```bash
    git clone https://github.com/YOUR_USERNAME/eloquence.git
    cd eloquence
    ```

2.  **Download Dependencies**:
    ```bash
    go mod download
    ```

3.  **Run the Sanity Check**:
    Ensure the current build passes before you make changes.
    ```bash
    go test ./...
    ```

---

## 🧪 Testing Strategy

Eloquence relies on a 4-tier testing strategy. **All PRs must pass these tests.**

| Tier | Scope | Command |
| :--- | :--- | :--- |
| **Unit** | Individual components (Lexer/Parser logic) | `go test ./...` |
| **System** | Full program execution (Fibonacci, etc.) | `go test -v ./tests` |
| **Benchmarks** | Performance regression checks | `go test -bench=. ./...` |
| **Race** | Concurrency safety | `go test -race ./...` |

### Adding New Tests
*   If you fix a bug, add a test case in the relevant `*_test.go` file to prevent regression.
*   If you add a new language feature, add a comprehensive scenario in `tests/system_test.go`.

---

## 🏗️ Development Workflow

1.  **Create a Branch**:
    ```bash
    git checkout -b feature/my-new-feature
    ```

2.  **Make Changes**:
    *   Follow standard Go coding conventions.
    *   Use `gofmt` to format your code.
    *   Write comments for exported functions.

3.  **Verify Locally**:
    Run the REPL to manually test your changes.
    ```bash
    go run main.go
    ```

4.  **Commit**:
    Write clear, concise commit messages.
    *   *Good:* `feat: add support for modulo operator`
    *   *Bad:* `fixed math`

5.  **Push and PR**:
    Push your branch and open a Pull Request against `main`.

---

## 📂 Project Navigation

To help you find your way around:

*   **`token/`**: Define new keywords here.
*   **`lexer/`**: If adding new syntax rules, update the state machine here.
*   **`ast/`**: Define new tree nodes (Statements/Expressions) here.
*   **`parser/`**: Implement the parsing logic (Pratt Parsing) here.
*   **`evaluator/`**: Implement the runtime logic here.
*   **`object/`**: Define new data types here.

---

## 🐛 Reporting Bugs

If you find a bug, please create an Issue using the following template:

1.  **Description**: What happened?
2.  **Code Snippet**: The specific `.eq` code that caused the crash.
3.  **Expected Behavior**: What should have happened?
4.  **Actual Behavior**: What actually happened (Stack trace, error message)?

---

## 📜 Pull Request Checklist

Before submitting your PR, please ensure:

- [ ] The code compiles (`go build ./...`).
- [ ] All tests pass (`go test ./...`).
- [ ] You have added tests for your new feature.
- [ ] You have run `gofmt` on your code.
- [ ] You have updated the documentation (if syntax changed).

---

We look forward to seeing your code! Let's make programming eloquent.