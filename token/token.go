package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // unknown token
	EOF     = "EOF"     // end of file

	// identifiers and literals
	IDENT = "IDENT" // named var/fns
	INT   = "INT"

	// Ops
	ASSIGN = "="
	PLUS   = "+"

	// Delims
	COMMA     = "," // var delimiter
	SEMICOLON = ";" // line end (along with \n)

	// Scopes
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}