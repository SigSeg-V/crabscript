package ast

import (
	"bytes"
	"crabscript.rs/token"
	"strings"
)

type DictLiteral struct {
	Token token.Token // opening dict token {
	Pairs map[Expression]Expression
}

func (dl *DictLiteral) expressionNode() {}

func (dl *DictLiteral) TokenLiteral() string {
	return dl.Token.Literal
}

func (dl *DictLiteral) String() string {
	var out bytes.Buffer

	// stringifying dict pairs
	pairs := []string{}
	for key, value := range dl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
