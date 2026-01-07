// ==============================================================================================
// FILE: parser/parser.go
// ==============================================================================================
// PACKAGE: parser
// PURPOSE: Transforms the linear stream of tokens into a hierarchical Abstract Syntax Tree (AST).
//          It implements a Recursive Descent Pratt Parser to handle operator precedence
//          and complex nesting efficiently.
// ==============================================================================================

package parser

import (
	"fmt"

	"eloquence/ast"
	"eloquence/lexer"
	"eloquence/token"
)

// ----------------------------------------------------------------------------------------------
// OPERATOR PRECEDENCE
// ----------------------------------------------------------------------------------------------

const (
	_ int = iota
	LOWEST
	EQUALS      // ==, !=
	LESSGREATER // >, <, >=, <=
	SUM         // +, -
	PRODUCT     // *, /, %
	PREFIX      // -X, !X
	CALL        // myFunction(X)
	INDEX       // array[index], struct.field
)

// precedences maps token types to their parsing priority.
var precedences = map[token.TokenType]int{
	token.EQUALS:        EQUALS,
	token.NOT_EQUALS:    EQUALS,
	token.LESS:          LESSGREATER,
	token.GREATER:       LESSGREATER,
	token.LESS_EQUAL:    LESSGREATER,
	token.GREATER_EQUAL: LESSGREATER,
	token.ADDS:          SUM,
	token.SUBTRACTS:     SUM,
	token.MINUS:         SUM,
	token.TIMES:         PRODUCT,
	token.DIVIDES:       PRODUCT,
	token.MODULO:        PRODUCT,
	token.LPAREN:        CALL,
	token.LBRACKET:      INDEX,
	token.DOT:           INDEX,
	token.AND:           LESSGREATER,
	token.OR:            LESSGREATER,
}

// Function types for Pratt Parsing
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser represents the state of the parsing process.
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken token.Token

	// peekTokens is a fixed-size buffer for lookahead.
	// We need 3 tokens of lookahead to distinguish between:
	// 1. "User { key : val }" (Struct Instantiation)
	// 2. "while x < y { show... }" (Block start)
	peekTokens [3]token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New initializes the parser and fills the lookahead buffer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Register Prefix Parsers (for tokens that start an expression)
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.CHAR, p.parseCharLiteral)
	p.registerPrefix(token.BOOL, p.parseBooleanLiteral)
	p.registerPrefix(token.NIL, p.parseNilLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.TAKES, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral) // Maps { key: val }
	p.registerPrefix(token.POINTING_TO, p.parsePointerReference)
	p.registerPrefix(token.POINTING_FROM, p.parsePointerDereference)

	// Register Infix Parsers (for tokens that sit between expressions)
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
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseFieldAccessExpression)

	// Fill lookahead buffer
	p.peekTokens[0] = p.l.NextToken()
	p.peekTokens[1] = p.l.NextToken()
	p.peekTokens[2] = p.l.NextToken()

	// Load first token into current
	p.nextToken()

	return p
}

// ----------------------------------------------------------------------------------------------
// TOKEN MANAGEMENT HELPERS
// ----------------------------------------------------------------------------------------------

// nextToken shifts the lookahead window forward.
func (p *Parser) nextToken() {
	p.curToken = p.peekTokens[0]
	p.peekTokens[0] = p.peekTokens[1]
	p.peekTokens[1] = p.peekTokens[2]
	p.peekTokens[2] = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekTokens[0].Type == t
}

// peekTokenAt allows inspecting tokens deeper in the buffer.
// 0 = next token, 1 = token after next, etc.
func (p *Parser) peekTokenAt(i int) token.Token {
	if i >= 0 && i < 3 {
		return p.peekTokens[i]
	}
	return token.Token{Type: token.EOF}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("line %d:%d - expected next token to be %s, got %s instead",
		p.peekTokens[0].Line, p.peekTokens[0].Column, t, p.peekTokens[0].Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// ----------------------------------------------------------------------------------------------
// STATEMENT PARSING
// ----------------------------------------------------------------------------------------------

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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.RETURN:
		return p.parseReturnStatement()
	case token.DEFINE:
		return p.parseStructDefinition()
	case token.WHILE, token.REPEAT:
		return p.parseLoopStatement()
	case token.FOR:
		return p.parseRangeLoopStatement()
	case token.POINTING_FROM:
		if p.isPointerAssignment() {
			return p.parsePointerAssignment()
		}
		return p.parseExpressionStatement()
	case token.TRY:
		return p.parseTryCatchStatement()
	case token.INCLUDE:
		return p.parseIncludeStatement()
	default:
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.IS) {
			return p.parseAssignmentStatement()
		}
		return p.parseExpressionStatement()
	}
}

