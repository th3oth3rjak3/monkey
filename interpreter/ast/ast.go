package ast

import (
	"bytes"
	"monkey/interpreter/token"
)

// Node represents a node in the abstract syntax tree.
type Node interface {
	TokenLiteral() string
	String() string
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

// String returns a string representation of the program statements
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement represents a binding from a user defined identifier to an expression.
type LetStatement struct {
	Token token.Token // The original token.LET token.
	Name  *Identifier // The identifier given by the user.
	Value Expression  // The Expression that is represented by the Name.
}

// statementNode is a placeholder function for the Statement interface.
func (l *LetStatement) statementNode() {}

// TokenLiteral returns the literal value of the token in the let statement.
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

// String returns a string representation of the LetStatement
func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Identifier represents a user defined name that is mapped to an expression through a let binding.
type Identifier struct {
	Token token.Token // The original token.IDENT token.
	Value string      // The string literal value with the identifier name.
}

// statementNode is a placeholder function for the Statement interface.
func (i *Identifier) expressionNode() {}

// TokenLiteral returns the literal value of the token of the identifier.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of the Identifier
func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement represents an expression to be evaluated and returned from a function.
type ReturnStatement struct {
	Token       token.Token // the 'return' token.
	ReturnValue Expression  // The expression to be returned.
}

// statementNode is a placeholder function for the Statement interface.
func (r *ReturnStatement) statementNode() {}

// TokenLiteral returns the literal value of the token of the return statement.
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

// String returns a string representation of the ReturnStatement
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

// ExpressionStatement is a simple wrapper around an expression to type coerce it into a Program.
type ExpressionStatement struct {
	Token      token.Token // The first token of the expression
	Expression Expression  // The expression
}

// statementNode is a placeholder function for the Statement interface.
func (e *ExpressionStatement) statementNode() {}

// TokenLiteral returns the literal value of the token of the expression statement.
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

// String returns a string representation of the ExpressionStatement
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

// IntegerLiteral represents a node that is an integer value.
type IntegerLiteral struct {
	Token token.Token // The token for the integer literal.
	Value int64       // The actual literal value.
}

// expressionNode is a placeholder function for the Expression interface.
func (i *IntegerLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token for the integer literal.
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of the IntegerLiteral
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}
