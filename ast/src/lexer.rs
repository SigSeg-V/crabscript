use crate::token::{Token, TokenKind};

pub struct Lexer {
    input: Box<str>,       // full file taken in
    position: usize,      // current position in the input (points to current char)
    read_position: usize, // current reading position in input (char - 1)
    ch: char,             // current char under examination
}

impl Iterator for Lexer {
    type Item = Token;

    fn next(&mut self) -> Option<Self::Item> {
        match self.next_token() {
            Token{kind: TokenKind::Eof, literal: _} => None,
            tok => Some(tok)
        }
    }
}

impl Lexer {
    // creates new lexer instance
    pub fn new(input: Box<str>) -> Lexer {
        let mut l = Lexer {
            input,
            position: 0,
            read_position: 0,
            ch: char::default(),
        };
        l.read_char();
        l
    }

    // read in the next char in the input
    pub fn read_char(&mut self) {
        // extract value of char or 0 if None
        self.ch = self.peek_char();

        self.position = self.read_position;
        self.read_position += 1;
    }

    pub fn peek_char(&self) -> char {
        // read next char but DO NOT increment position
        self.input.chars().nth(self.read_position).unwrap_or_default()
    }

    // generate new tokenfrom char of input
    pub fn next_token(&mut self) -> Token {

        // go to next potential token
        self.swallow_whitespace();

        let tok = match self.ch {
            '=' => {
                if self.peek_char() == '=' {
                    self.read_char();
                    Token{kind: TokenKind::Eq, literal: "==".into()}
                } else {
                    self.new_token(TokenKind::Bind, self.ch)
                }
            },
            ';' => self.new_token(TokenKind::Semicolon, self.ch),
            '(' => self.new_token(TokenKind::LParen, self.ch),
            ')' => self.new_token(TokenKind::RParen, self.ch),
            '{' => self.new_token(TokenKind::LBrace, self.ch),
            '}' => self.new_token(TokenKind::RBrace, self.ch),
            ',' => self.new_token(TokenKind::Comma, self.ch),
            '+' => self.new_token(TokenKind::Plus, self.ch),
            '-' => self.new_token(TokenKind::Minus, self.ch),
            '/' => self.new_token(TokenKind::Slash, self.ch),
            '*' => self.new_token(TokenKind::Aster, self.ch),
            '!' => {
                if self.peek_char() == '=' {
                    self.read_char();
                    Token{kind: TokenKind::NEq, literal: "!=".into()}
                } else {
                    self.new_token(TokenKind::Bang, self.ch)
                } 
            },
            '<' => self.new_token(TokenKind::Lt, self.ch),
            '>' => self.new_token(TokenKind::Gt, self.ch),
            '\x00' => self.new_token(TokenKind::Eof, self.ch),
            'a'..='z' | 'A'..='Z' => {
                let literal: Box<str> = Box::from(self.read_identifier());
                return Token {
                    kind: Token::lookup_identifier(literal.clone()),
                    literal,
                }
            },
            '0'..='9' => {
                return Token {
                    kind: TokenKind::Integer,
                    literal: self.read_identifier().into(),
                }
            },
            _ => self.new_token(TokenKind::Illegal, self.ch)
        };

        self.read_char();
        tok
    }

    // check valid identifier string
    fn is_letter(&self, ch: char) -> bool {
        ch.is_alphanumeric() || ch == '_'
    }

    // retrieve string of identifier
    fn read_identifier(&mut self) -> &str {
        let position = self.position;
        while self.is_letter(self.ch) {
            self.read_char()
        }
        &self.input[position..self.position]
    }

    // retrieve string of identifier
    fn read_integer(&mut self) -> &str {
        let position = self.position;
        while self.ch.is_ascii_digit() {
            self.read_char()
        }
        &self.input[position..self.position]
    }

    // creating new token from char
    fn new_token(&self, token_kind: TokenKind, ch: char) -> Token {
        Token {
            kind: token_kind,
            literal: Box::from(ch.to_string()),
        }
    }

