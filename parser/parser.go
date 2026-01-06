// ==============================================================================================
// FILE: parser/parser.go
// ==============================================================================================
// PACKAGE: parser
// PURPOSE: Implements a Recursive Descent Parser with Pratt Parsing for expressions.
//          It converts a stream of Tokens (from the Lexer) into an Abstract Syntax Tree (AST).
//          This component defines the grammar and syntax rules of Eloquence.
// ==============================================================================================

package parser

import (
	"fmt"
	"strconv"

	"eloquence/ast"
	"eloquence/lexer"
	"eloquence/token"
)

// Precedence constants determine the order of operations in expressions.
// Higher values mean the operator binds more tightly.
const (
	_ int = iota
	LOWEST
	EQUALS      // equals, not_equals
	LESSGREATER // less, greater
	SUM         // adds, subtracts
	PRODUCT     // times, divides
	PREFIX      // -x, !x
	CALL        // myFunction(x)
	INDEX       // myArray[i], myStruct.field
)

// precedences maps token types to their integer precedence level.
var precedences = map[token.TokenType]int{
	token.EQUALS:        EQUALS,
	token.NOT_EQUALS:    EQUALS,
	token.LESS:          LESSGREATER,
	token.GREATER:       LESSGREATER,
	token.LESS_EQUAL:    LESSGREATER,
	token.GREATER_EQUAL: LESSGREATER,
	token.ADDS:          SUM,
	token.SUBTRACTS:     SUM,
	token.MINUS:         SUM, // 'minus' can be infix
	token.TIMES:         PRODUCT,
	token.DIVIDES:       PRODUCT,
	token.MODULO:        PRODUCT,
	token.LPAREN:        CALL,
	token.LBRACKET:      INDEX,
	token.DOT:           CALL,
	token.LBRACE:        CALL, // For struct instantiation
	token.AND:           EQUALS,
	token.OR:            EQUALS,
}

// Function types for Pratt Parsing
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser struct holds the state of the parsing process.
type Parser struct {
	l         *lexer.Lexer // Pointer to the lexer
	curToken  token.Token  // The current token under examination
	peekToken token.Token  // The next token (lookahead)
	errors    []string     // Collection of syntax errors found

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New initializes a new Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Register Prefix Parsing Functions (nuds)
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.CHAR, p.parseCharLiteral)
	p.registerPrefix(token.BOOL, p.parseBooleanLiteral)
	p.registerPrefix(token.NIL, p.parseNilLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.POINTING_TO, p.parsePointerReferenceExpression)
	p.registerPrefix(token.POINTING_FROM, p.parsePointerDereferenceExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseMapLiteral)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.TAKES, p.parseFunctionLiteral)

	// Register Infix Parsing Functions (leds)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.ADDS, p.parseInfixExpression)
	p.registerInfix(token.SUBTRACTS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.TIMES, p.parseInfixExpression)
	p.registerInfix(token.DIVIDES, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.EQUALS, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUALS, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)
	p.registerInfix(token.LESS_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.GREATER_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseFieldAccessExpression)
	p.registerInfix(token.LBRACE, p.parseStructInstantiation)

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }

// expectPeek asserts that the next token is of a specific type.
// If it is, it advances the parser. If not, it records an error.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("line %d:%d - expected next token to be %s, got %s instead",
		p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram is the entry point for parsing. It iterates through tokens
