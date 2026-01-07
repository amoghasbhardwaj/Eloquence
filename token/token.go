// ==============================================================================================
// FILE: token/token.go
// ==============================================================================================
// PACKAGE: token
// PURPOSE: Defines the vocabulary of the Eloquence programming language.
//          It maps raw source code text to semantic meanings (Tokens).
//          This file acts as the dictionary for the Lexer and Parser.
// ==============================================================================================

package token

// TokenType is a string type alias that represents the category of a token.
// We use strings (instead of integers) for easier debugging and readability
// during the development of the language core.
type TokenType string

// Token represents a single lexical unit scanned from the source code.
// It acts as the "atom" of the language that the Parser will consume.
type Token struct {
	Type    TokenType // The category of the token (e.g., IDENT, KEYWORD, INT)
	Literal string    // The actual text found in the source code (e.g., "myVar", "10")
	Line    int       // The line number where the token was found (for error reporting)
	Column  int       // The column number where the token starts (for precise error pointing)
}

// ----------------------------------------------------------------------------------------------
// TOKEN CONSTANTS
// ----------------------------------------------------------------------------------------------
// These constants define the exhaustive list of valid tokens in Eloquence.
// ----------------------------------------------------------------------------------------------

const (
	// Special Tokens
	// ----------------
	ILLEGAL = "ILLEGAL" // Represents any character or sequence that the Lexer cannot recognize
	EOF     = "EOF"     // End Of File - signals the Parser to stop processing

	// Identifiers & Literals
	// ----------------------
	IDENT  = "IDENT"  // User-defined names (variables, functions, e.g., "calculate_tax")
	INT    = "INT"    // Integer numbers (e.g., 10, 42)
	FLOAT  = "FLOAT"  // Floating point numbers (e.g., 3.14, 0.001)
	STRING = "STRING" // Text strings (e.g., "Hello World")
	CHAR   = "CHAR"   // Single characters (e.g., 'a')
	BOOL   = "BOOL"   // Boolean values (true, false)
	NIL    = "NIL"    // Represents the absence of value (syntax: "none")

	// Operators (Natural Language)
	// ----------------------------
	// Eloquence replaces cryptic symbols with English words to lower the barrier to entry.
	IS            = "IS"            // Assignment operator (replaces '=')
	ADDS          = "ADDS"          // Addition operator (replaces '+')
	SUBTRACTS     = "SUBTRACTS"     // Subtraction operator (replaces '-')
	TIMES         = "TIMES"         // Multiplication operator (replaces '*')
	DIVIDES       = "DIVIDES"       // Division operator (replaces '/')
	MODULO        = "MODULO"        // Modulo operator (replaces '%')
	MINUS         = "MINUS"         // Unary minus or subtraction (context dependent)
	EQUALS        = "EQUALS"        // Equality check (replaces '==')
	NOT_EQUALS    = "NOT_EQUALS"    // Inequality check (replaces '!=')
	GREATER       = "GREATER"       // Greater than (replaces '>')
	LESS          = "LESS"          // Less than (replaces '<')
	GREATER_EQUAL = "GREATER_EQUAL" // Greater than or equal (replaces '>=')
	LESS_EQUAL    = "LESS_EQUAL"    // Less than or equal (replaces '<=')
	AND           = "AND"           // Logical AND (replaces '&&')
	OR            = "OR"            // Logical OR (replaces '||')
	NOT           = "NOT"           // Logical NOT (replaces '!')

	// Delimiters
	// ----------
	// Standard punctuation to structure the code.
	LPAREN   = "(" // Start of function parameters or grouping
	RPAREN   = ")" // End of function parameters or grouping
	LBRACKET = "[" // Start of array definition or index
	RBRACKET = "]" // End of array definition or index
	LBRACE   = "{" // Start of hash map or struct definition
	RBRACE   = "}" // End of hash map or struct definition
	COMMA    = "," // Separator for elements
	COLON    = ":" // Separator for key-value pairs
	DOT      = "." // Accessor for struct fields or methods

	// Keywords (Control Flow & Definitions)
	// -------------------------------------
	IF      = "IF"      // Start of conditional block
	ELSE    = "ELSE"    // Alternative conditional block
	END     = "END"     // Universal block closer (replaces '}')
	RETURN  = "RETURN"  // Returns a value from a function
	FOR     = "FOR"     // Start of loop
	WHILE   = "WHILE"   // Start of while-loop (syntactic sugar)
	REPEAT  = "REPEAT"  // Start of repeat-loop
	TAKES   = "TAKES"   // Function definition keyword (replaces 'func'/'def')
	RETURNS = "RETURNS" // Function return type definition (optional)
	// NOTE: 'SHOW' is purposefully absent. It is handled as a built-in function (IDENT), not a keyword.
	TRY     = "TRY"     // Start of error handling block
	CATCH   = "CATCH"   // Handle errors
	THROW   = "THROW"   // Raise errors
	FINALLY = "FINALLY" // Always execute block
	IN      = "IN"      // Used in range loops (for x IN list)

	// Pointer Keywords
	// ----------------
	// Eloquence uses explicit phrases for pointers to make memory logic readable.
	POINTING_TO   = "POINTING_TO"   // Reference operator (replaces '&')
	POINTING_FROM = "POINTING_FROM" // Dereference operator (replaces '*')

	// Data Structure & Module Keywords
	// --------------------------------
	STRUCT  = "STRUCT"  // Defines a composite data type
	DEFINE  = "DEFINE"  // Starts a definition statement
	AS      = "AS"      // Linking word for definitions
	INCLUDE = "INCLUDE" // Imports code from another file
)

// keywords map connects the English string literals to their internal TokenType.
// This is used by the Lexer to determine if an identifier is actually a reserved keyword.
var keywords = map[string]TokenType{
	// Operators
	"is":            IS,
	"adds":          ADDS,
	"subtracts":     SUBTRACTS,
	"times":         TIMES,
	"divides":       DIVIDES,
	"modulo":        MODULO,
	"minus":         MINUS,
	"equals":        EQUALS,
	"not_equals":    NOT_EQUALS,
	"greater":       GREATER,
	"less":          LESS,
	"greater_equal": GREATER_EQUAL,
	"less_equal":    LESS_EQUAL,
	"and":           AND,
	"or":            OR,
	"not":           NOT,

	// Control Flow
	"if":      IF,
	"else":    ELSE,
	"end":     END,
	"return":  RETURN,
	"for":     FOR,
	"while":   WHILE,
	"repeat":  REPEAT,
	"takes":   TAKES,
	"returns": RETURNS,
	// "show" is deliberately excluded so it parses as a function identifier
	"try":     TRY,
	"catch":   CATCH,
	"throw":   THROW,
	"finally": FINALLY,
	"in":      IN,

	// Complex Keywords (Handled via specific lexer logic usually, but mapped here for consistency)
	"pointing to":   POINTING_TO,
	"pointing from": POINTING_FROM,

	// Literals
	"true":  BOOL,
	"false": BOOL,
	"none":  NIL,

	// Structs & Modules
	"struct":  STRUCT,
	"define":  DEFINE,
	"as":      AS,
	"include": INCLUDE,
}

// LookupIdent checks if a given identifier string is a reserved keyword.
// If it is a keyword (e.g., "if", "adds"), it returns the specific TokenType.
// If it is not found in the map, it returns token.IDENT, signifying a user-defined name.
//
// Parameters:
//
//	ident (string): The word captured by the lexer.
//
// Returns:
//
//	TokenType: The corresponding token type.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
