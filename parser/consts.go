package parser

import "crabscript.rs/token"

const (
	_ int = iota
	Lowest
	Eq
	Ltgt
	Sum
	Prod
	Prefix
	Call
)

// Precedence of binary operations
var precedences = map[token.TokenType]int{
	token.Eq:       Eq,
	token.NEq:      Eq,
	token.Lt:       Ltgt,
	token.Gt:       Ltgt,
	token.Plus:     Sum,
	token.Minus:    Sum,
	token.Slash:    Prod,
	token.Asterisk: Prod,
	token.LParen:   Call,
}
