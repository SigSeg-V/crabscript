package lexer

import (
	"crabscript.rs/token"
	"unicode"
	"unicode/utf8"
)

// A UTF-8 lexer.
// Yes, that means you can write code in emoji, for better or worse.
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

// Gets the next char and increments index
func (l *Lexer) readChar() {
	inpSlice := l.input[l.readPosition:]
	runeChar, runeSize := utf8.DecodeRune([]byte(inpSlice))
	if len(inpSlice) <= 0 {
		l.ch = 0
	} else {
		l.ch = runeChar
	}
	l.position = l.readPosition
	l.readPosition += runeSize
}

// Gets current char being read
func (l *Lexer) peekChar() rune {
	inpSlice := l.input[l.readPosition:]
	runeChar, _ := utf8.DecodeRune([]byte(inpSlice))
	if len(inpSlice) <= 0 {
		return 0
	} else {
		return runeChar
	}
}

// moves to the next token, in cases such as '==' this would move 2 bytes
// instead of 1.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.swallowWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// 2B width token need to read next char
			oldCh := l.ch
			l.readChar()
			tok.Type = token.Eq
			tok.Literal = string(oldCh) + string(l.ch)
		} else {
			tok = newToken(token.Assign, l.ch)
		}
	case ';':
		tok = newToken(token.Semicolon, l.ch)
	case '(':
		tok = newToken(token.LParen, l.ch)
	case ')':
		tok = newToken(token.RParen, l.ch)
	case '{':
		tok = newToken(token.LBrace, l.ch)
	case '}':
		tok = newToken(token.RBrace, l.ch)
	case '+':
		tok = newToken(token.Plus, l.ch)
	case '-':
		tok = newToken(token.Minus, l.ch)
	case '>':
		tok = newToken(token.Gt, l.ch)
	case '<':
		tok = newToken(token.Lt, l.ch)
	case '!':
		if l.peekChar() == '=' {
			// 2B width token need to read next char
			oldCh := l.ch
			l.readChar()
			tok.Type = token.NEq
			tok.Literal = string(oldCh) + string(l.ch)
		} else {
			tok = newToken(token.Bang, l.ch)
		}
	case '*':
		tok = newToken(token.Asterisk, l.ch)
	case '/':
		tok = newToken(token.Slash, l.ch)
	case ',':
		tok = newToken(token.Comma, l.ch)
  case '"': 
    tok.Type = token.String
    tok.Literal = l.readString()
	case 0:
		tok = newToken(token.Eof, l.ch)
	default: // character
		if unicode.IsDigit(l.ch) {
			tok.Type = token.Int
			tok.Literal = l.readNumber()
			return tok
			// TODO assume all non-digit valid chars are usable letters
			// TODO this will allow emojis as bindings
		} else if unicode.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else {
			tok = newToken(token.Illegal, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
  position := l.position + 1

  for {
    l.readChar()
    if l.ch == '"' || l.ch == 0 {
      break
    }
  }

  return l.input[position:l.position]
}

// creates a token from a rune type
func newToken(tokenType token.TokenType, ch rune) token.Token {
	str := ""
	if ch != 0 {
		str = string(ch)
	}
	return token.Token{Type: tokenType, Literal: str}
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

func (l *Lexer) readNumber() string {
	position := l.position
	for unicode.IsNumber(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
