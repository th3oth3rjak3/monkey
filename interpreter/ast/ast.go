package ast

// Node represents a node in the abstract syntax tree.
type Node interface {
	TokenLiteral() string
}

// Statement represents a statement in the abstract syntax tree.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression in the abstract syntax tree.
type Expression interface {
	Node
	expressionNode()
}

// Program represents a program in the abstract syntax tree.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the token literal of the first statement in the program.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