    fn swallow_whitespace(&mut self) {
        while self.ch.is_whitespace() {
            self.read_char()
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    struct TestCase {
        expected_type: TokenKind,
        expected_literal: Box<str>,
    }

    #[test]
    fn test_next_token() {
        let input = r"
        let five = 5;
        let ten = 10;
        
        let add = fn(x, y) {
            x + y;
        };
        let result = add(five, ten);
        !-/*5;
        5 < 10 > 5;

        if (5 < 10) {
          return true;
        } else {
             return false;
        }

        10 == 10;    
        10 != 1;
        ";

        let tests = vec![
            TestCase {
                expected_type: TokenKind::Let,
                expected_literal: Box::from("let"),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("five"),
            },
            TestCase {
                expected_type: TokenKind::Bind,
                expected_literal: Box::from("="),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("5"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::Let,
                expected_literal: Box::from("let"),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("ten"),
            },
            TestCase {
                expected_type: TokenKind::Bind,
                expected_literal: Box::from("="),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("10"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::Let,
                expected_literal: Box::from("let"),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("add"),
            },
            TestCase {
                expected_type: TokenKind::Bind,
                expected_literal: Box::from("="),
            },
            TestCase {
                expected_type: TokenKind::Function,
                expected_literal: Box::from("fn"),
            },
            TestCase {
                expected_type: TokenKind::LParen,
                expected_literal: Box::from("("),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("x"),
            },
            TestCase {
                expected_type: TokenKind::Comma,
                expected_literal: Box::from(","),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("y"),
            },
            TestCase {
                expected_type: TokenKind::RParen,
                expected_literal: Box::from(")"),
            },
            TestCase {
                expected_type: TokenKind::LBrace,
                expected_literal: Box::from("{"),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("x"),
            },
            TestCase {
                expected_type: TokenKind::Plus,
                expected_literal: Box::from("+"),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("y"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::RBrace,
                expected_literal: Box::from("}"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::Let,
                expected_literal: Box::from("let"),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("result"),
            },
            TestCase {
                expected_type: TokenKind::Bind,
                expected_literal: Box::from("="),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("add"),
            },
            TestCase {
                expected_type: TokenKind::LParen,
                expected_literal: Box::from("("),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("five"),
            },
            TestCase {
                expected_type: TokenKind::Comma,
                expected_literal: Box::from(","),
            },
            TestCase {
                expected_type: TokenKind::Identifier,
                expected_literal: Box::from("ten"),
            },
            TestCase {
                expected_type: TokenKind::RParen,
                expected_literal: Box::from(")"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::Bang,
                expected_literal: Box::from("!"),
            },
            TestCase {
                expected_type: TokenKind::Minus,
                expected_literal: Box::from("-"),
            },
            TestCase {
                expected_type: TokenKind::Slash,
                expected_literal: Box::from("/"),
            },
            TestCase {
                expected_type: TokenKind::Aster,
                expected_literal: Box::from("*"),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("5"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("5"),
            },
            TestCase {
                expected_type: TokenKind::Lt,
                expected_literal: Box::from("<"),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("10"),
            },
            TestCase {
                expected_type: TokenKind::Gt,
                expected_literal: Box::from(">"),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("5"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::If,
                expected_literal: Box::from("if"),
            },
            TestCase {
                expected_type: TokenKind::LParen,
                expected_literal: Box::from("("),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("5"),
            },
            TestCase {
                expected_type: TokenKind::Lt,
                expected_literal: Box::from("<"),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("10"),
            },
            TestCase {
                expected_type: TokenKind::RParen,
                expected_literal: Box::from(")"),
            },
            TestCase {
                expected_type: TokenKind::LBrace,
                expected_literal: Box::from("{"),
            },
            TestCase {
                expected_type: TokenKind::Return,
                expected_literal: Box::from("return"),
            },
            TestCase {
                expected_type: TokenKind::True,
                expected_literal: Box::from("true"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::RBrace,
                expected_literal: Box::from("}"),
            },
            TestCase {
                expected_type: TokenKind::Else,
                expected_literal: Box::from("else"),
            },
            TestCase {
                expected_type: TokenKind::LBrace,
                expected_literal: Box::from("{"),
            },
            TestCase {
                expected_type: TokenKind::Return,
                expected_literal: Box::from("return"),
            },
            TestCase {
                expected_type: TokenKind::False,
                expected_literal: Box::from("false"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::RBrace,
                expected_literal: Box::from("}"),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("10"),
            },
            TestCase {
                expected_type: TokenKind::Eq,
                expected_literal: Box::from("=="),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("10"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("10"),
            },
            TestCase {
                expected_type: TokenKind::NEq,
                expected_literal: Box::from("!="),
            },
            TestCase {
                expected_type: TokenKind::Integer,
                expected_literal: Box::from("1"),
            },
            TestCase {
                expected_type: TokenKind::Semicolon,
                expected_literal: Box::from(";"),
            },
            TestCase {
                expected_type: TokenKind::Eof,
                expected_literal: Box::from("\0"),
            },
        ];

        let mut l = Lexer::new(Box::from(input));
        for t in tests {
            let tok = l.next_token();
            assert_eq!(
                tok.kind, t.expected_type,
                "Token type wrong, got {:?} expected {:?}",
                tok, t.expected_type
            );
            assert_eq!(
                tok.literal, t.expected_literal,
                "Token literal wrong, got {:?} expected {:?}",
                tok.literal, t.expected_literal
            )
        }
    }
}
