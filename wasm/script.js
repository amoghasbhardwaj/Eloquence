// ==================================================================
// FILE: script.js
// PURPOSE: Logic for Eloquence Playground & Syntax Tabs
// ==================================================================

// --- 1. DATA: EXAMPLES ---
const EXAMPLES = {
    "1. Hello World": `// 1. Basic Output
show("Hello, Eloquence!")

x is 10
y is 20
// Natural language math
result is x adds y

show("Result:", result)`,

    "2. Logic": `// 2. Logic & Control Flow
score is 85

if score greater 80 {
    show("Pass")
} else {
    show("Fail")
}`,

    "3. Loops": `// 3. Loops
n is 3
while n greater 0 {
    show(n)
    n is n minus 1
}
show("Liftoff!")`,

    "4. Functions": `// 4. Functions
add is takes(a, b) {
    return a adds b
}

res is add(10, 20)
show(res)`,

    "5. Structs": `// 5. Structs
define User as struct { name, role }

u is User { 
    name: "Amogh", 
    role: "Admin" 
}
show(u.name)`,

    "6. Pointers": `// 6. Memory
val is 100
ptr is pointing to val

// Mutate via pointer
pointing from ptr is 500
show(val)`
};

// --- 2. DOM ELEMENTS ---
const codeInput = document.getElementById('code');
const outputDiv = document.getElementById('output');
const exampleList = document.getElementById('example-list');
const statusDot = document.querySelector('#status .dot');
const statusText = document.querySelector('#status .text');
const lineNumbers = document.getElementById('lineNumbers');

// --- 3. EDITOR SYNTAX HIGHLIGHTING ---
function update(text) {
    let result_element = document.getElementById("highlighting-content");
    
    // Handle final newline
    if(text[text.length-1] == "\n") text += " ";
    
    // Escape HTML
    result_element.innerHTML = text.replace(/&/g, "&amp;").replace(/</g, "&lt;");

    let data = result_element.innerHTML;

    // 1. Strings
    data = data.replace(/"(.*?)"/g, '<span class="tok-string">"$1"</span>');
    // 2. Comments
    data = data.replace(/(\/\/.*)/g, '<span class="tok-comment">$1</span>');
    // 3. Keywords
    const keywords = ["is", "adds", "subtracts", "times", "divides", "if", "else", "while", "for", "in", "return", "takes", "define", "struct", "as", "true", "false", "none", "pointing", "to", "from", "and", "or", "not", "greater", "less", "equals", "not_equals"];
    keywords.forEach(kw => {
        let regex = new RegExp(`\\b${kw}\\b`, 'g');
        data = data.replace(regex, `<span class="tok-keyword">${kw}</span>`);
    });
    // 4. Functions
    const funcs = ["show", "count", "append", "upper", "lower", "split", "join", "str", "ask"];
    funcs.forEach(fn => {
        data = data.replace(new RegExp(`\\b${fn}\\b`, 'g'), `<span class="tok-func">${fn}</span>`);
    });
    // 5. Numbers
    data = data.replace(/\b(\d+)\b/g, '<span class="tok-number">$1</span>');

    result_element.innerHTML = data;
    updateLineNumbers(text);
}

function syncScroll(element) {
    let result_element = document.querySelector("#highlighting");
    result_element.scrollTop = element.scrollTop;
    result_element.scrollLeft = element.scrollLeft;
    lineNumbers.scrollTop = element.scrollTop;
}

function checkTab(element, event) {
    if(event.key == "Tab") {
        event.preventDefault();
        let start = element.selectionStart;
        let end = element.selectionEnd;
        element.value = element.value.substring(0, start) + "    " + element.value.substring(end);
        element.selectionStart = element.selectionEnd = start + 4;
        update(element.value);
    }
}

// --- 4. APP LOGIC ---

function updateLineNumbers(text) {
    const lines = text.split('\n').length;
    lineNumbers.innerHTML = Array.from({length: lines}, (_, i) => i + 1).join('<br>');
}

// Load Examples
Object.keys(EXAMPLES).forEach((key, index) => {
    const item = document.createElement('div');
    item.innerText = key;
    item.className = "example-item";
    if (index === 0) item.classList.add('active');
    
    item.onclick = () => {
        document.querySelectorAll('.example-item').forEach(i => i.classList.remove('active'));
        item.classList.add('active');
        codeInput.value = EXAMPLES[key];
        update(codeInput.value);
    };
    exampleList.appendChild(item);
});

// Init
codeInput.value = EXAMPLES["1. Hello World"];
update(codeInput.value);

// Run Logic
async function run() {
    outputDiv.innerHTML = "";
    const source = codeInput.value;

    if (!window.runEloquence) {
        log("❌ Engine not loaded.", "log-error");
        return;
    }

    try {
        const response = window.runEloquence(source);
        if (response.error) {
            response.error.forEach(err => log("❌ " + err, "log-error"));
            return;
        }
        if (response.logs) {
            response.logs.split('\n').forEach(line => { if(line) log(line); });
        }
        if (response.result && response.result !== "none") {
            log("➤ " + response.result, "log-result");
        } else if (!response.logs) {
            log("// Done.", "console-placeholder");
        }
    } catch (e) {
        log("System Error: " + e.message, "log-error");
    }
}

function log(msg, cls = "log-entry") {
    const div = document.createElement('div');
    div.className = cls;
    div.innerText = msg;
    outputDiv.appendChild(div);
    outputDiv.scrollTop = outputDiv.scrollHeight;
}

function clearOutput() {
    outputDiv.innerHTML = '<div class="console-placeholder">// Output cleared</div>';
}

// --- 5. DOCS LOGIC ---

function toggleModal(id) {
    const m = document.getElementById(id);
    m.style.display = m.style.display === 'block' ? 'none' : 'block';
}

function showDoc(id) {
    // Hide all pages
    document.querySelectorAll('.doc-page').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.doc-tab').forEach(t => t.classList.remove('active'));
    
    // Show selected
    document.getElementById(id).classList.add('active');
    
    // Highlight tab
    // Find the button that called this function based on onclick text
    const buttons = document.querySelectorAll('.doc-tab');
    for (let btn of buttons) {
        if(btn.getAttribute('onclick').includes(id)) {
            btn.classList.add('active');
            break;
        }
    }
}

window.onclick = (e) => {
    if (e.target.classList.contains('modal')) e.target.style.display = 'none';
}

// --- 6. WASM INIT ---
const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((res) => {
    go.run(res.instance);
    statusDot.classList.add('ready');
    statusText.innerText = "Engine Ready";
    statusDot.style.background = "var(--success)";
});