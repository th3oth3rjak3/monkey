package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression               // A function that is used to parse an expression with a prefix operator.
	infixParseFn  func(ast.Expression) ast.Expression // A function that is used to parse an expression with an infix operator.
)

const (
	_           int = iota
	LOWEST          // The lowest precedence possible
	EQUALS          // ==
	LESSGREATER     // < or >
	SUM             // +
	PRODUCT         // *
	PREFIX          // -x or !x
	CALL            // myFunction(x)
	INDEX           // array[index]
)

// precedence defines the operator precedence for each given token type.
var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NE:       EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

// Parser represents the monkey language parser to convert tokens into a runnable program.
type Parser struct {
	lex            *lexer.Lexer                      // The lexer that is used for generating tokens.
	errors         []string                          // The list of parse errors encountered.
	curToken       token.Token                       // The current token to be parsed.
	peekToken      token.Token                       // The next token to be parsed.
	prefixParseFns map[token.TokenType]prefixParseFn // The map of token types to a predetermined prefix parsing function.
	infixParseFns  map[token.TokenType]infixParseFn  // The map of token types to a predetermined infix parsing function.
}

// New creates a new Parser.
//
// Parameters:
//   - l: A lexer used to handle the source code lexing.
//
// Returns:
//   - *Parser: a new parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lex: l, errors: []string{}}

	// Read two tokens so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	// Register prefix parsing functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.TRUE, p.parseBoolean)
	p.registerPrefixFn(token.FALSE, p.parseBoolean)
	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixFn(token.IF, p.parseIfExpression)
	p.registerPrefixFn(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefixFn(token.STRING, p.parseStringLiteral)
	p.registerPrefixFn(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefixFn(token.LBRACE, p.parseHashLiteral)

	// Register infix parsing functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NE, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.LPAREN, p.parseCallExpression)
	p.registerInfixFn(token.LBRACKET, p.parseIndexExpression)

	return p
}

// registerPrefixFn adds an associated parsing function for a given prefix token type with the parser.
func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfixFn adds an associated parsing function for a given infix token type with the parser.
func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Errors returns any errors encountered during parsing.
//
// Returns:
//   - []string: The errors encountered during parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError creates an error message and adds it to the parser's error list.
//
// Parameters:
//   - t: The expected token type.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, but got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken uses the lexer to get the next token and update its internal state.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

// ParseProgram uses the parser to parse the source code.
//
// Returns:
//   - *ast.Program: The program code represented as an AST.
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

// parseStatement parses a statement from the input source code.
//
// Returns:
//   - ast.Statement: a statement generated by parsing.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement evaluates tokens and constructs a let statement from the input.
//
// Returns:
//   - *ast.LetStatement: a statement that represents a let binding.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// curTokenIs checks to see if the current token's type matches the expected input.
//
// Parameters:
//   - t: The token type that is expected.
//
// Returns:
//   - bool: true when the input matches the current token type, otherwise false.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks to see if the peek token's type matches the expected input.
//
// Parameters:
//   - t: The token type that is expected.
//
// Returns:
//   - bool: true when the input matches the peek token type, otherwise false.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks the peek token to see if it matches the expected type.
// If so, it advances the tokens on the parser state and returns true.
//
// Returns:
//   - bool: true when the peek token matches the expected input TokenType, otherwise false.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// parseReturnStatement constructs a return statement at the current parser position.
//
// Returns:
//   - *ast.ReturnStatement: The return statement parsed from the current parser position.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement constructs an Expression Statement at the current parser position.
//
// Returns:
//   - *ast.ExpressionStatement: The new expression statement.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpression parses the expression with the given precedence level.
//
// Parameters:
//   - precedence: The precedence to use for determining parsing order. (higher is more important)
//
// Returns:
//   - ast.Expression: The expression parsed from the current token position.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseIdentifier parses the current token as an Identifier.
//
// Returns:
//   - ast.Expression: The expression parsed from the current token position.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral parses the current token as an integer.
//
// Returns:
//   - ast.Expression: The integer literal expression parsed from the current token position.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}

	lit.Value = value
	return lit
}

// parseStringLiteral parses the current token as a literal string.
//
// Returns:
//   - ast.Expression: The string literal expression parsed from the current token position.
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// noPrefixParseFnError adds an error message to the parser error list when no registered parsing function was found.
//
// Parameters:
//   - t: The token type that was missing the registered prefixParseFn.
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parsePrefixExpression creates a prefix expression by parsing the operator and expression.
//
// Returns:
//   - ast.Expression: The expression parsed as a prefix expression.
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// peekPrecedence finds the precedence of the parsers peekToken.
//
// Returns:
//   - int: The precedence for the peekToken.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// curPrecedence finds the precedence for the parser's curToken.
//
// Returns:
//   - int: The precedence for the curToken.
func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// parseInfixExpression constructs an infix expression.
//
// Parameters:
//   - left: An expression used as the left hand side of the infix expression.
//
// Returns:
//   - ast.Expression: The resulting infix expression.
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

// parseBoolean parses a boolean expression.
//
// Returns:
//   - ast.Expression: A boolean ast expression.
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parseGroupedExpression parses expression that are grouped.
//
// Returns:
//   - ast.Expression: The grouping ast expression.
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression parses an if expression into its separate parts.
//
// Returns:
//   - ast.Expression: the parsed if expression.
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// parseBlockStatement parses a block statement like the Consequence of an if expression.
//
// Returns:
//   - *ast.BlockStatement: The block statement after parsing.
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

	return block
}

// parseFunctionLiteral parses a user defined function.
//
// Returns:
//   - ast.Expression: The function literal as an expression.
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

// parseFunctionParameters handles the parsing task for function input parameters.
//
// Returns:
//   - []*ast.Identifier: A slice of identifiers that are function parameters.
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

// parseCallExpression handles parsing for a callable expression.
//
// Parameters:
//   - function: The input function that is to be called.
//
// Returns:
//   - ast.Expression: The parsed call expression.
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

// parseArrayLiteral handles the parsing of an array of values.
//
// Returns:
//   - ast.Expression: The output array expression.
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// parseExpressionList parses a series of expressions joined by commas as a go slice.
//
// Parameters:
//   - end: The token that is expected at the end of the list which would mean we are done.
//
// Returns:
//   - []ast.Expression: The list of expressions that were parsed.
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

// parseIndexExpression parses an index expression like items[0].
//
// Parameters:
//   - left: The left side of the expression. e.g. items in items[0].
//
// Returns:
//   - ast.Expression: The parsed index expression.
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
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
