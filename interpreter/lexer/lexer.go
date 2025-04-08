package lexer

import (
	"monkey/interpreter/token"
	"slices"
)

// Lexer represents a lexical analyzer for the Monkey language.
type Lexer struct {
	input        string // The source code being lexed.
	position     int    // The current position in the input (points to current char).
	readPosition int    // The next position in the input (after current char).
	ch           byte   // The current char under examination.
}

// New creates a new Lexer instance.
// It initializes the Lexer with the input string.
//
// Parameters:
//   - input: The source code string to be lexed.
//
// Returns:
//   - *Lexer: A pointer to the newly created Lexer instance.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character from the input and updates the Lexer's state.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken evaluates the current character and returns the corresponding token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NE, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// skipWhitespace is a function that is used to iterate over any
// whitespace characters to exclude them from token generation.
func (l *Lexer) skipWhitespace() {
	whitespace := []byte{' ', '\r', '\t', '\n'}
	for slices.Contains(whitespace, l.ch) {
		l.readChar()
	}
}

// readIdentifier reads characters until it finds no more letters.
// It then produces the string identifier between the two positions.
//
// Returns:
//   - string: A new user-defined identifier.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads characters until it finds no more numeric characters.
// It then produces the string representation of the number between the two
// positions.
//
// Returns:
//   - string: The string representation of a user provided number.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// newToken is a helper function to generate a Token.
// This is used to generate single character operator type tokens.
//
// Example: '=' for Assignment.
//
// Parameters:
//   - tokenType: The TokenType of the new token.
//   - ch: The character that represents a single character token.
//
// Returns:
//   - Token: A new token.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// isLetter checks to see if a character is a letter or an underscore.
// These values are considered valid when used by an identifier.
//
// Parameters:
//   - ch: the character to check.
//
// Returns:
//   - bool: True when the character is a letter, otherwise false.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks to see if a character is considered to be numeric.
// This is typically used to parse a literal number.
//
// Parameters:
//   - ch: the character to check.
//
// Returns:
//   - bool: True when the character is numeric, otherwise false.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// peekChar returns the next character in the source input.
//
// Returns:
//   - byte: The next character beyond the current read position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