func (p *Parser) isPointerAssignment() bool {
	if p.peekTokenIs(token.IDENT) && p.peekTokenAt(1).Type == token.IS {
		return true
	}
	return false
}

func (p *Parser) parsePointerAssignment() ast.Statement {
	stmt := &ast.PointerAssignmentStatement{Token: p.curToken}
	p.nextToken() // eat 'pointing from'

	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.IS) {
		return nil
	}
	p.nextToken() // eat 'is'

	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.IS) {
		return nil
	}
	p.nextToken() // eat 'is'

	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// Allow return; (void return)
	if p.curTokenIs(token.EOF) || p.curTokenIs(token.END) || p.curTokenIs(token.ELSE) {
		stmt.ReturnValue = nil
		return stmt
	}
	stmt.ReturnValue = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseStructDefinition() *ast.StructDefinitionStatement {
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
	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		stmt.Attributes = append(stmt.Attributes, ident)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return stmt
}

func (p *Parser) parseLoopStatement() *ast.LoopStatement {
	stmt := &ast.LoopStatement{Token: p.curToken}
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseRangeLoopStatement() *ast.RangeLoopStatement {
	stmt := &ast.RangeLoopStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Iterator = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	stmt.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseTryCatchStatement() *ast.TryCatchStatement {
	stmt := &ast.TryCatchStatement{Token: p.curToken}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	stmt.TryBlock = p.parseBlockStatement()

	if p.peekTokenIs(token.CATCH) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		stmt.CatchBlock = p.parseBlockStatement()
	}

	if p.peekTokenIs(token.FINALLY) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		stmt.FinallyBlock = p.parseBlockStatement()
	}
	return stmt
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}
	p.nextToken()
	stmt.Path = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	// CHECK FOR UNTERMINATED BLOCK:
	// If we hit EOF instead of RBRACE, we report an error.
	if p.curTokenIs(token.EOF) {
		p.errors = append(p.errors, "unterminated block: expected '}', got EOF")
	}

	return block
}

// ----------------------------------------------------------------------------------------------
// EXPRESSION PARSING
// ----------------------------------------------------------------------------------------------

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekTokens[0].Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekTokens[0].Type]; ok {
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

// ----------------------------------------------------------------------------------------------
// PREFIX FUNCTIONS
// ----------------------------------------------------------------------------------------------

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.LBRACE) {
		peek1 := p.peekTokenAt(1)
		peek2 := p.peekTokenAt(2)

		if peek1.Type == token.RBRACE {
			return p.parseStructInstantiation(ident)
		}
		if peek1.Type == token.IDENT && peek2.Type == token.COLON {
			return p.parseStructInstantiation(ident)
		}
	}
	return ident
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	var value int64
	fmt.Sscanf(p.curToken.Literal, "%d", &value)
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}
	var value float64
	fmt.Sscanf(p.curToken.Literal, "%f", &value)
	lit.Value = value
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
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
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
	expression := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.peekTokenIs(token.LBRACE) && !p.peekTokenIs(token.IF) {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
		} else {
			if p.peekTokenIs(token.LBRACE) {
				p.nextToken()
			}
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.MapLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

// ----------------------------------------------------------------------------------------------
// INFIX FUNCTIONS
// ----------------------------------------------------------------------------------------------

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

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

func (p *Parser) parseStructInstantiation(name *ast.Identifier) ast.Expression {
	p.nextToken()
	exp := &ast.StructInstantiationExpression{Token: p.curToken, Name: name}
	exp.Fields = []ast.StructField{}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			return nil
		}
		fieldKey := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		fieldVal := p.parseExpression(LOWEST)
		exp.Fields = append(exp.Fields, ast.StructField{Name: fieldKey, Value: fieldVal})

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return exp
}

func (p *Parser) parsePointerReference() ast.Expression {
	exp := &ast.PointerReferenceExpression{Token: p.curToken}
	p.nextToken()
	exp.Value = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parsePointerDereference() ast.Expression {
	exp := &ast.PointerDereferenceExpression{Token: p.curToken}
	p.nextToken()
	exp.Value = p.parseExpression(PREFIX)
	return exp
}
