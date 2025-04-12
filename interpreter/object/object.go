package object

import (
	"bytes"
	"fmt"
	"monkey/interpreter/ast"
	"strings"
)

// ObjectType represents the underlying object's type.
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
)

// Object represents our universal type.
type Object interface {
	Type() ObjectType // Type gets the underlying object type.
	Inspect() string  // Inspect represents the object as a string.
}

// Integer represents an integer type.
type Integer struct {
	Value int64 // The value of the integer.
}

// Inspect represents the object as a string.
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type gets the underlying object type.
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// Boolean represents a boolean value.
type Boolean struct {
	Value bool // The actual value.
}

// Inspect represents the object as a string.
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type gets the underlying object type.
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// Null represents no value.
type Null struct{}

// Inspect represents the object as a string.
func (n *Null) Inspect() string {
	return "null"
}

// Type gets the underlying object type.
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// ReturnValue represents a value that should be returned from a block or function.
type ReturnValue struct {
	Value Object // The object that should be returned.
}

// Type gets the underlying object type.
func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

// Inspect represents the object as a string.
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

// Error represents an error that occurs during interpretation.
type Error struct {
	Message string // The error message.
}

// Type gets the underlying object type.
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

// Inspect represents the object as a string.
func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s", e.Message)
}

// Function represents a callable function.
type Function struct {
	Parameters []*ast.Identifier   // The parameters that were passed to the function.
	Body       *ast.BlockStatement // The block of statements to execute.
	Env        *Environment        // The environment containing the current state.
}

// Type gets the underlying object type.
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

// Inspect represents the object as a string.
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
