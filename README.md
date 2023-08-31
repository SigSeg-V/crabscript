# Crabscript 🦀
A language based on Thorsten Ball's "Writing an Interpreter in Go".
The current implemented features are:

- [x] Tokeniser (AST)
- [x] Lexer
- [x] Parser
- [x] Evaluator
- [x] REPL
- [ ] Files
- [x] Int
- [x] String
- [x] Bool
- [x] Variable binding
- [x] Functions
- [x] Closures
- [ ] Arrays
- [ ] Builtins

The parser is using Pratt's algorithm, which is modular and easily extensible.
The evaluator is an implementation of a tree-walking interpreter and no 
byte-code is generated. There are no primitive types - everything is an object 
a la Ruby.