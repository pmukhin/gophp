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
	tokenEof     = token.Token{Type: token.EOF}
	tokenIllegal = token.Token{Type: token.ILLEGAL}

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

func (s *Scanner) Next() (tok token.Token) {
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
	case '\n':
		tok = tokenSemicolon
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
		} else if s.peek() == '=' {
			s.nextChar() // eat `=`
			if s.peek() == '=' {
				s.nextChar()
				tok = tokenIdentical
			} else {
				tok = tokenEqual
			}
		} else {
			tok = tokenAssign
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
func (s *Scanner) scanNumber(neg bool) token.Token {
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
	return token.Token{Type: token.NUMBER, Literal: string(literal)}
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
	return unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r)
}

func (s *Scanner) scanIdentifier() (token.Token) {
	identifier := make([]rune, 0, 32)
	for s.isIdentifier(s.cur) {
		identifier = append(identifier, s.cur)
		if tok, ok := tokens[string(identifier)]; ok {
			return token.Token{Type: tok, Literal: string(identifier)}
		}
		s.nextChar()
	}
	s.backup()
	return token.Token{Type: token.IDENT, Literal: string(identifier)}
}

func (s *Scanner) parseComment() token.Token {
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
	return token.Token{Type: token.COMMENT, Literal: string(com)}
}

func (s *Scanner) readString(op rune) (tok token.Token) {
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
