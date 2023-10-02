module crabscript.rs/object

replace (
	crabscript.rs/ast => ../ast
	crabscript.rs/code => ../code
	crabscript.rs/token => ../token
)

go 1.21.0

require (
	crabscript.rs/ast v0.0.0-00010101000000-000000000000
	crabscript.rs/code v0.0.0-00010101000000-000000000000
)

require crabscript.rs/token v0.0.0-00010101000000-000000000000 // indirect
