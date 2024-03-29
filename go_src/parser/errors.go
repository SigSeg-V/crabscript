package parser

import (
	"crabscript.rs/token"
	"fmt"
)

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token %v, got %v", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse fn available for %v", t)
	p.errors = append(p.errors, msg)
}
