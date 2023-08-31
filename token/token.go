package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal = "Illegal" // unknown token
	Eof     = "Eof"     // end of file

	// identifiers and literals
	Ident = "Ident" // named var/fns
	Int   = "Int"
  String = "String"

	// Ops
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"
	Lt       = "<"
	Gt       = ">"
	Eq       = "=="
	NEq      = "!="

	// Delims
	Comma     = "," // var delimiter
	Semicolon = ";" // line end (along with \n)

	// Scopes
	LParen = "("
	RParen = ")"
	LBrace = "{"
	RBrace = "}"

	// Keywords
	Function = "Function"
	Let      = "Let"
	If       = "If"
	Else     = "Else"
	True     = "True"
	False    = "False"
	Return   = "Return"
)

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"if":     If,
	"else":   Else,
	"true":   True,
	"false":  False,
	"return": Return,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
