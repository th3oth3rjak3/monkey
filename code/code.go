package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Opcode is an instruction to perform.
type Opcode byte

const (
	OpConstant      Opcode = iota // represents a constant value
	OpPop                         // Tells the vm to pop the stack.
	OpAdd                         // Represents an addition operation
	OpSub                         // Represents a subtraction operation
	OpMul                         // Represents a multiplication operation
	OpDiv                         // Represents a division operation
	OpTrue                        // Represents the value true
	OpFalse                       // Represents the value false
	OpEqual                       // Represents equality
	OpNotEqual                    // Represents inequality
	OpGreaterThan                 // Represents left > right
	OpMinus                       // Represents integer negation
	OpBang                        // Represents boolean negation
	OpJump                        // Represents a Jump
	OpJumpNotTruthy               // Represents a jump when a condition is false
	OpNull                        // Represents no value
	OpGetGlobal                   // Get a global binding
	OpSetGlobal                   // Set a global binding
	OpArray                       // Construct an array from N elements off of the stack
	OpHash                        // Construct a hash
	OpIndex                       // Index expression
	OpCall                        // Call a function
	OpReturnValue                 // Return a value from the function call by popping the last value on the stack
	OpReturn                      // Return from a function, no value is on the stack
	OpGetLocal                    // Get a local binding
	OpSetLocal                    // Set a local binding
	OpGetBuiltin                  // Get a builtin function
)

// Instructions represent virtual machine instructions.
type Instructions []byte

// Implement the string interface.
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0

	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return fmt.Sprintf("%s", def.Name)
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Definition is a way to keep track opcode metadata
type Definition struct {
	Name          string // The human readable name for an opcode.
	OperandWidths []int  // Contains the number of bytes each operand takes up.
}

// definitions contains a map of all the opcode types to some metadata about their use.
var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpPop:           {"OpPop", []int{}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJump:          {"OpJump", []int{2}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{}},
	OpCall:          {"OpCall", []int{1}},
	OpReturnValue:   {"OpReturnValue", []int{}},
	OpReturn:        {"OpReturn", []int{}},
	OpGetLocal:      {"OpGetLocal", []int{1}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
	OpGetBuiltin:    {"OpGetBuiltin", []int{1}},
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

// Make creates an instruction which consists of an Opcode and all of its operands.
//
// Parameters:
//   - op: The Opcode which represents the operation to perform.
//   - operands: Any number of operands required for the operation.
//
// Returns:
//   - []byte: The instruction to execute encoded in a byte slice.
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}
