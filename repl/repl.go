package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

// Start begins the REPL for the user.
//
// Parameters:
//   - in: The input reader interface.
//   - out: The output writer interface.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// env := object.NewEnvironment()
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Compilation failed:\n%s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed:\n%s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")

		// evaluated := evaluator.Eval(program, env)
		// if evaluated != nil {
		// 	io.WriteString(out, evaluated.Inspect())
		// 	io.WriteString(out, "\n")
		// }
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

// printParserErrors prints the errors out to the writer for the user.
//
// Parameters:
//   - out: The output writer.
//   - errors: The errors to print to the output.
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Whoops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
