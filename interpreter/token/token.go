package token

// TokenType represents the type of a token.
type TokenType string

// Token represents a lexer token in the Monkey language.
type Token struct {
	Type    TokenType // The type of the token.
	Literal string    // The literal value of the token.
}

const (
	// Token Flow Control
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT = "IDENT" // add, x, y, z
	INT   = "INT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"
	EQ = "=="
	NE = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Grouping
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	STRING   = "STRING"
)

// keywords is the map of reserved keywords to their token types.
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent returns the token type for the given identifier.
// It returns IDENT if the identifier is not a reserved keyword because it must be user defined.
//
// Parameters:
//   - ident: The identifier to look up.
//
// Returns:
//   - TokenType: The token type of the identifier.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
