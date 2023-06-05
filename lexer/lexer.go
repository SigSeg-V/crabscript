package lexer

import (
	"crabscript.rs/token"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input        string
	position     int  // current index into input
	readPosition int  //current reading pos in input (position + 1)
	ch           rune // current char
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	inpSlice := l.input[l.readPosition:]
	runeChar, runeSize := utf8.DecodeRune([]byte(inpSlice))
	if len(l.input) <= 0 {
		l.ch = 0
	} else {
		l.ch = runeChar
	}
	l.position = l.readPosition
	l.readPosition += runeSize
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	
	l.swallowWhitespace()
	
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case 0:
		tok = newToken(token.EOF, l.ch)
	default: // character
		if unicode.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.ch){
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for unicode.IsLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) swallowWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

