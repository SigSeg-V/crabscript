module crabscript.rs/bench

go 1.21.1

replace crabscript.rs/repl => ../repl

replace crabscript.rs/token => ../token

replace crabscript.rs/lexer => ../lexer

replace crabscript.rs/parser => ../parser

replace crabscript.rs/ast => ../ast

replace crabscript.rs/evaluator => ../evaluator

replace crabscript.rs/object => ../object

replace crabscript.rs/compiler => ../compiler

replace crabscript.rs/code => ../code

replace crabscript.rs/vm => ../vm

require (
	crabscript.rs/compiler v0.0.0-00010101000000-000000000000
	crabscript.rs/evaluator v0.0.0-00010101000000-000000000000
	crabscript.rs/lexer v0.0.0-00010101000000-000000000000
	crabscript.rs/object v0.0.0-00010101000000-000000000000
	crabscript.rs/parser v0.0.0-00010101000000-000000000000
	crabscript.rs/vm v0.0.0-00010101000000-000000000000
)

require (
	crabscript.rs/ast v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/code v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/token v0.0.0-00010101000000-000000000000 // indirect
)
