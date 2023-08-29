module evaluator

replace (
	crabscript.rs/ast => ../ast
	crabscript.rs/lexer => ../lexer
	crabscript.rs/object => ../object
	crabscript.rs/parser => ../parser
	crabscript.rs/token => ../token
)

go 1.21.0

require (
	crabscript.rs/lexer v0.0.0-00010101000000-000000000000
	crabscript.rs/object v0.0.0-00010101000000-000000000000
	crabscript.rs/parser v0.0.0-00010101000000-000000000000
)

require (
	crabscript.rs/ast v0.0.0-00010101000000-000000000000 // indirect
	crabscript.rs/token v0.0.0-00010101000000-000000000000 // indirect
)
