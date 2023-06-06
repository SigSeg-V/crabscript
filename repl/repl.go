package repl

import (
	"bufio"
	"fmt"
	"io"
	"crabscript.rs/lexer"
	"crabscript.rs/token"
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
		
		for tok := l.NextToken(); tok.Type != token.Eof; tok = l.NextToken(){
			fmt.Printf("%+v\n", tok)
		}
	}
}