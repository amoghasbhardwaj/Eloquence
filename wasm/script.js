// wasm/script.js

// ---------------------------------------------------------
// EXAMPLES LIBRARY
// ---------------------------------------------------------
const EXAMPLES = {
    "1. Hello World": `
// Basic output
show("Hello, Web World!")
x is 10
y is 20
show("Calculation:", x adds y)
`,
    "2. String Manipulation": `
// Built-in string functions
name is "Eloquence"
show("Original:", name)
show("Upper:", upper(name))
show("Lower:", lower(name))

// Splitting and Joining
sentence is "Code is Art"
parts is split(sentence, " ")
show("Split Array:", parts)
show("Joined:", join(parts, "-"))
`,
    "3. Control Flow (Loops)": `
// While Loop
counter is 0
show("Starting Loop...")

// Note: Use 'less' instead of '<'
while counter less 3 {
    show("Counter is:", counter)
    counter is counter adds 1
}

// Range Loop
fruits is ["Apple", "Banana", "Cherry"]
show("My Fruits:")
for fruit in fruits {
    show("- " adds fruit)
}
`,
    "4. Arrays & Maps": `
// Array Operations
nums is [10, 20, 30]
show("Count:", count(nums))

// Append returns a new array
nums is append(nums, 40)
show("Updated:", nums)

// Hash Maps
user is {
    "name": "Amogh",
    "role": "Creator",
    "level": 99
}
show("User Info:", user)
show("User Name:", user["name"])
`,
    "5. Structs (Custom Types)": `
// Define a structure
define Person as struct { name, age, city }

// Instantiate
p1 is Person { name: "Alice", age: 25, city: "NYC" }
p2 is Person { name: "Bob", age: 30, city: "London" }

show(p1)
show(p2.name adds " lives in " adds p2.city)
`,
    "6. Pointers": `
// Memory Management
x is 100
ptr is pointing to x

show("X is:", x)
show("Pointer address:", ptr)
show("Dereferenced:", pointing from ptr)

// Modify via pointer
pointing from ptr is 500
show("X is now:", x)
`,
    "7. Linked List (Advanced)": `
// Implementing a Linked List
define Node as struct { value, next }

// Create Nodes
n3 is Node { value: 30, next: none }
n2 is Node { value: 20, next: n3 }
head is Node { value: 10, next: n2 }

show("Traversing Linked List:")

// Function to traverse
traverse is takes(current) {
    while current not_equals none {
        show("Node Value:", current.value)
        current is current.next
    }
}

traverse(head)
`
};

// ---------------------------------------------------------
// WASM LOADER
// ---------------------------------------------------------
const go = new Go();
let wasmReady = false;

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    wasmReady = true;
    document.getElementById("status").innerText = "Ready";
    document.getElementById("status").style.color = "#4ec9b0";
});

// ---------------------------------------------------------
// UI LOGIC
// ---------------------------------------------------------

function loadExample(key) {
    document.getElementById("code").value = EXAMPLES[key].trim();
}

function run() {
    if (!wasmReady) {
        alert("WASM is still loading...");
        return;
    }

    const input = document.getElementById("code").value;
    const outputDiv = document.getElementById("output");
    
    // Clear previous output
    outputDiv.innerHTML = "";

    try {
        // CALL GO FUNCTION
        const response = runEloquence(input);

        // Print Logs (Output from show)
        if (response.logs) {
            const lines = response.logs.split("\n");
            lines.forEach(line => {
                if(line) outputDiv.innerHTML += `<div class="log-line">${line}</div>`;
            });
        }

        // Print Result
        if (response.result) {
            outputDiv.innerHTML += `<div class="result-line">=> ${response.result}</div>`;
        }

        // Print Errors
        if (response.error) {
            response.error.forEach(err => {
                outputDiv.innerHTML += `<div class="error-line">${err}</div>`;
            });
        }
    } catch (e) {
        outputDiv.innerHTML += `<div class="error-line">System Error: ${e}</div>`;
    }
}

// Populate Sidebar
window.onload = function() {
    const sidebar = document.getElementById("example-list");
    for (const key in EXAMPLES) {
        const btn = document.createElement("button");
        btn.innerText = key;
        btn.className = "example-btn";
        btn.onclick = () => loadExample(key);
        sidebar.appendChild(btn);
    }
    // Load first example by default
    loadExample("1. Hello World");
};