package main

import (
	"fmt"
	"os"
	"os/user"
	"crabscript.rs/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome, %v, to 🦀script!\n", user.Name)
	repl.Start(os.Stdin, os.Stdout)
}
