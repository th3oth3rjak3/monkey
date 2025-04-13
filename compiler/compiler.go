package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
	"sort"
)

// Compiler represents the bytecode compiler.
type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction // The instruction that was most recently emitted to the compiler
	previousInstruction EmittedInstruction // The instruction that was emitted just before the lastInstruction
	symbolTable         *SymbolTable       // Where identifiers are stored.
}

// Bytecode represents instructions for our bytecode vm.
type Bytecode struct {
	Instructions code.Instructions // Instructions to execute
	Constants    []object.Object   // Constants to reference from the instructions by position number.
}

// EmittedInstruction is a record of a previously emitted instruction
type EmittedInstruction struct {
	Opcode   code.Opcode // The opcode that was emitted.
	Position int         // The position in the instructions list where the instruction was emitted to.
}

// New creates a new compiler instance.
//
// Returns:
//   - *Compiler: The new compiler instance.
func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.constants = constants
	compiler.symbolTable = s
	return compiler
}

// Compile produces bytecode from our ast.
//
// Parameters:
//   - node: The root node of the program to compile.
//
// Returns:
//   - error: An error is returned when compilation fails.
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "-":
			c.emit(code.OpMinus)
		case "!":
			c.emit(code.OpBang)
		}

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
		} else {
			err := c.Compile(node.Left)
			if err != nil {
				return err
			}

			err = c.Compile(node.Right)
			if err != nil {
				return err
			}

			switch node.Operator {
			case "+":
				c.emit(code.OpAdd)
			case "-":
				c.emit(code.OpSub)
			case "*":
				c.emit(code.OpMul)
			case "/":
				c.emit(code.OpDiv)
			case ">":
				c.emit(code.OpGreaterThan)
			case "==":
				c.emit(code.OpEqual)
			case "!=":
				c.emit(code.OpNotEqual)
			default:
				return fmt.Errorf("unknown operator %s", node.Operator)
			}
		}
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an OpJumpNotTruthy with a bogus value 9999
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		// emit an `OpJump` with a bogus value to be updated later
		opJumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			// compile the alternative
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.instructions)
		c.changeOperand(opJumpPos, afterAlternativePos)

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, e := range node.Elements {
			err := c.Compile(e)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2) // each pair has a key and value, so its pairs * 2 elements.

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.emit(code.OpGetGlobal, symbol.Index)

	// Statements
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)

		// Placeholder to keep this bracket down.
	}

	return nil
}

// Bytecode outputs the bytecode generated by the compiler.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// addConstant appends a constant value to the end of the compiler's constants slice.
//
// Parameters:
//   - obj: the constant object to add the the compiler's constant slice.
//
// Returns:
//   - int: The index position of the newly added object.
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit adds a new opcode and operands to the compiler instructions
//
// Parameters:
//   - op: The opcode instruction to add.
//   - operands: The operands to add to the instruction.
//
// Returns:
//   - int: Essentially a stack pointer to the start of the current instruction.
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setLastInstruction(op, pos)
	return pos
}

// addInstruction adds an instruction to the compiler and returns the start
// position for the new instruction.
//
// Parameters:
//   - ins: The instruction to add
//
// Returns:
//   - int: The starting position for the current instruction.
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

// setLastInstruction keeps track of the instructions that were recently emited in the compiler.
//
// Parameters:
//   - op: The opcode most recently emitted.
//   - pos: The position in the instruction set where this was emitted.
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

// lastInstructionIsPop evaluates if the last emitted instruction was an OpPop.
func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

// removeLastPop removes the last pop instruction from the stack when it's not needed.
func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

// replaceInstruction allows us to replace an instruction with a set of new instructions.
// WARNING: This function assumes that the newInstruction is the same width as the one being replaced.
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

// changeOperand allows us to update the operand of an existing instruction.
//
// Parameters:
//   - opPos: The position of the opcode in the instruction set.
//   - operand: The new operand.
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}
