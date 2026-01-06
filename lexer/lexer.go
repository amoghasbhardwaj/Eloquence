// ----------------------------------------------------------------------------
// FILE: lexer/lexer.go
// ----------------------------------------------------------------------------
package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"eloquence/token"
)

// Lexer represents the state of the source code scanner.
// It iterates through the input string and produces a stream of tokens.
type Lexer struct {
	input        string
	position     int  // Current position in input (points to current char)
	readPosition int  // Current reading position in input (after current char)
	ch           rune // Current char under examination
	line         int  // Line number for error reporting
	column       int  // Column number for error reporting
}

// New initializes a new Lexer with the given input string.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position indices.
// It handles ASCII and UTF-8 characters.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for NUL (signifies EOF)
		l.position = l.readPosition
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size

		if r == '\n' {
			l.line++
			l.column = 0
		} else {
			l.column++
		}
	}
}

// peekChar returns the next character without advancing the lexer's position.
// Useful for lookahead logic (e.g., distinguishing '=' from '==').
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// NextToken inspects the current character and returns the corresponding Token.
// It handles whitespace skipping, comment ignoring, and delegates to specific
// reader methods for identifiers, numbers, and strings.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	// Check for comments (Single line // and Multi line /* */)
	if l.ch == '/' {
		if l.peekChar() == '/' {
			l.skipSingleLineComment()
			return l.NextToken()
		}
		if l.peekChar() == '*' {
			l.readChar()
			l.readChar()
			if !l.skipMultiLineComment() {
				return l.newToken(token.ILLEGAL, "unterminated comment")
			}
			return l.NextToken()
		}
	}

	var tok token.Token

	switch l.ch {
	case '-':
		tok = l.newToken(token.MINUS, string(l.ch))
	case '!':
		tok = l.newToken(token.NOT, string(l.ch))
	case '(':
		tok = l.newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.ch))
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.ch))
	case '{':
		tok = l.newToken(token.LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(token.RBRACE, string(l.ch))
	case ',':
		tok = l.newToken(token.COMMA, string(l.ch))
	case ':':
		tok = l.newToken(token.COLON, string(l.ch))
	case '.':
		// Distinguish between DOT (access) and floats starting with DOT (.5)
		if unicode.IsDigit(l.peekChar()) {
			return l.readNumberToken()
		}
		tok = l.newToken(token.DOT, string(l.ch))
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = l.line
		tok.Column = l.column
	case '\'':
		tok.Type = token.CHAR
		tok.Literal = l.readCharLiteral()
		tok.Line = l.line
		tok.Column = l.column
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.ch) {
			return l.readNumberToken()
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

// newToken creates a Token instance with the given type and literal.
func (l *Lexer) newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

// readIdentifier reads in an identifier and advances the lexer's position
// until it encounters a non-letter-character.
// It also handles multi-word keywords like "pointing to".
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || unicode.IsDigit(l.ch) {
		l.readChar()
	}
	literal := l.input[position:l.position]

	// Handle compound keyword "pointing to/from"
	if literal == "pointing" {
		savedPos := l.position
		savedReadPos := l.readPosition
		savedCh := l.ch
		savedLine := l.line
		savedCol := l.column

		// Look ahead skipping whitespace
		for l.ch == ' ' || l.ch == '\t' {
			l.readChar()
		}

		if isLetter(l.ch) {
			nextStart := l.position
			for isLetter(l.ch) {
				l.readChar()
			}
			nextWord := l.input[nextStart:l.position]

			if nextWord == "to" {
				return "pointing to"
			}
			if nextWord == "from" {
				return "pointing from"
			}
		}

		// Backtrack if not a compound keyword
		l.position = savedPos
		l.readPosition = savedReadPos
		l.ch = savedCh
		l.line = savedLine
		l.column = savedCol
	}

	return literal
}

// readNumberToken reads a number (integer or float) from the input.
func (l *Lexer) readNumberToken() token.Token {
	line := l.line
	column := l.column
	position := l.position
	isFloat := false

	for unicode.IsDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && unicode.IsDigit(l.peekChar()) {
		isFloat = true
		l.readChar()
		for unicode.IsDigit(l.ch) {
			l.readChar()
		}
	}

	literal := l.input[position:l.position]
	if isFloat {
		return token.Token{Type: token.FLOAT, Literal: literal, Line: line, Column: column}
	}
	return token.Token{Type: token.INT, Literal: literal, Line: line, Column: column}
}

// readString reads a string literal enclosed in double quotes.
func (l *Lexer) readString() string {
	var out strings.Builder
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				out.WriteRune('\n')
			case 't':
				out.WriteRune('\t')
			case 'r':
				out.WriteRune('\r')
			case '"':
				out.WriteRune('"')
			case '\\':
				out.WriteRune('\\')
			default:
				out.WriteRune(l.ch)
			}
		} else {
			out.WriteRune(l.ch)
		}
	}
	return out.String()
}

// readCharLiteral reads a single character literal enclosed in single quotes.
func (l *Lexer) readCharLiteral() string {
	l.readChar() // skip opening '
	char := l.ch
	l.readChar() // skip char
	l.readChar() // skip closing '
	return string(char)
}

// skipWhitespace skips over whitespace characters.
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// skipSingleLineComment consumes characters until a newline is found.
func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.skipWhitespace()
}

// skipMultiLineComment consumes characters until "*/" is found.
func (l *Lexer) skipMultiLineComment() bool {
	for {
		if l.ch == 0 {
			return false
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			return true
		}
		l.readChar()
	}
}

// isLetter checks if a rune is a letter or underscore (valid for identifiers).
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}
