package repl

import (
	"bufio"
	"crabscript.rs/compiler"
	"crabscript.rs/lexer"
	"crabscript.rs/parser"
	"crabscript.rs/vm"
	"fmt"
	"io"
)

const Prompt = "=> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	// env := object.NewEnvironment()

	for {
		fmt.Printf(Prompt)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Compilation failed: %s", err)
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Bytecode failed to execute: %s", err)
			continue
		}
		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		_, ok := io.WriteString(out, err)
		if ok != nil {
			panic("Cannot print error!")
		}
		_, ok = io.WriteString(out, "\n")
		if ok != nil {
			panic("Cannot print error!")
		}
	}
}
