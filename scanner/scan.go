package scanner

import (
	"io"
	"unicode"
	"lang/token"
)

var tokens = map[string]TokenType{
	"return":     token.RETURN,
	"include":    token.INCLUDE,
	"require":    token.REQUIRE,
	"namespace":  token.NAMESPACE,
	"use":        token.USE,
	"class":      token.CLASS,
	"implements": token.IMPLEMENTS,
	"protected":  token.PROTECTED,
	"public":     token.PUBLIC,
	"private":    token.PRIVATE,
	"function":   token.FUNCTION,
	"as":         token.AS,
	"if":         token.IF,
	"else":       token.ELSE,
	"extends":    token.EXTENDS,
	"for":        token.FOR,
	"foreach":    token.FOREACH,
	"instanceof": token.INSTANCEOF,
	"const":      token.CONST,
	"throw":      token.THROW,
	"new":        token.NEW,
}

var (
	// errors
	eof = io.EOF

	// static tokens
	// system tokens
	tokenEof     = Token{Type: token.EOF}
	tokenIllegal = Token{Type: token.ILLEGAL}

	// arithmetic
	tokenPlus        = Token{Type: token.PLUS, Literal: "+"}
	tokenMinus       = Token{Type: token.MINUS, Literal: "-"}
	tokenMulti       = Token{Type: token.MUL, Literal: "*"}
	tokenDiv         = Token{Type: token.DIV, Literal: "/"}
	tokenVariable    = Token{Type: token.VAR, Literal: "$"}
	tokenCurlyOpen   = Token{Type: token.CURLY_OPENING, Literal: "{"}
	tokenCurlyClose  = Token{Type: token.CURLY_CLOSING, Literal: "}"}
	tokenSmaller     = Token{Type: token.IS_SMALLER, Literal: "<"}
	tokenGreater     = Token{Type: token.IS_GREATER, Literal: ">"}
	tokenMultiAssign = Token{Type: token.MUL_EQUAL, Literal: "*="}

	tokenSubAssign = Token{Type: token.MINUS_EQUAL, Literal: "-="}
	tokenDecrement = Token{Type: token.DEC, Literal: "--"}

	tokenMod      = Token{Type: token.MOD, Literal: "%"}
	tokenModEqual = Token{Type: token.MOD_EQUAL, Literal: "%="}

	tokenDivAssign = Token{Type: token.DIV_EQUAL, Literal: "/="}

	tokenAddAssign = Token{Type: token.PLUS_EQUAL, Literal: "+="}
	tokenIncrement = Token{Type: token.INC, Literal: "++"}

	tokenSmallerOrEqual = Token{Type: token.IS_SMALLER_OR_EQUAL, Literal: "<="}
	tokenGreaterOrEqual = Token{Type: token.IS_GREATER_OR_EQUAL, Literal: "<="}

	tokenEqual       = Token{Type: token.EQUAL, Literal: "="}
	tokenDoubleArrow = Token{Type: token.DOUBLE_ARROW, Literal: "=>"}

	tokenComma = Token{Type: token.COMMA, Literal: ","}

	tokenSemicolon = Token{Type: token.SEMICOLON, Literal: ";"}
	tokenColon     = Token{Type: token.COLON, Literal: ":"}

	tokenNot          = Token{Type: token.NOT, Literal: "!"}
	tokenNotEqual     = Token{Type: token.IS_NOT_EQUAL, Literal: "!="}
	tokenNotIdentical = Token{Type: token.IS_NOT_IDENTICAL, Literal: "!=="}

	tokenFetch       = Token{Type: token.OBJECT_OPERATOR, Literal: "->"}
	tokenStaticFetch = Token{Type: token.PAAMAYIM_NEKUDOTAYIM, Literal: "::"}
	tokenBackslash   = Token{Type: token.BACKSLASH, Literal: "\\"}

	// paren
	tokenParenOpen  = Token{Type: token.PARENTHESIS_OPENING, Literal: "("}
	tokenParenClose = Token{Type: token.PARENTHESIS_CLOSING, Literal: ")"}

	// sq brackets
	tokenSquareBracketOpening = Token{Type: token.SQUARE_BRACKET_OPENING, Literal: "["}
	tokenSquareBracketClosing = Token{Type: token.SQUARE_BRACKET_CLOSING, Literal: "]"}
)

type (
	TokenType uint8
	Token struct {
		Type    TokenType
		Literal string
	}
)

func New(text []rune) *Scanner {
	s := new(Scanner)
	s.len = len(text)
	s.src = text
	s.cur = text[s.offset]

	return s
}

type Scanner struct {
	src    []rune
	len    int
	offset int
	cur    rune
}

// HasNext checks if the string is over
func (s *Scanner) HasNext() bool {
	return s.offset < s.len-1
}

func (s *Scanner) peek() rune {
	return s.src[s.offset+1]
}

func (s *Scanner) backup() {
	s.offset--
	s.cur = s.src[s.offset]
}

