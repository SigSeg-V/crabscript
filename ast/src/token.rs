use std::collections::HashMap;
use crate::TokenKind::{Function, Let};

#[derive(Debug, Eq, PartialEq)]
pub enum TokenKind {
    // terminating tokens
    Illegal,
    Eof,

    // identifiers / literals
    Identifier,
    Integer,
    True,
    False,

    // operators
    Bind,
    Plus,
    Minus,
    Bang,
    Aster,
    Slash,
    Lt,
    Gt,
    Eq,
    NEq,

    // delimiters
    Comma,
    Semicolon,
    LParen,
    RParen,
    LBrace,
    RBrace,

    // keywords
    Function,
    Let,
    If,
    Else,
    Return,
}

#[derive(Debug)]
pub struct Token {
    pub kind: TokenKind,
    pub literal: Box<str>,
}

impl Default for Token {
    fn default() -> Token {
        Token {
            kind: TokenKind::Illegal,
            literal: Box::from(""),
        }
    }
}

impl Token {
    pub fn lookup_identifier<S: AsRef<str>>(id: S) -> TokenKind {
        Self::get_keyword(id).unwrap_or(TokenKind::Identifier)
    }

    pub fn get_keyword<S: AsRef<str>>(kw: S) -> Option<TokenKind> {
        match kw.as_ref() {
            "fn" => Some(TokenKind::Function),
            "let" => Some(TokenKind::Let),
            "true" => Some(TokenKind::True),
            "false" => Some(TokenKind::False),
            "return" => Some(TokenKind::Return),
            "if" => Some(TokenKind::If),
            "else" => Some(TokenKind::Else),
            _ => None,
        }
    }
}