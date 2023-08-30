package repl

import (
	"bufio"
	"crabscript.rs/evaluator"
	"crabscript.rs/lexer"
	"crabscript.rs/parser"
	"fmt"
	"io"
)

const Prompt = "=> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
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