func (s *Scanner) Next() (tok Token) {
	s.skipWhitespace()

	switch s.cur {
	case -1:
		tok = tokenEof
	case '\\':
		tok = tokenBackslash
		if s.peek() == '\\' {
			s.nextChar()
		}
	case '\'':
		tok = s.readString('\'')
	case '"':
		tok = s.readString('"')
	case ',':
		tok = tokenComma
	case ':':
		if s.peek() == ':' {
			s.nextChar()
			tok = tokenStaticFetch
		} else {
			tok = tokenColon
		}
	case ';':
		tok = tokenSemicolon
	case '!':
		if s.peek() == '=' {
			s.nextChar()
			if s.peek() == '=' {
				s.nextChar()
				tok = tokenNotIdentical
			} else {
				tok = tokenNotEqual
			}
		} else {
			tok = tokenNot
		}
	case '=':
		if s.peek() == '>' {
			s.nextChar()
			tok = tokenDoubleArrow
		} else {
			tok = tokenEqual
		}
	case '<':
		if s.peek() == '=' {
			s.nextChar()
			tok = tokenSmallerOrEqual
		} else {
			tok = tokenSmaller
		}
	case '>':
		if s.peek() == '=' {
			s.nextChar()
			tok = tokenGreaterOrEqual
		} else {
			tok = tokenGreater
		}
	case '%':
		if s.peek() == '=' {
			s.nextChar()
			tok = tokenModEqual
		} else {
			tok = tokenMod
		}
	case '[':
		tok = tokenSquareBracketOpening
	case ']':
		tok = tokenSquareBracketClosing
	case '(':
		tok = tokenParenOpen
	case ')':
		tok = tokenParenClose
	case '{':
		tok = tokenCurlyOpen
	case '}':
		tok = tokenCurlyClose
	case '$':
		tok = tokenVariable
	case '+':
		if s.peek() == '+' {
			s.nextChar()
			tok = tokenIncrement
		} else if s.peek() == '=' {
			s.nextChar()
			tok = tokenAddAssign
		} else {
			tok = tokenPlus
		}
	case '-':
		next := s.peek()
		if unicode.IsDigit(next) {
			tok = s.scanNumber(true)
		} else if s.peek() == '>' {
			s.nextChar()
			tok = tokenFetch
		} else if s.peek() == '=' {
			s.nextChar()
			tok = tokenSubAssign
		} else if next == '-' {
			s.nextChar()
			tok = tokenDecrement
		} else {
			tok = tokenMinus
		}
	case '*':
		if s.peek() == '=' {
			s.nextChar()
			tok = tokenMultiAssign
		} else {
			tok = tokenMulti
		}
	case '/':
		if s.peek() == '*' {
			tok = s.parseComment()
		} else if s.peek() == '=' {
			s.nextChar()
			tok = tokenDivAssign
		} else {
			tok = tokenDiv
		}
	default:
		switch {
		case unicode.IsNumber(s.cur):
			tok = s.scanNumber(false)
		case s.isIdentifier(s.cur):
			tok = s.scanIdentifier()
		default:
			tok = tokenIllegal
			tok.Literal = string(s.cur)
		}
	}
	s.nextChar()
	return
}

//
func (s *Scanner) scanNumber(neg bool) (Token) {
	literal := make([]rune, 0, 16)
	if neg {
		literal = append(literal, '-')
		s.nextChar()
	}
	for unicode.IsNumber(s.cur) {
		literal = append(literal, s.cur)
		if !s.HasNext() {
			break
		} else {
			s.nextChar()
		}
	}
	s.backup()
	return Token{Type: token.NUMBER, Literal: string(literal)}
}

func (s *Scanner) nextChar() {
	s.offset++
	if s.offset >= s.len {
		s.cur = -1
	} else {
		s.cur = s.src[s.offset]
	}
}

func (s *Scanner) skipWhitespace() {
	for {
		if s.cur == ' ' || s.cur == '\n' || s.cur == '\t' {
			s.nextChar()
		} else {
			break
		}
	}
}

func (s *Scanner) isIdentifier(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (s *Scanner) scanIdentifier() (Token) {
	identifier := make([]rune, 0, 32)
	for s.isIdentifier(s.cur) {
		identifier = append(identifier, s.cur)
		if tok, ok := tokens[string(identifier)]; ok {
			return Token{Type: tok, Literal: string(identifier)}
		}
		s.nextChar()
	}
	s.backup()
	return Token{Type: token.IDENT, Literal: string(identifier)}
}

func (s *Scanner) parseComment() Token {
	com := make([]rune, 0, 256)
	for {
		com = append(com, s.cur)
		if s.cur == '*' && s.peek() == '/' {
			com = append(com, '/')
			s.nextChar()
			goto exit
		}
		s.nextChar()
	}
exit:
	return Token{Type: token.COMMENT, Literal: string(com)}
}

func (s *Scanner) readString(op rune) (tok Token) {
	s.nextChar()
	tok.Type = token.STRING
	str := make([]rune, 0, 32)
	for s.cur != op {
		str = append(str, s.cur)
		s.nextChar()
	}
	tok.Literal = string(str)
	return
}
