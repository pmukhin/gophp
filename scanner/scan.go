package scanner

import (
	"io"
	"unicode"
	"github.com/pmukhin/gophp/token"
)

var tokens = map[string]token.TokenType{
	"return":     token.RETURN,
	"include":    token.INCLUDE,
	"require":    token.REQUIRE,
	"namespace":  token.NAMESPACE,
	"use":        token.USE,
	"final":      token.FINAL,
	"abstract":   token.ABSTRACT,
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
	"each":       token.EACH,
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
	tokenEof     = token.Token{Type: token.EOF, Literal: "EOF"}
	tokenIllegal = token.Token{Type: token.ILLEGAL}

	tokenDoubleDot = token.Token{Type: token.DOUBLE_DOT, Literal: ".."}

	// arithmetic
	tokenPlus        = token.Token{Type: token.PLUS, Literal: "+"}
	tokenMinus       = token.Token{Type: token.MINUS, Literal: "-"}
	tokenMulti       = token.Token{Type: token.MUL, Literal: "*"}
	tokenDiv         = token.Token{Type: token.DIV, Literal: "/"}
	tokenVariable    = token.Token{Type: token.VAR, Literal: "$"}
	tokenCurlyOpen   = token.Token{Type: token.CURLY_OPENING, Literal: "{"}
	tokenCurlyClose  = token.Token{Type: token.CURLY_CLOSING, Literal: "}"}
	tokenSmaller     = token.Token{Type: token.IS_SMALLER, Literal: "<"}
	tokenGreater     = token.Token{Type: token.IS_GREATER, Literal: ">"}
	tokenMultiAssign = token.Token{Type: token.MUL_EQUAL, Literal: "*="}

	tokenSubAssign = token.Token{Type: token.MINUS_EQUAL, Literal: "-="}
	tokenDecrement = token.Token{Type: token.DEC, Literal: "--"}

	tokenMod      = token.Token{Type: token.MOD, Literal: "%"}
	tokenModEqual = token.Token{Type: token.MOD_EQUAL, Literal: "%="}

	tokenDivAssign = token.Token{Type: token.DIV_EQUAL, Literal: "/="}

	tokenAddAssign = token.Token{Type: token.PLUS_EQUAL, Literal: "+="}
	tokenIncrement = token.Token{Type: token.INC, Literal: "++"}

	tokenSmallerOrEqual = token.Token{Type: token.IS_SMALLER_OR_EQUAL, Literal: "<="}
	tokenGreaterOrEqual = token.Token{Type: token.IS_GREATER_OR_EQUAL, Literal: "<="}

	tokenAssign      = token.Token{Type: token.EQUAL, Literal: "="}
	tokenDoubleArrow = token.Token{Type: token.DOUBLE_ARROW, Literal: "=>"}

	tokenComma = token.Token{Type: token.COMMA, Literal: ","}

	tokenSemicolon = token.Token{Type: token.SEMICOLON, Literal: ";"}
	tokenColon     = token.Token{Type: token.COLON, Literal: ":"}

	tokenNot = token.Token{Type: token.NOT, Literal: "!"}

	tokenEqual        = token.Token{Type: token.IS_EQUAL, Literal: "=="}
	tokenIdentical    = token.Token{Type: token.IS_EQUAL, Literal: "==="}
	tokenNotEqual     = token.Token{Type: token.IS_NOT_EQUAL, Literal: "!="}
	tokenNotIdentical = token.Token{Type: token.IS_NOT_IDENTICAL, Literal: "!=="}

	tokenFetch       = token.Token{Type: token.OBJECT_OPERATOR, Literal: "->"}
	tokenStaticFetch = token.Token{Type: token.PAAMAYIM_NEKUDOTAYIM, Literal: "::"}
	tokenBackslash   = token.Token{Type: token.BACKSLASH, Literal: "\\"}

	// paren
	tokenParenOpen  = token.Token{Type: token.PARENTHESIS_OPENING, Literal: "("}
	tokenParenClose = token.Token{Type: token.PARENTHESIS_CLOSING, Literal: ")"}

	// sq brackets
	tokenSquareBracketOpening = token.Token{Type: token.SQUARE_BRACKET_OPENING, Literal: "["}
	tokenSquareBracketClosing = token.Token{Type: token.SQUARE_BRACKET_CLOSING, Literal: "]"}
)

func New(text []rune) *Scanner {
	s := new(Scanner)
	s.len = len(text)
	s.src = text
	s.ch = text[s.offset]

	return s
}

type Scanner struct {
	src        []rune
	len        int
	insertSemi bool
	offset     int
	ch         rune
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
	s.ch = s.src[s.offset]
}

