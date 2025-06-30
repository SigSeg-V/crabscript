package main

import (
	"crabscript.rs/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {

	handleInput(os.Args)

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome, %v, to 🦀script!\n", user.Name)
	repl.Start(os.Stdin, os.Stdout)
}

// returns code based on option
// 0 - compile an object
// 1 - repl
// 2 - interpret
func handleInput(args []string) (int, error) {
	if len(args) == 1 {
		return 1, nil
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" {
			if i+1 >= len(args) {
				return -1, fmt.Errorf("missing output file name")
			} else {

			}
		}
	}
	return -1, fmt.Errorf("unable to parse options")
}
