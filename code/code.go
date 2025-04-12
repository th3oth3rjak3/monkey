package code

import "fmt"

// Instructions represent virtual machine instructions.
type Instructions []byte

// Opcode is an instruction to perform.
type Opcode byte

const (
	OpConstant Opcode = iota // represents a constant value
)

// Definition is a way to keep track opcode metadata
type Definition struct {
	Name          string // The human readable name for an opcode.
	OperandWidths []int  // Contains the number of bytes each operand takes up.
}

// definitions contains a map of all the opcode types to some metadata about their use.
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup is used to access opcode definitions from other packages.
//
// Parameters:
//   - op: The opcode to lookup information for.
//
// Returns:
//   - *Definition: The definition of the opcode
//   - error: An error that would be returned if the opcode is not found.
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}
