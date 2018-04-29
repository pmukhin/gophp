package parser

import (
	"github.com/pmukhin/gophp/token"
	"github.com/pmukhin/gophp/ast"
)

var modifierTable = map[token.TokenType]int32{
	token.PUBLIC:    ast.ModPublic,
	token.PROTECTED: ast.ModProtected,
	token.PRIVATE:   ast.ModPrivate,
	token.FINAL:     ast.ModFinal,
	token.ABSTRACT:  ast.ModAbstract,
}

func modifier(t token.TokenType) int32 {
	if node, ok := modifierTable[t]; ok {
		return node
	}
	panic("nonexistent tokenType")
}

func isModifier(t token.TokenType) bool {
	_, ok := modifierTable[t]
	return ok
}
