package ast

import (
	"bytes"
	"crabscript.rs/token"
	"strings"
)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elem := []string{}
	for _, e := range al.Elements {
		elem = append(elem, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elem, ", "))
	out.WriteString("]")

	return out.String()
}
