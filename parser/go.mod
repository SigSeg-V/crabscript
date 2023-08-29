module crabscript.rs/parser

go 1.20

replace crabscript.rs/token => ../token

replace crabscript.rs/lexer => ../lexer

replace crabscript.rs/ast => ../ast

require (
	crabscript.rs/ast v0.0.0-00010101000000-000000000000
	crabscript.rs/lexer v0.0.0-00010101000000-000000000000
	crabscript.rs/token v0.0.0-00010101000000-000000000000
)
