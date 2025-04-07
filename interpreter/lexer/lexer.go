package lexer

import (
	"monkey/interpreter/token"
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
	l.ReadChar()
	return l
}

// ReadChar reads the next character from the input and updates the Lexer's state.
func (l *Lexer) ReadChar() {
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

	// TODO: left off on page 23 of the book. I was about to add more lexer features.
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.ReadChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
