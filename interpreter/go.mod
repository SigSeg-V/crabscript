module crabscript.rs/interpreter

go 1.21.0

replace crabscript.rs/repl => ../repl

replace crabscript.rs/token => ../token

replace crabscript.rs/lexer => ../lexer

replace crabscript.rs/parser => ../parser

replace crabscript.rs/ast => ../ast

replace crabscript.rs/evaluator => ../evaluator

replace crabscript.rs/object => ../object

require crabscript.rs/repl v0.0.0-00010101000000-000000000000

require (
	crabscript.rs/ast v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/evaluator v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/lexer v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/object v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/parser v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/token v0.0.0-00010101000000-000000000000 // indirect
)
