// ==============================================================================================
// FILE: repl/repl.go
// ==============================================================================================
// PACKAGE: repl
// PURPOSE: The Read-Eval-Print Loop interface.
//          It connects the user input stream to the compiler pipeline (Lexer->Parser->Evaluator)
//          and manages the persistent session state.
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
	PROMPT = ">> "
	LOGO   = `
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
// It listens to 'in', evaluates code, and writes results to 'out'.
// The 'env' persists across the session to allow variable storage.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment() // Persistent memory for the session
	debugMode := false

	// Print Welcome Header
	fmt.Fprint(out, LOGO)
	printHelp(out)

	for {
		fmt.Fprint(out, Cyan+PROMPT+Reset)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// --- COMMAND HANDLING ---
		if strings.HasPrefix(line, ".") {
			switch line {
			case ".exit":
				fmt.Fprintln(out, Yellow+"Goodbye!"+Reset)
				return
			case ".clear":
				env = object.NewEnvironment() // Reset environment
				fmt.Fprintln(out, Green+"Environment cleared (memory reset)."+Reset)
				continue
			case ".debug":
				debugMode = !debugMode
				status := "DISABLED"
				if debugMode {
					status = "ENABLED"
				}
				fmt.Fprintf(out, Gray+"Debug mode %s\n"+Reset, status)
				continue
			case ".help":
				printHelp(out)
				continue
			default:
				fmt.Fprintf(out, Red+"Unknown command: %s. Type .help for info.\n"+Reset, line)
				continue
			}
		}

		// --- 1. LEXER DEBUG (Optional) ---
		if debugMode {
			printTokens(out, line)
		}

		// --- 2. PARSER ---
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// --- 3. AST DEBUG (Optional) ---
		if debugMode {
			printAST(out, program)
		}

		// --- 4. EVALUATOR ---
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			printEvalResult(out, evaluated)
		}
	}
}

// ----------------------------------------------------------------------------
// HELPER FUNCTIONS
// ----------------------------------------------------------------------------

func printHelp(out io.Writer) {
	fmt.Fprintln(out, Gray+"Commands:")
	fmt.Fprintln(out, "  .exit   Quit the REPL")
	fmt.Fprintln(out, "  .clear  Reset memory")
	fmt.Fprintln(out, "  .debug  Toggle verbose AST/Token output")
	fmt.Fprintln(out, "  .help   Show this message"+Reset)
	fmt.Fprintln(out)
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
	// We check for non-empty string to avoid printing blank lines
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
