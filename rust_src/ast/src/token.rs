use std::{collections::HashMap};
use crate::TokenKind::{Function, Let};

#[derive(Debug, Eq, PartialEq)]
pub enum TokenKind {
    // terminating tokens
    Illegal,
    Eof,

    // identifiers / literals
    Identifier,
    Integer,

    // operators
    Bind,
    Plus,
    Minus,
    Bang,
    Aster,
    Slash,
    Lt,
    Gt,

    // delimiters
    Comma,
    Semicolon,
    LParen,
    RParen,
    LBrace,
    RBrace,
    Function,
    Let,
}

#[derive(Debug)]
pub struct Token {
    pub kind: TokenKind,
    pub literal: Box<str>,
    pub keywords: HashMap<Box<str>, TokenKind>, // relating the string to it's enum
}

impl Default for Token {
    fn default() -> Token {
        Token {
            kind: TokenKind::Illegal,
            literal: Box::from(""),
            keywords: HashMap::new(),
        }
    }
}

impl Token {
    pub fn lookup_identifier<S: AsRef<str>>(id: S) -> TokenKind {
        Self::get_keyword(id).unwrap_or(TokenKind::Identifier)
    }

    pub fn get_keyword<S: AsRef<str>>(kw: S) -> Option<TokenKind> {
        match kw.as_ref() {
            "fn" => Some(Function),
            "let" => Some(Let),
            _ => None,
        }
    }
}