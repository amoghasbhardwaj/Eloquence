package main

import (
	"fmt"
	"os"
	"os/user"

	"eloquence/evaluator"
	"eloquence/lexer"
	"eloquence/object"
	"eloquence/parser"
	"eloquence/repl"
)

func main() {
	// 1. Script Mode: go run main.go myfile.eq
	if len(os.Args) > 1 {
		runFile(os.Args[1])
		return
	}

	// 2. REPL Mode: go run main.go
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Welcome to the Eloquence programming language.\n", currentUser.Username)
	fmt.Println("Type your commands below (or 'go run main.go <file>' to execute a script).")

	repl.Start(os.Stdin, os.Stdout)
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	input := string(data)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Println("Parser Errors:")
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		fmt.Println(evaluated.Inspect())
		os.Exit(1)
	}
}
