package parser

import (
	"github.com/pmukhin/gophp/token"
	"github.com/pmukhin/gophp/ast"
)

func modifier(t token.TokenType) int32 {
	switch t {
	case token.PUBLIC:
		return ast.ModPublic
	case token.PROTECTED:
		return ast.ModProtected
	case token.PRIVATE:
		return ast.ModPrivate
	}
	panic("nonexistent tokenType")
}
