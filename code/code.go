package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Opcode is an instruction to perform.
type Opcode byte

const (
	OpConstant Opcode = iota // represents a constant value
	OpPop                    // Tells the vm to pop the stack.
	OpAdd                    // Represents an addition operation
	OpSub                    // Represents a subtraction operation
	OpMul                    // Represents a multiplication operation
	OpDiv                    // Represents a division operation
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
	OpConstant: {"OpConstant", []int{2}},
	OpPop:      {"OpPop", []int{}},
	OpAdd:      {"OpAdd", []int{}},
	OpSub:      {"OpSub", []int{}},
	OpMul:      {"OpMul", []int{}},
	OpDiv:      {"OpDiv", []int{}},
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
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
