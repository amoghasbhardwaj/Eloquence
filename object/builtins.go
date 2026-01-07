// ==============================================================================================
// FILE: object/builtins.go
// ==============================================================================================
package object

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Builtins is the list of available native functions
var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"show",
		&Builtin{Fn: func(args ...Object) Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			// Print all arguments separated by space
			fmt.Println(strings.Join(parts, " "))
			return &Null{}
		}},
	},
	{
		"count",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newBuiltinError("argument to `count` not supported, got %s", args[0].Type())
			}
		}},
	},
	{
		"append",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newBuiltinError("first argument to `append` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)

			// FIX for S1019: redundant capacity argument removed
			newElements := make([]Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &Array{Elements: newElements}
		}},
	},
	{
		"ask",
		&Builtin{Fn: func(args ...Object) Object {
			// Print prompt if provided
			if len(args) > 0 {
				fmt.Print(args[0].Inspect() + " ")
			}

			// IMPROVEMENT: Use bufio to read the full line (including spaces)
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				return &Null{}
			}

			// Trim the newline character from the input
			text = strings.TrimSpace(text)
			return &String{Value: text}
		}},
	},
	{
		"upper",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 || args[0].Type() != STRING_OBJ {
				return newBuiltinError("upper takes a string")
			}
			return &String{Value: strings.ToUpper(args[0].(*String).Value)}
		}},
	},
	{
		"lower",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 || args[0].Type() != STRING_OBJ {
				return newBuiltinError("lower takes a string")
			}
			return &String{Value: strings.ToLower(args[0].(*String).Value)}
		}},
	},
	{
		"split",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newBuiltinError("split requires (string, separator)")
			}

			parts := strings.Split(str.Value, sep.Value)
			elements := make([]Object, len(parts))
			for i, p := range parts {
				elements[i] = &String{Value: p}
			}
			return &Array{Elements: elements}
		}},
	},
	{
		"join",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*Array)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newBuiltinError("join requires (array, separator)")
			}

			var parts []string
			for _, el := range arr.Elements {
				// Convert every element to string for joining
				if strVal, ok := el.(*String); ok {
					parts = append(parts, strVal.Value)
				} else {
					parts = append(parts, el.Inspect())
				}
			}
			return &String{Value: strings.Join(parts, sep.Value)}
		}},
	},
	{
		"str", // Converts integers/bools/etc to string
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &String{Value: args[0].Inspect()}
		}},
	},
}

// GetBuiltin is a helper to find a function by name
func GetBuiltin(name string) (*Builtin, bool) {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin, true
		}
	}
	return nil, false
}

// Helper function to create errors inside the object package
func newBuiltinError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
