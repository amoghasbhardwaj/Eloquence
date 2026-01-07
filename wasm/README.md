<!-- ======================================================== -->
<!-- WASM Package README — Eloquence Programming Language -->
<!-- ======================================================== -->

<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-English--First%20Language-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Package-WASM-f2c94c?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Stage-Web%20Integration-111111?style=for-the-badge" />
</p>

---

# WebAssembly Interface  
## Eloquence Programming Language

The **wasm** package provides a bridge between the Eloquence Go compiler and the JavaScript browser environment.  

It allows Eloquence code to be compiled, run, and interacted with entirely **client-side**, without requiring a backend server. This is the engine behind the **online Playground**.

---

## Table of Contents

1. [Overview](#1-overview)  
2. [Architecture](#2-architecture)  
3. [JavaScript Bridge (`runEloquence`)](#3-javascript-bridge-rune eloquence)  
4. [Output Buffering](#4-output-buffering)  
5. [Built-in Overrides](#5-built-in-overrides)  
6. [Building the WASM Binary](#6-building-the-wasm-binary)  
7. [Integration with HTML/JS](#7-integration-with-htmljs)  

---

## 1. Overview

Go supports compilation to **WebAssembly (.wasm)**, a binary format for a stack-based virtual machine that runs in modern browsers.  

The wasm package wraps the core Eloquence interpreter so that:

- **JavaScript** can pass source code strings into Go.  
- **Results & logs** are returned as JSON objects.

---

## 2. Architecture

The entry point is `wasm_main.go`. Unlike `main.go` which interacts with terminal I/O:

- **Channel Initialization**: Keeps Go runtime alive to listen for JS calls indefinitely.  
- **Exposed Functions**: Attaches `runCode` to the global JS object as `runEloquence`.  
- **I/O Overrides**: Redirects `show()` to a string buffer instead of printing to terminal.  

---

## 3. JavaScript Bridge (`runEloquence`)

The exposed function signature in JavaScript:

```ts
runEloquence(sourceCode: string) -> { logs: string, result: string, error?: string[] }
```

**Input:**

- `sourceCode`: Raw Eloquence code entered by the user.

**Output (JSON Object):**

| Property | Description |
|----------|------------|
| `logs`   | Accumulated output from `show()` calls |
| `result` | Final return value of the script (if any) |
| `error`  | Array of strings for parser/runtime errors (optional) |

---

## 4. Output Buffering

Browsers lack a standard stdout, so `strings.Builder` is used as `outputBuffer`.

- **Before Execution**: `outputBuffer.Reset()`  
- **During Execution**: `show()` appends text to the buffer  
- **After Execution**: Contents returned to JS as `logs`

---

## 5. Built-in Overrides

Some standard library functions are modified for browser execution:

| Function | Standard Behavior | WASM Behavior |
|----------|-----------------|---------------|
| `show(...)` | Prints to `os.Stdout` | Appends to `outputBuffer` |
| `ask(...)`  | Reads from `os.Stdin` | Returns a placeholder string (blocking input not supported) |

---

## 6. Building the WASM Binary

Set Go target environment variables and build:

```bash
GOOS=js GOARCH=wasm go build -o main.wasm wasm/wasm_main.go
```

This produces `main.wasm` (~2–3MB) which can be served **statically**.

---

## 7. Integration with HTML/JS

**Step 1: Copy `wasm_exec.js`**

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

**Step 2: Load in HTML**

```html
<script src="wasm_exec.js"></script>
<script>
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
        console.log("Eloquence Engine Loaded");
    });
</script>
```

**Step 3: Execute Eloquence Code**

```js
const output = runEloquence('show("Hello World")');
console.log(output.logs); // "Hello World\n"
```

---

### Summary

The **WASM package** enables:

- Running Eloquence entirely in the browser  
- Real-time code execution for Playground and demos  
- Seamless I/O handling via buffers  
- Lightweight integration with HTML/JS  