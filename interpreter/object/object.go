package object

import "fmt"

// ObjectType represents the underlying object's type.
type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
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
func (n *Null) Type() string {
	return NULL_OBJ
}