// and constructs the root AST node (Program).
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement determines the type of statement based on the current token.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		// Lookahead to distinguish assignments ("x is 5") from expressions ("x")
		if p.peekTokenIs(token.IS) {
			return p.parseAssignmentStatement()
		}
		return p.parseExpressionStatement()
	case token.POINTING_FROM:
		// "pointing from ptr is 5" -> Pointer Assignment
		if p.peekTokenIs(token.IDENT) {
			return p.parsePointerAssignmentStatement()
		}
		return p.parseExpressionStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.SHOW:
		return p.parseShowStatement()
	case token.TRY:
		return p.parseTryCatchStatement()
	case token.FOR, token.WHILE:
		return p.parseLoopStatement()
	case token.DEFINE:
		return p.parseStructDefinitionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseBlockStatement parses a block of code terminated by keywords like END, ELSE, etc.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	for !p.curTokenIs(token.END) && !p.curTokenIs(token.EOF) &&
		!p.curTokenIs(token.ELSE) && !p.curTokenIs(token.CATCH) && !p.curTokenIs(token.FINALLY) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // Skip IDENT
	p.nextToken() // Skip IS
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parsePointerAssignmentStatement() *ast.PointerAssignmentStatement {
	stmt := &ast.PointerAssignmentStatement{Token: p.curToken}
	// pointing from [IDENT] is [VALUE]
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.IS) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseShowStatement() *ast.ShowStatement {
	stmt := &ast.ShowStatement{Token: p.curToken}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseLoopStatement() ast.Statement {
	stmt := &ast.LoopStatement{Token: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	p.nextToken()
	stmt.Body = p.parseBlockStatement()

	if !p.curTokenIs(token.END) {
		p.peekError(token.END)
		return nil
	}
	return stmt
}

func (p *Parser) parseTryCatchStatement() ast.Statement {
	stmt := &ast.TryCatchStatement{Token: p.curToken}
	p.nextToken()
	stmt.TryBlock = p.parseBlockStatement()

	if p.curTokenIs(token.CATCH) {
		p.nextToken()
		stmt.CatchBlock = p.parseBlockStatement()
	}
	if p.curTokenIs(token.FINALLY) {
		p.nextToken()
		stmt.FinallyBlock = p.parseBlockStatement()
	}
	if !p.curTokenIs(token.END) {
		p.peekError(token.END)
		return nil
	}
	return stmt
}

func (p *Parser) parseStructDefinitionStatement() *ast.StructDefinitionStatement {
	stmt := &ast.StructDefinitionStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.AS) {
		return nil
	}
	if !p.expectPeek(token.STRUCT) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Attributes = []*ast.Identifier{}
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return stmt
	}

	p.nextToken()
	stmt.Attributes = append(stmt.Attributes, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		stmt.Attributes = append(stmt.Attributes, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return stmt
}

// parseExpression manages precedence to parse expressions correctly.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// --- Prefix Parsing Functions ---

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as float", p.curToken.Literal))
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseCharLiteral() ast.Expression {
	return &ast.CharLiteral{Token: p.curToken, Value: []rune(p.curToken.Literal)[0]}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == token.BOOL && p.curToken.Literal == "true"}
}

func (p *Parser) parseNilLiteral() ast.Expression {
	return &ast.NilLiteral{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parsePointerReferenceExpression() ast.Expression {
	exp := &ast.PointerReferenceExpression{Token: p.curToken}
	p.nextToken()
	exp.Value = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parsePointerDereferenceExpression() ast.Expression {
	exp := &ast.PointerDereferenceExpression{Token: p.curToken}
	p.nextToken()
	exp.Value = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)
	p.nextToken()
	exp.Consequence = p.parseBlockStatement()
	if p.curTokenIs(token.ELSE) {
		p.nextToken()
		exp.Alternative = p.parseBlockStatement()
	}
	if !p.curTokenIs(token.END) {
		p.peekError(token.END)
		return nil
	}
	return exp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	p.nextToken()
	lit.Body = p.parseBlockStatement()
	if !p.curTokenIs(token.END) {
		p.peekError(token.END)
		return nil
	}
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}
	p.nextToken()
	params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return params
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.curToken}
	arr.Elements = p.parseExpressionList(token.RBRACKET)
	return arr
}

func (p *Parser) parseMapLiteral() ast.Expression {
	m := &ast.MapLiteral{Token: p.curToken}
	m.Pairs = make(map[ast.Expression]ast.Expression)
	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		val := p.parseExpression(LOWEST)
		m.Pairs[key] = val
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return m
}

// Helper to parse comma-separated lists (arrays, arguments)
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// --- Infix Parsing Functions ---

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: fn}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseFieldAccessExpression(left ast.Expression) ast.Expression {
	exp := &ast.FieldAccessExpression{Token: p.curToken, Object: left}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	exp.Field = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return exp
}

func (p *Parser) parseStructInstantiation(left ast.Expression) ast.Expression {
	nameIdent, ok := left.(*ast.Identifier)
	if !ok {
		p.errors = append(p.errors, "struct instantiation expects identifier")
		return nil
	}
	stmt := &ast.StructInstantiationExpression{Token: p.curToken, Name: nameIdent}
	stmt.Fields = []ast.StructField{}

	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return stmt
	}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			return nil
		}
		field := ast.StructField{Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}}
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		field.Value = p.parseExpression(LOWEST)
		stmt.Fields = append(stmt.Fields, field)
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return stmt
}
