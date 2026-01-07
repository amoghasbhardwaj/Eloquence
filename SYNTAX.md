<!-- ======================================================== -->
<!-- SYNTAX GUIDE — Eloquence Programming Language -->
<!-- ======================================================== -->

<p align="center">
  <img src="https://img.shields.io/badge/Eloquence-Syntax%20Guide-2f80ed?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Philosophy-English--First-6fcf97?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Version-1.0-111111?style=for-the-badge" />
</p>

# 📘 The Eloquence Syntax Book

This document is the definitive reference manual for the Eloquence programming language. It covers **100% of the language's surface area**, from basic literals to advanced pointer arithmetic and closure mechanics.

---

## 📑 Table of Contents

1.  [Comments](#1-comments)
2.  [Variables & Assignment](#2-variables--assignment)
3.  [Data Types & Literals](#3-data-types--literals)
4.  [Operators (The English Layer)](#4-operators-the-english-layer)
5.  [Control Flow](#5-control-flow)
6.  [Functions & Closures](#6-functions--closures)
7.  [Data Structures](#7-data-structures)
8.  [Object-Oriented Programming (Structs)](#8-object-oriented-programming-structs)
9.  [Memory Management (Pointers)](#9-memory-management-pointers)
10. [Error Handling](#10-error-handling)
11. [Modules System](#11-modules-system)
12. [Standard Library (Built-ins)](#12-standard-library-built-ins)
13. [Operator Precedence](#13-operator-precedence)

---

## 1. Comments

Eloquence supports C-style comments. These are stripped out by the Lexer before compilation.

    // This is a single-line comment. Use it for brief notes.

    /* 
       This is a multi-line block comment.
       It is useful for commenting out large chunks of code
       or writing detailed documentation.
    */

---

## 2. Variables & Assignment

Variables in Eloquence are **dynamically typed**. Types are inferred from assigned values.  
The assignment keyword is **is**.

    // Declaration and Assignment
    age is 25
    name is "Eloquence"
    is_active is true

    // Re-assignment (Dynamic Typing)
    age is "Twenty Five"   // Valid: 'age' is now a String

---

## 3. Data Types & Literals

Type        | Description                       | Examples
----------- | --------------------------------- | -------------------
Integer     | 64-bit signed integers             | 0, 10, -42, 9999
Float       | 64-bit floating point numbers      | 3.14, -0.001, 10.5
String      | UTF-8 sequences in double quotes, supports escapes | "Hello", "Line\nBreak"
Boolean     | Logical truth values               | true, false
Null        | Represents absence of value        | none

---

## 4. Operators (The English Layer)

Eloquence replaces cryptic symbols with **readable English keywords**.

### Arithmetic Operators

Operator     | Keyword     | Syntax Example      | Standard Equivalent
------------ | ---------- | ----------------- | -----------------
Addition     | adds       | x is 10 adds 5     | +
Subtraction  | subtracts / minus | x is 10 minus 2 | -
Multiplication | times    | x is 5 times 5     | *
Division     | divides    | x is 100 divides 2 | /
Modulo       | modulo     | x is 10 modulo 3   | %

### Comparison Operators

Operator         | Keyword       | Syntax Example       | Standard Equivalent
---------------- | ------------- | ------------------ | -----------------
Equals           | equals        | if x equals 10      | ==
Not Equals       | not_equals    | if x not_equals y   | !=
Greater Than     | greater       | if x greater 5      | >
Less Than        | less          | if x less 10        | <
Greater/Eq       | greater_equal | if x greater_equal 1| >=
Less/Eq          | less_equal    | if x less_equal 0   | <=

### Logical Operators

Operator | Keyword        | Syntax Example     | Standard Equivalent
-------- | -------------- | ---------------- | -----------------
AND      | and            | if x and y        | &&
OR       | or             | if x or y         | ||
NOT      | not / !        | if not valid      | !

---

## 5. Control Flow

**Crucial Rule:** All control flow blocks **must be wrapped in curly braces { ... }**.

### Conditional (If / Else)

    score is 85

    if score greater_equal 90 {
        show("Grade: A")
    } else {
        if score greater_equal 80 {
            show("Grade: B")
        } else {
            show("Grade: C")
        }
    }

### While Loop

Executes while a condition is true. Alias: **repeat**.

    counter is 5
    while counter greater 0 {
        show("Countdown:", counter)
        counter is counter subtracts 1
    }

### Range Loop (For In)

Iterates over arrays.

    fruits is ["Apple", "Banana", "Cherry"]

    for fruit in fruits {
        show("Current fruit:", fruit)
    }

---

## 6. Functions & Closures

Functions are defined using **takes**. First-class citizens.

### Definition

    add is takes(x, y) {
        return x adds y
    }

### Invocation

    result is add(10, 20)

### Anonymous Functions & Closures

    make_multiplier is takes(factor) {
        return takes(x) {
            return x times factor
        }
    }

    doubler is make_multiplier(2)
    show(doubler(5))  // 10

---

## 7. Data Structures

### Arrays

Ordered, can hold mixed types.

    list is [10, "Hello", true]
    show(list[0])          // 10
    list is append(list, 42)
    count(list)             // 4

### Hash Maps (Dictionaries)

Key-Value pairs. Keys: String, Integer, Boolean.

    user is {
        "name": "Amogh",
        "role": "Admin",
        1: "Integer Key Supported"
    }

    show(user["name"])     // "Amogh"
    data is { "meta": { "id": 101 } }
    show(data["meta"]["id"]) // 101

---

## 8. Object-Oriented Programming (Structs)

Structs define custom data types with fixed fields.

### Definition

    define Person as struct { firstName, lastName, age }

### Instantiation

    p is Person { 
        firstName: "John", 
        lastName: "Doe", 
        age: 30 
    }

### Field Access

    fullname is p.firstName adds " " adds p.lastName
    show(fullname)

---

## 9. Memory Management (Pointers)

Safe pointer system.

Syntax | Description
------ | -----------
ptr is pointing to x | Reference: link to variable
pointing from ptr   | Dereference (read)
pointing from ptr is val | Dereference (write)

Example: Mutating global state

    x is 100

    mutate is takes() {
        p is pointing to x
        pointing from p is 999
    }

    mutate()
    show(x) // 999

---

## 10. Error Handling

    try {
        result is 10 divides 0
    } catch {
        show("Caught a math error!")
        result is 0
    } finally {
        show("Cleanup complete.")
    }

---

## 11. Modules System

Include external files:

    // math_lib.eq
    // ... define functions ...

    // main.eq
    include "math_lib.eq"
    // functions are available here

---

## 12. Standard Library (Built-ins)

Function | Signature | Description
-------- | -------- | -----------
show     | show(arg1, arg2, ...) | Prints to console/buffer
count    | count(collection)       | Length of array/string
append   | append(array, item)     | Adds item to end
upper    | upper(string)           | Uppercase
lower    | lower(string)           | Lowercase
split    | split(string, sep)      | Splits string into array
join     | join(array, sep)        | Joins array of strings
str      | str(value)              | Converts to string
ask      | ask(prompt)             | Prompt for user input

---

## 13. Operator Precedence

From **highest → lowest**:

* Index / Call: arr[i], func(), struct.field  
* Prefix: pointing to, not, -5  
* Product: times, divides, modulo  
* Sum: adds, subtracts  
* Comparison: less, greater, equals, not_equals...  
* Logic: and, or  
* Assignment: is  

Example:

    x is 5 adds 10 times 2
    // Parsed as: x is (5 adds (10 times 2)) → Result: 25