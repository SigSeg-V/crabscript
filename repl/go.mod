module crabscript.rs/repl

go 1.20

replace crabscript.rs/token => ../token

replace crabscript.rs/lexer => ../lexer

require (
	crabscript.rs/lexer v0.0.0-00010101000000-000000000000
	crabscript.rs/token v0.0.0-00010101000000-000000000000
)