// Next ...
func (s *Scanner) Next() (tok token.Token) {
	s.skipWhitespace()

	insertSemi := false
	switch s.ch {
	case -1:
		tok = tokenEof
	case '\\':
		tok = tokenBackslash
		if s.peek() == '\\' {
			s.next() // eat another one if exists
		}
	case '\'':
		insertSemi = true
		tok = s.readString('\'')
	case '"':
		insertSemi = true
		tok = s.readString('"')
	case ',':
		tok = tokenComma
	case ':':
		if s.peek() == ':' {
			s.next()
			tok = tokenStaticFetch
		} else {
			tok = tokenColon
		}
	case '.':
		if s.peek() == '.' {
			s.next() // eat `.`
			tok = tokenDoubleDot
		} else {
			tok = tokenIllegal
		}
	case ';':
		tok = tokenSemicolon
	case '!':
		if s.peek() == '=' {
			s.next()
			if s.peek() == '=' {
				s.next()
				tok = tokenNotIdentical
			} else {
				tok = tokenNotEqual
			}
		} else {
			tok = tokenNot
		}
	case '=':
		if s.peek() == '>' {
			s.next()
			tok = tokenDoubleArrow
		} else if s.peek() == '=' {
			s.next() // eat `=`
			if s.peek() == '=' {
				s.next()
				tok = tokenIdentical
			} else {
				tok = tokenEqual
			}
		} else {
			tok = tokenAssign
		}
	case '<':
		if s.peek() == '=' {
			s.next()
			tok = tokenSmallerOrEqual
		} else {
			tok = tokenSmaller
		}
	case '>':
		if s.peek() == '=' {
			s.next()
			tok = tokenGreaterOrEqual
		} else {
			tok = tokenGreater
		}
	case '%':
		if s.peek() == '=' {
			s.next()
			tok = tokenModEqual
		} else {
			tok = tokenMod
		}
	case '[':
		tok = tokenSquareBracketOpening
	case ']':
		insertSemi = true
		tok = tokenSquareBracketClosing
	case '(':
		tok = tokenParenOpen
	case ')':
		insertSemi = true
		tok = tokenParenClose
	case '{':
		tok = tokenCurlyOpen
	case '}':
		insertSemi = true
		tok = tokenCurlyClose
	case '$':
		tok = tokenVariable
	case '+':
		if s.peek() == '+' {
			s.next()
			insertSemi = true
			tok = tokenIncrement
		} else if s.peek() == '=' {
			s.next()
			tok = tokenAddAssign
		} else {
			tok = tokenPlus
		}
	case '-':
		next := s.peek()
		if unicode.IsDigit(next) {
			tok = s.scanNumber(true)
		} else if s.peek() == '>' {
			s.next()
			tok = tokenFetch
		} else if s.peek() == '=' {
			s.next()
			tok = tokenSubAssign
		} else if next == '-' {
			s.next()
			insertSemi = true
			tok = tokenDecrement
		} else {
			tok = tokenMinus
		}
	case '\n':
		tok = tokenSemicolon
	case '*':
		if s.peek() == '=' {
			s.next()
			tok = tokenMultiAssign
		} else {
			tok = tokenMulti
		}
	case '/':
		if s.peek() == '*' {
			tok = s.parseCommentMultiLine()
		} else if s.peek() == '/' {
			tok = s.parseLineComment()
		} else if s.peek() == '=' {
			s.next()
			tok = tokenDivAssign
		} else {
			tok = tokenDiv
		}
	default:
		switch {
		case unicode.IsNumber(s.ch):
			insertSemi = true
			tok = s.scanNumber(false)
		case s.isIdentifier(s.ch):
			tok = s.scanIdentifier()
			if tok.Type == token.RETURN || tok.Type == token.IDENT {
				insertSemi = true
			}
		default:
			tok = tokenIllegal
			tok.Literal = string(s.ch)
		}
	}
	s.insertSemi = insertSemi
	tok.Pos = s.offset

	s.next()
	return
}

//
func (s *Scanner) scanNumber(neg bool) token.Token {
	literal := make([]rune, 0, 16)
	if neg {
		literal = append(literal, '-')
		s.next()
	}
	for unicode.IsNumber(s.ch) {
		literal = append(literal, s.ch)
		if !s.HasNext() {
			break
		} else {
			s.next()
		}
	}
	s.backup()
	return token.Token{Type: token.NUMBER, Literal: string(literal)}
}

func (s *Scanner) next() {
	s.offset++
	if s.offset >= s.len {
		s.ch = -1
	} else {
		s.ch = s.src[s.offset]
	}
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' && !s.insertSemi || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) isIdentifier(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r)
}

func (s *Scanner) scanIdentifier() (token.Token) {
	identifier := make([]rune, 0, 32)
	for s.isIdentifier(s.ch) {
		identifier = append(identifier, s.ch)
		if tok, ok := tokens[string(identifier)]; ok {
			return token.Token{Type: tok, Literal: string(identifier)}
		}
		s.next()
	}
	s.backup()
	return token.Token{Type: token.IDENT, Literal: string(identifier)}
}

func (s *Scanner) parseLineComment() token.Token {
	com := make([]rune, 0, 256)
	for {
		if s.ch == '\n' {
			break
		}
		com = append(com, s.ch)
		s.next()
	}
	return token.Token{Type: token.COMMENT, Literal: string(com)}
}

func (s *Scanner) parseCommentMultiLine() token.Token {
	com := make([]rune, 0, 256)
	for {
		com = append(com, s.ch)
		if s.ch == '*' && s.peek() == '/' {
			com = append(com, '/')
			s.next()
			goto exit
		}
		s.next()
	}
exit:
	return token.Token{Type: token.COMMENT, Literal: string(com)}
}

func (s *Scanner) readString(op rune) (tok token.Token) {
	s.next()
	tok.Type = token.STRING
	str := make([]rune, 0, 32)
	for s.ch != op {
		str = append(str, s.ch)
		s.next()
	}
	tok.Literal = string(str)
	return
}
