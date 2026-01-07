// ============================================================
// FILE: test_release.eq
// PURPOSE: Comprehensive Verification Suite for Eloquence
// ============================================================

show("════════════════════════════════════════════════════")
show("         ELOQUENCE SYSTEM DIAGNOSTIC v1.0         ")
show("════════════════════════════════════════════════════")

// ------------------------------------------------------------
// SECTION 1: ARITHMETIC & PRECEDENCE
// ------------------------------------------------------------
show("\n[1] ARITHMETIC & PRECEDENCE")

// Test Order of Operations (PEMDAS)
// 10 + (5 * 2) = 20. If precedence is broken, it might equal 30.
calc is 10 adds 5 times 2
show("  10 adds 5 times 2 (Expect 20):", calc)

// Test Negative Numbers and Float Math
neg is -5.5
pos is 10.5
res is pos adds neg
show("  10.5 adds -5.5 (Expect 5):", res)

// ------------------------------------------------------------
// SECTION 2: STRINGS
// ------------------------------------------------------------
show("\n[2] STRING MANIPULATION")

str1 is "Hello"
str2 is "World"
combined is str1 adds " " adds str2
show("  Concatenation:", combined)

txt is "eloquence"
show("  Upper:", upper(txt))
show("  Lower:", lower("TEST"))

csv is "apple,banana,orange"
arr is split(csv, ",")
show("  Split CSV:", arr)
show("  Re-joined:", join(arr, " | "))

// ------------------------------------------------------------
// SECTION 3: LOGIC & CONTROL FLOW
// ------------------------------------------------------------
show("\n[3] LOGIC & CONTROL FLOW")

age is 20
has_ticket is true
is_banned is false

// Test Complex Boolean Logic
if age greater 18 and has_ticket and not is_banned {
    show("  Access Granted (Complex AND/NOT logic works)")
} else {
    show("  Access Denied (Logic Failed)")
}

// Test Else Block
x is 5
if x greater 10 {
    show("  Error: 5 is not greater than 10")
} else {
    show("  Else block executed correctly")
}

// ------------------------------------------------------------
// SECTION 4: LOOPS
// ------------------------------------------------------------
show("\n[4] LOOPS")

// WHILE LOOP
counter is 3
show("  While Loop Countdown:")
while counter greater 0 {
    show("    T-minus:", counter)
    counter is counter subtracts 1
}

// RANGE LOOP (FOR IN)
nums is [10, 20, 30]
sum is 0
for n in nums {
    sum is sum adds n
}
show("  Sum of [10, 20, 30] using Range Loop:", sum)

// ------------------------------------------------------------
// SECTION 5: ARRAYS & MAPS
// ------------------------------------------------------------
show("\n[5] DATA STRUCTURES")

// Array Manipulation
list is [1, 2]
list is append(list, 3)
show("  Appended List:", list)
show("  Count:", count(list))
show("  Index Access [1]:", list[1])

// Map (Hash) Manipulation
hero is {
    "name": "Batman",
    "city": "Gotham",
    "power": "Money"
}
show("  Map Retrieval:", hero["name"], "defends", hero["city"])

// ------------------------------------------------------------
// SECTION 6: FUNCTIONS & RECURSION
// ------------------------------------------------------------
show("\n[6] FUNCTIONS")

// Basic Function
greet is takes(name) {
    return "Welcome " adds name
}
show("  Function Call:", greet("Tester"))

// Recursion (Factorial)
factorial is takes(n) {
    if n less_equal 1 { return 1 }
    return n times factorial(n minus 1)
}
show("  Recursive Factorial(5):", factorial(5))

// Higher Order Function (Passing function as argument)
apply is takes(val, func) {
    return func(val)
}
double is takes(x) { return x times 2 }

res is apply(10, double)
show("  Higher Order Func (Double 10):", res)

// ------------------------------------------------------------
// SECTION 7: STRUCTS
// ------------------------------------------------------------
show("\n[7] STRUCTS")

define User as struct { name, role, active }

// Instantiation
admin is User { name: "Root", role: "SuperAdmin", active: true }

show("  Struct Field Access:", admin.name)
show("  Boolean Field:", admin.active)

// ------------------------------------------------------------
// SECTION 8: POINTERS (MEMORY)
// ------------------------------------------------------------
show("\n[8] MEMORY & POINTERS")

origin is 100
ptr is pointing to origin

show("  Origin:", origin)
show("  Pointer:", ptr)
show("  Dereference:", pointing from ptr)

// Mutation via pointer
pointing from ptr is 500
show("  Origin after pointer write (Expect 500):", origin)

// ------------------------------------------------------------
// SECTION 9: ERROR HANDLING
// ------------------------------------------------------------
show("\n[9] ERROR HANDLING")

try {
    // We simulate an error logic here if throw existed, 
    // or rely on runtime protection.
    // For now, testing the block structure.
    show("  Inside Try Block")
    // Uncommenting line below would trigger catch in a real scenario if x undefined
    // x is unknown_variable 
} catch {
    show("  Inside Catch Block")
} finally {
    show("  Inside Finally Block (Always runs)")
}

show("\n════════════════════════════════════════════════════")
show("SYSTEM CHECK COMPLETE")
show("════════════════════════════════════════════════════")