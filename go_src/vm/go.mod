module vm

go 1.21.0

replace (
	crabscript.rs/ast => ../ast
	crabscript.rs/code => ../code
	crabscript.rs/compiler => ../compiler
	crabscript.rs/lexer => ../lexer
	crabscript.rs/object => ../object
	crabscript.rs/parser => ../parser
	crabscript.rs/token => ../token
)

require (
	crabscript.rs/ast v0.0.0-00010101000000-000000000000
	crabscript.rs/code v0.0.0-00010101000000-000000000000
	crabscript.rs/compiler v0.0.0-00010101000000-000000000000
	crabscript.rs/lexer v0.0.0-00010101000000-000000000000
	crabscript.rs/object v0.0.0-00010101000000-000000000000
	crabscript.rs/parser v0.0.0-00010101000000-000000000000
)

require crabscript.rs/token v0.0.0-00010101000000-000000000000 // indirect
