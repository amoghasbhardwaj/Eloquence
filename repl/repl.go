// ==============================================================================================
// FILE: repl/repl.go
// ==============================================================================================
// PACKAGE: repl
// PURPOSE: The Read-Eval-Print Loop interface.
//          It handles multiline input buffering and connects to the interpreter.
// ==============================================================================================

package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"eloquence/evaluator"
	"eloquence/lexer"
	"eloquence/object"
	"eloquence/parser"
	"eloquence/token"
)

// ----------------------------------------------------------------------------
// UI CONSTANTS & CONFIGURATION
// ----------------------------------------------------------------------------

const (
	PROMPT      = ">> "
	CONT_PROMPT = "... " // Continuation prompt for multiline blocks
	LOGO        = `
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃  _____ _                                           ┃
┃ | ____| | ___   __ _ _   _  ___ _ __   ___ ___     ┃
┃ |  _| | |/ _ \ / _` + "`" + ` | | | |/ _ \ '_ \ / __/ _ \    ┃
┃ | |___| | (_) | (_| | |_| |  __/ | | | (_|  __/    ┃
┃ |_____|_|\___/ \__, |\__,_|\___|_| |_|\___\___|    ┃
┃                   |_|                              ┃
┃                                                    ┃
┃ The Eloquence Language v0.1                        ┃
┃ Built by Amogh S Bharadwaj                         ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
)

// ANSI Color Codes for terminal output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	Bold   = "\033[1m"
)

// ----------------------------------------------------------------------------
// REPL LOGIC
// ----------------------------------------------------------------------------

// Start launches the Read-Eval-Print Loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment() // Persistent memory for the session
	debugMode := false

	// Print Welcome Header
	fmt.Fprint(out, LOGO)
	fmt.Fprintln(out, "Type .help or .helper for syntax guide.")

	// Buffer to store code across multiple lines (for loops/functions)
	var codeBuffer strings.Builder
	braceCount := 0

	// Initial Prompt
	fmt.Fprint(out, Cyan+PROMPT+Reset)

	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// --- COMMAND HANDLING (Only if not inside a code block) ---
		if braceCount == 0 && strings.HasPrefix(trimmedLine, ".") {
			switch trimmedLine {
			case ".exit":
				fmt.Fprintln(out, Yellow+"Goodbye!"+Reset)
				return
			case ".clear":
				env = object.NewEnvironment() // Reset environment
				codeBuffer.Reset()
				fmt.Fprintln(out, Green+"Environment cleared (memory reset)."+Reset)
				fmt.Fprint(out, Cyan+PROMPT+Reset)
				continue
			case ".debug":
				debugMode = !debugMode
				status := "DISABLED"
				if debugMode {
					status = "ENABLED"
				}
				fmt.Fprintf(out, Gray+"Debug mode %s\n"+Reset, status)
				fmt.Fprint(out, Cyan+PROMPT+Reset)
				continue
			case ".help", ".helper":
				printHelp(out)
				fmt.Fprint(out, Cyan+PROMPT+Reset)
				continue
			default:
				fmt.Fprintf(out, Red+"Unknown command: %s. Type .help for info.\n"+Reset, trimmedLine)
				fmt.Fprint(out, Cyan+PROMPT+Reset)
				continue
			}
		}

		// --- MULTILINE DETECTION ---
		// We count open/close braces to know if the user is finished typing a block.
		braceCount += strings.Count(line, "{")
		braceCount -= strings.Count(line, "}")

		// Append current line to buffer
		codeBuffer.WriteString(line + "\n")

		// If braces are unbalanced (e.g., "while x < 10 {"), wait for more input
		if braceCount > 0 {
			fmt.Fprint(out, Gray+CONT_PROMPT+Reset)
			continue
		}

		// --- EXECUTION PHASE ---
		// User is done typing (braces balanced), let's run the code.
		fullCode := codeBuffer.String()
		codeBuffer.Reset() // Clear buffer for next command
		braceCount = 0     // Reset count safety

		// 1. LEXER DEBUG (Optional)
		if debugMode {
			printTokens(out, fullCode)
		}

		// 2. PARSER
		l := lexer.New(fullCode)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			// Reset prompt and continue loop
			fmt.Fprint(out, Cyan+PROMPT+Reset)
			continue
		}

		// 3. AST DEBUG (Optional)
		if debugMode {
			printAST(out, program)
		}

		// 4. EVALUATOR
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			printEvalResult(out, evaluated)
		}

		// Ready for next input
		fmt.Fprint(out, Cyan+PROMPT+Reset)
	}
}

// ----------------------------------------------------------------------------
// HELPER FUNCTIONS
// ----------------------------------------------------------------------------

func printHelp(out io.Writer) {
	fmt.Fprintln(out, "\n"+Bold+"══════════ ELOQUENCE SYNTAX GUIDE ══════════"+Reset)

	fmt.Fprintln(out, Cyan+"\n[ REPL Commands ]"+Reset)
	fmt.Fprintln(out, "  .exit           Quit the REPL")
	fmt.Fprintln(out, "  .clear          Reset variables/memory")
	fmt.Fprintln(out, "  .debug          Toggle detailed AST/Token view")

	fmt.Fprintln(out, Cyan+"\n[ Variables & Math ]"+Reset)
	fmt.Fprintln(out, "  Assignment      "+Green+"x is 10"+Reset)
	fmt.Fprintln(out, "  Arithmetic      "+Green+"x adds 5, y subtracts 2, z times 3"+Reset)
	fmt.Fprintln(out, "  Comparison      "+Green+"x equals 10, y greater 5"+Reset)

	fmt.Fprintln(out, Cyan+"\n[ Control Flow ]"+Reset)
	fmt.Fprintln(out, "  If/Else         "+Green+"if x < 10 { ... } else { ... }"+Reset)
	fmt.Fprintln(out, "  While Loop      "+Green+"while x < 100 { x is x adds 1 }"+Reset)
	fmt.Fprintln(out, "  Range Loop      "+Green+"for item in myList { show(item) }"+Reset)

	fmt.Fprintln(out, Cyan+"\n[ Functions ]"+Reset)
	fmt.Fprintln(out, "  Define          "+Green+"add is takes(a, b) { return a adds b }"+Reset)
	fmt.Fprintln(out, "  Call            "+Green+"result is add(10, 20)"+Reset)

	fmt.Fprintln(out, Cyan+"\n[ Data Structures ]"+Reset)
	fmt.Fprintln(out, "  Arrays          "+Green+"list is [1, 2, 3]"+Reset)
	fmt.Fprintln(out, "  Maps            "+Green+"dict is {\"name\": \"Amogh\", \"age\": 30}"+Reset)
	fmt.Fprintln(out, "  Struct Def      "+Green+"define User as struct { name, age }"+Reset)
	fmt.Fprintln(out, "  Struct Init     "+Green+"u is User { name: \"Amogh\", age: 30 }"+Reset)

	fmt.Fprintln(out, Cyan+"\n[ Memory ]"+Reset)
	fmt.Fprintln(out, "  Reference       "+Green+"ptr is pointing to x"+Reset)
	fmt.Fprintln(out, "  Dereference     "+Green+"val is pointing from ptr"+Reset)

	fmt.Fprintln(out, Cyan+"\n[ Built-in Functions ]"+Reset)
	fmt.Fprintln(out, "  IO              "+Green+"show(x, y), ask(\"Name?\")"+Reset)
	fmt.Fprintln(out, "  Utils           "+Green+"count(arr), append(arr, item), str(10)"+Reset)
	fmt.Fprintln(out, "  Strings         "+Green+"upper(s), lower(s), split(s, \" \"), join(arr, \",\")"+Reset)

	fmt.Fprintln(out, "\n"+Bold+"════════════════════════════════════════════"+Reset)
}

func printTokens(out io.Writer, line string) {
	fmt.Fprintln(out, Gray+"┌── [ TOKENS ] ──────────────────────────────────────────┐"+Reset)
	l := lexer.New(line)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Fprintf(out, "│ %-15s : %s\n", tok.Type, tok.Literal)
	}
	fmt.Fprintln(out, Gray+"└────────────────────────────────────────────────────────┘"+Reset)
}

func printAST(out io.Writer, program fmt.Stringer) {
	fmt.Fprintln(out, Gray+"┌── [ AST TREE ] ────────────────────────────────────────┐"+Reset)
	if str := program.String(); str != "" {
		fmt.Fprintf(out, "%s\n", str)
	}
	fmt.Fprintln(out, Gray+"└────────────────────────────────────────────────────────┘"+Reset)
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, Red+Bold+"Whoops! Parser Errors:"+Reset)
	for _, msg := range errors {
		fmt.Fprintf(out, Red+"  ✖ %s\n"+Reset, msg)
	}
}

// printEvalResult formats the output based on object type
func printEvalResult(out io.Writer, obj object.Object) {
	if obj == nil || obj.Type() == object.NULL_OBJ {
		return
	}

	str := obj.Inspect()

	switch obj := obj.(type) {
	case *object.Error:
		fmt.Fprintf(out, Red+Bold+"ERROR: "+Reset+Red+"%s\n"+Reset, obj.Message)
	case *object.Integer, *object.Float:
		fmt.Fprintf(out, Yellow+"%s\n"+Reset, str)
	case *object.Boolean:
		color := Green
		if !obj.Value {
			color = Red
		}
		fmt.Fprintf(out, color+"%s\n"+Reset, str)
	case *object.String:
		fmt.Fprintf(out, Green+"%s\n"+Reset, str)
	case *object.ReturnValue:
		printEvalResult(out, obj.Value)
	case *object.Function:
		fmt.Fprintf(out, Purple+"(function)\n"+Reset)
	case *object.Array:
		fmt.Fprintf(out, Blue+"%s\n"+Reset, str)
	case *object.Map:
		fmt.Fprintf(out, Blue+"%s\n"+Reset, str)
	case *object.StructInstance:
		fmt.Fprintf(out, Cyan+"%s\n"+Reset, str)
	case *object.Pointer:
		fmt.Fprintf(out, Gray+"(ptr -> %s)\n"+Reset, obj.Name)
	default:
		fmt.Fprintf(out, "%s\n", str)
	}
}
