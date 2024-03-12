# Crabscript 🦀
A complete interpreter and compiler for a language based on Thorsten Ball's "Writing an Interpreter in Go" using only 
the standard library. The current implemented features are:

## Interpreter

- [x] Tokeniser (AST)
- [x] Lexer
- [x] Parser
- [x] Evaluator
- [x] REPL
- [x] Files
- [x] Int
- [x] String
- [x] Bool
- [x] Variable binding
- [x] Functions
- [x] Closures
- [x] Arrays
- [x] Builtins (len, first, last, tail, puts)

## Compiler

- [x] Compiler
- [x] Virtual Machine

## About
The parser is using [Pratt's algorithm](https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html), 
which is modular and easily extensible.
The evaluator is an implementation of a tree-walking interpreter and no 
byte-code is generated. There are no primitive types - everything is an object 
a la Ruby.

The compiler is less complete than the interpreter - it doesn't fully support fns.
