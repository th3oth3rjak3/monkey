package ast

import (
	"bytes"
	"monkey/interpreter/token"
	"strings"
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

// PrefixExpression represents an expression that has a prefix operator.
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. ! or -
	Operator string      // The string representation of the operator.
	Right    Expression  // The expression to be evaluated.
}

// expressionNode is a placeholder function for the Expression interface.
func (p *PrefixExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token for the prefix expression.
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// String returns a string representation of the PrefixExpression
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an expression that is bound by an infix operator.
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression  // The left hand side expression to be evaluated.
	Operator string      // The operator literal value.
	Right    Expression  // The right hand side expression to be evaluated.
}

// expressionNode is a placeholder function for the Expression interface.
func (i *InfixExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token for the infix expression.
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of the InfixExpression
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

// Boolean represents a boolean literal.
type Boolean struct {
	Token token.Token // The token that represents the boolean literal.
	Value bool        // The actual value of the boolean: true or false.
}

// expressionNode is a placeholder function for the Expression interface.
func (b *Boolean) expressionNode() {}

// TokenLiteral returns the literal value of the token for the boolean expression.
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String returns a string representation of the boolean expression
func (b *Boolean) String() string {
	return b.Token.Literal
}

// IfExpression represents a decision to make based on a conditional expression.
// When true, Consequence is evaluated, when false Alternative is evaluated.
type IfExpression struct {
	Token       token.Token     // The 'if' Token
	Condition   Expression      // The condition to evaluate
	Consequence *BlockStatement // The result when the condition is true.
	Alternative *BlockStatement // The result when the condition is false.
}

// expressionNode is a placeholder function for the Expression interface.
func (i *IfExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token for the infix expression.
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of the if expression
func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

// BlockStatement represents a scoped section of code that contains
// additional statements.
type BlockStatement struct {
	Token      token.Token // The { token
	Statements []Statement // A collection of scoped statements.
}

// statementNode is a placeholder function for the Statement interface.
func (b *BlockStatement) statementNode() {}

// TokenLiteral returns the literal value of the token for the block statement.
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

// String returns a string representation of the block statement.
func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// FunctionLiteral represents a user defined function.
type FunctionLiteral struct {
	Token      token.Token     // The 'fn' token
	Parameters []*Identifier   // The list of parameters which can be empty.
	Body       *BlockStatement // The body of the function
}

// expressionNode is a placeholder function for the Expression interface.
func (f *FunctionLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token for the function literal.
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

// String returns a string representation of the function literal.
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}

// CallExpression represents a function and a set of arguments that can be called.
type CallExpression struct {
	Token     token.Token  // The '(' Token
	Function  Expression   // The function to be called
	Arguments []Expression // The arguments to be passed into the function.
}

// expressionNode is a placeholder function for the Expression interface.
func (c *CallExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token for the call expression.
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

// String returns a string representation of the call expression.
func (c *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}

	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// StringLiteral represents a string in the monkey language.
type StringLiteral struct {
	Token token.Token // The token.
	Value string      // The literal string value.
}

// expressionNode is a placeholder function for the Expression interface.
func (s *StringLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token for the StringLiteral expression.
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

// String returns a string representation of the StringLiteral expression.
func (s *StringLiteral) String() string {
	return s.Token.Literal
}

// ArrayLiteral represents an array of data.
type ArrayLiteral struct {
	Token    token.Token  // The '[' token.
	Elements []Expression // A collection of expressions.
}

// expressionNode is a placeholder function for the Expression interface.
func (a *ArrayLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the token for the ArrayLiteral expression.
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

// String returns a string representation of the ArrayLiteral expression.
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// IndexExpression represents using indexing to get a value from an array. Example: items[1]
type IndexExpression struct {
	Token token.Token // The '[' Token
	Left  Expression  // The identifier, function, or array that evaluates to an array of items.
	Index Expression  // The index expression that is used to find the item in the list. e.g. 1 + 1 in items[1 + 1]
}

// expressionNode is a placeholder function for the Expression interface.
func (i *IndexExpression) expressionNode() {}

// TokenLiteral returns the literal value of the token for the IndexExpression.
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

// String returns a string representation of the ArrayLiteral expression.
func (i *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")

	return out.String()
}
