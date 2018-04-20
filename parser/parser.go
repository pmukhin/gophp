package parser

import (
	"github.com/pmukhin/gophp/scanner"
	"github.com/pmukhin/gophp/ast"
	"github.com/pmukhin/gophp/token"
	"fmt"
	"strings"
	"io"
	"strconv"
)

var precedences = map[token.TokenType]int{
	token.CURLY_OPEN: 4,


	token.EQUAL:               1,
	token.IS_EQUAL:            1,
	token.IS_IDENTICAL:        1,
	token.IS_SMALLER:          1,
	token.IS_SMALLER_OR_EQUAL: 1,
	token.IS_GREATER:          1,
	token.IS_GREATER_OR_EQUAL: 1,


	token.PARENTHESIS_OPENING: 1,
	//
	token.PLUS:  2,
	token.MINUS: 2,
	token.MOD:   3,

	// biggest
	token.DIV: 100,
	token.MUL: 100,
}

type (
	prefixParser func() (ast.Expression, error)
	infixParser func(ast.Expression) (ast.Expression, error)
)

type Parser struct {
	nextToken *token.Token
	curToken  token.Token

	prefixExpressionParsers map[token.TokenType]prefixParser
	infixExpressionParsers  map[token.TokenType]infixParser

	scn *scanner.Scanner
}

func New(s *scanner.Scanner) *Parser {
	p := new(Parser)
	p.scn = s
	p.init()

	return p
}

func (p *Parser) next() {
	if p.nextToken == nil {
		tok := p.scn.Next()
		p.curToken = tok
	} else {
		p.curToken = *p.nextToken
		p.nextToken = nil
	}
}

func (p *Parser) peek() token.Token {
	if p.nextToken == nil {
		tok := p.scn.Next()
		p.nextToken = &tok

		return tok
	}
	return *p.nextToken
}

func (p *Parser) init() {
	p.prefixExpressionParsers = make(map[token.TokenType]prefixParser)
	// prefix parsers
	p.prefixExpressionParsers[token.VAR] = p.parseVariable
	p.prefixExpressionParsers[token.CONST] = p.parseConstant
	p.prefixExpressionParsers[token.FUNCTION] = p.parseFunctionDeclaration
	p.prefixExpressionParsers[token.CLASS] = p.parseClassDeclaration
	p.prefixExpressionParsers[token.TRAIT] = p.parseTraitDeclaration
	p.prefixExpressionParsers[token.SQUARE_BRACKET_OPENING] = p.parseArrayInitialization
	p.prefixExpressionParsers[token.STRING] = p.parseStringLiteral
	p.prefixExpressionParsers[token.IF] = p.parseConditionalExpression
	p.prefixExpressionParsers[token.PARENTHESIS_OPENING] = p.parseGroupedExpression
	p.prefixExpressionParsers[token.IDENT] = p.parseIdentifier
	p.prefixExpressionParsers[token.NUMBER] = p.parseInteger
	p.prefixExpressionParsers[token.FOR] = p.parseLoop

	p.infixExpressionParsers = make(map[token.TokenType]infixParser)
	// infix parsers
	p.infixExpressionParsers[token.EQUAL] = p.parseAssignment

	p.infixExpressionParsers[token.PLUS] = p.parseBinaryExpression
	p.infixExpressionParsers[token.MINUS] = p.parseBinaryExpression
	p.infixExpressionParsers[token.MUL] = p.parseBinaryExpression
	p.infixExpressionParsers[token.DIV] = p.parseBinaryExpression
	p.infixExpressionParsers[token.MOD] = p.parseBinaryExpression

	p.infixExpressionParsers[token.IS_GREATER] = p.parseBinaryExpression
	p.infixExpressionParsers[token.IS_GREATER_OR_EQUAL] = p.parseBinaryExpression
	p.infixExpressionParsers[token.IS_SMALLER] = p.parseBinaryExpression
	p.infixExpressionParsers[token.IS_SMALLER_OR_EQUAL] = p.parseBinaryExpression

	p.infixExpressionParsers[token.IS_EQUAL] = p.parseBinaryExpression
	p.infixExpressionParsers[token.IS_NOT_EQUAL] = p.parseBinaryExpression
	p.infixExpressionParsers[token.IS_IDENTICAL] = p.parseBinaryExpression
	p.infixExpressionParsers[token.IS_NOT_IDENTICAL] = p.parseBinaryExpression

	p.infixExpressionParsers[token.SQUARE_BRACKET_OPENING] = p.parseIndexExpression
	p.infixExpressionParsers[token.INSTANCEOF] = p.parseInstanceOfExpression
	p.infixExpressionParsers[token.OBJECT_OPERATOR] = p.parseObjectFetch
	p.infixExpressionParsers[token.PARENTHESIS_OPENING] = p.parseFunctionCall

	p.next()
}

// Parse
func (p *Parser) Parse() (*ast.Module, error) {
	program := &ast.Module{Statements: make([]ast.Statement, 0, 256)}
	for {
		st, err := p.parseStatement()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return program, err
		}
		// end
		program.Statements = append(program.Statements, st)
	}
	return program, nil
}

func (p *Parser) skipComment() {
	for p.curToken.Type == token.COMMENT {
		p.next()
	}
}

func (p *Parser) parseLoop() (ast.Expression, error) {
	p.next() // eat `for`

	if p.curToken.Type == token.EACH {
		foreach := &ast.ForEachExpression{Token: p.curToken}
		p.eatOfType(token.PARENTHESIS_OPENING)
		array, err := p.parseExpression(-1)
		if err != nil {
			return nil, err
		}
		foreach.Array = array

		return foreach, nil
	}

	panic("dsf")
}

func (p *Parser) parseStatement() (st ast.Statement, err error) {
	p.skipComment()
	switch p.curToken.Type {
	case token.EOF:
		err = io.EOF
		return
	case token.USE:
		// use NameSpace\\Class;
		st, err = p.parseUseStatement()
	case token.NAMESPACE:
		// namespace NameSpace;
		st, err = p.parseNamespaceStatement()
	case token.RETURN:
		// return $value;
		st, err = p.parseReturnStatement()
	default:
		st, err = p.parseExpressionStatement()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.next() // eat `;` if it's there
	}

	return
}

func (p *Parser) parseExpressionStatement() (es *ast.ExpressionStatement, err error) {
	es = new(ast.ExpressionStatement) // alloc
	expression, err := p.parseExpression(-1)
	if err != nil {
		return es, err
	}
	es.Expression = expression

	return
}

// eatOfType ...
func (p *Parser) eatOfType(tokenType token.TokenType) error {
	p.next()
	return p.assertTokenType(tokenType)
}

func (p *Parser) assertTokenType(tokenType token.TokenType) error {
	if p.curToken.Type != tokenType {
		return fmt.Errorf("expected token %d, %s given", tokenType, p.curToken.Literal)
	}
	return nil
}

// parseUseStatement parses statements like
// `use Symfony\Component\HttpFoundation\Response;`
func (p *Parser) parseNamespaceStatement() (ns *ast.NamespaceStatement, err error) {
	ns = new(ast.NamespaceStatement)
	namespace := make([]string, 0, 8)
	for {
		p.next()
		if err = p.assertTokenType(token.IDENT); err != nil {
			return
		}
		namespace = append(namespace, p.curToken.Literal)

		p.next()
		if p.curToken.Type == token.SEMICOLON {
			break
		}
		if err = p.assertTokenType(token.BACKSLASH); err != nil {
			return
		}
	}

	pathLen := len(namespace)
	if pathLen == 0 {
		return ns, fmt.Errorf("empty path in namespace directive")
	}
	ns.Namespace = strings.Join(namespace, "\\")
	// eat `;`
	err = p.assertTokenType(token.SEMICOLON)
	p.next()
	return
}

// parseUseStatement parses statements like
// `use Symfony\Component\HttpFoundation\Response;`
func (p *Parser) parseUseStatement() (us *ast.UseStatement, err error) {
	us = &ast.UseStatement{Token: p.curToken}
	namespace := make([]string, 0, 8)

	for {
		p.next()
		if err = p.assertTokenType(token.IDENT); err != nil {
			return
		}
		namespace = append(namespace, p.curToken.Literal)

		p.next()
		if p.curToken.Type == token.SEMICOLON {
			break
		}
		if err = p.assertTokenType(token.BACKSLASH); err != nil {
			return
		}
	}

	pathLen := len(namespace)
	if pathLen == 0 {
		return us, fmt.Errorf("empty classname in use directive")
	}

	us.Classes = []string{namespace[pathLen-1]}
	us.Namespace = strings.Join(namespace[0:pathLen-1], "\\")
	// eat `;`
	err = p.assertTokenType(token.SEMICOLON)
	p.next()
	return
}

// parseFunctionDeclaration
func (p *Parser) parseFunctionDeclaration() (ast.Expression, error) {
	fun := &ast.FunctionDeclarationExpression{Token: p.curToken}
	p.next() // eat `function`

	// function has a name
	if p.curToken.Type == token.IDENT {
		// give function a name
		fun.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		// just for making it explicit
		fun.Anonymous = false
		p.next() // eat IDENT
	}

	e := p.assertTokenType(token.PARENTHESIS_OPENING) // must be `(`
	if e != nil {
		return fun, e
	}
	args, err := p.parseArgs()
	if err != nil {
		return fun, err
	}
	fun.Args = args

	if p.curToken.Type == token.USE {
		p.next() // eat `use`
		//fun.BindVars = p.parse
	}

	// have a return type
	if p.curToken.Type == token.COLON {
		fun.ReturnType, err = p.parseReturnType()
	}

	if p.curToken.Type == token.CURLY_OPENING {
		fun.Block, err = p.parseBlock()
		if err != nil {
			return fun, err
		}
	} else {
		return fun, fmt.Errorf("expected : or {, got %s", p.curToken.Literal)
	}

	return fun, err
}

func (p *Parser) oneOf(ts ...token.TokenType) bool {
	for _, t := range ts {
		if p.curToken.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	r := &ast.ReturnStatement{Token: p.curToken}
	p.next() // eat `return`
	if p.oneOf(token.SEMICOLON, token.CURLY_CLOSING) {
		r.Value = &ast.Null{}
		return r, nil
	}
	val, e := p.parseExpression(-1)
	if e != nil {
		return nil, e
	}
	r.Value = val
	return r, nil
}

func (p *Parser) parseArgs() ([]ast.Arg, error) {
	args := make([]ast.Arg, 0, 4)
	p.next() // eat `(`
	if p.curToken.Type == token.PARENTHESIS_CLOSING { // if () empty args
		p.next() // eat `)`
		return args, nil
	}

	// if not
	for {
		// okay, we have an arg
		var arg ast.Arg
		var err error
		// we have a type!
		if p.curToken.Type == token.IDENT {
			arg, err = p.parseTypedArg()
		} /* no type, just var def */ else if p.curToken.Type == token.VAR {
			arg, err = p.parseArg()
		} else {
			return args, fmt.Errorf("unexpected token %s", p.curToken.Literal)
		}
		if err != nil {
			return args, err
		}
		args = append(args, arg)

		p.next()
		if p.curToken.Type == token.PARENTHESIS_CLOSING {
			break
		} else if p.curToken.Type == token.COMMA {
			p.next() // eat `,`
		} else {
			return args, fmt.Errorf("expected , or ), got %s", p.curToken.Literal)
		}
	}
	p.next() // eat `)`
	return args, nil
}

// parseReturnType parses return type declarations like
// function()`: ReturnTypeClass` {
func (p *Parser) parseReturnType() (*ast.Identifier, error) {
	p.next() // eat `:`
	if err := p.assertTokenType(token.IDENT); err != nil {
		return nil, err
	}
	returnType := p.curToken.Literal
	p.next() // eat IDENT
	// token.IDENT is return type
	return &ast.Identifier{Token: p.curToken, Value: returnType}, nil
}

func (p *Parser) parseBlock() (*ast.BlockStatement, error) {
	p.next() // eat `{`
	block := new(ast.BlockStatement)
	if p.curToken.Type == token.CURLY_CLOSING {
		goto exit
	}
	for {
		st, err := p.parseStatement()
		if err != nil {
			return block, err
		}
		block.Statements = append(block.Statements, st)
		if p.curToken.Type == token.CURLY_CLOSING {
			break
		}
	}
exit:
	p.next() // eat `}`
	return block, nil
}

// parseTypedArg parses typed arg like `array $values = []`
func (p *Parser) parseTypedArg() (arg ast.Arg, err error) {
	typeName := p.curToken.Literal
	if err = p.eatOfType(token.VAR); err != nil {
		return
	}
	if arg, err = p.parseArg(); err != nil {
		return
	}
	arg.Type = &ast.Identifier{Value: typeName, Token: p.curToken}
	return
}

// parseArg parses untyped arg like `$value = "someDefaultString"`
func (p *Parser) parseArg() (ast.Arg, error) {
	arg := ast.Arg{Token: p.curToken}
	p.next() // eat `$`
	// should get $`string_` here
	if err := p.assertTokenType(token.IDENT); err != nil {
		return arg, err
	}
	// assign name
	arg.Name = ast.VariableExpression{Name: p.curToken.Literal, Token: p.curToken}

	if p.peek().Type == token.EQUAL {
		// we have assign
		p.eatOfType(token.EQUAL)
		df, err := p.parseExpression(-1)
		if err != nil {
			return arg, err
		}
		arg.DefaultValue = df
	}

	return arg, nil
}

func (p *Parser) getPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return -999
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	// first get prefix to guess the expression kind
	prefix, ok := p.prefixExpressionParsers[p.curToken.Type]
	// not a prefix
	if !ok {
		return ast.Null{}, fmt.Errorf("%s is not a prefix operator", p.curToken.Literal, p.curToken)
	}
	// we got the left
	// e.g. for variable assignment it's $var
	left, err := prefix()
	if err != nil {
		return left, err
	}
	for precedence < p.getPrecedence() {
		if p.curToken.Type == token.SEMICOLON {
			return left, nil
		}
		infix, ok := p.infixExpressionParsers[p.curToken.Type]
		if !ok {
			return left, nil
		}
		//p.next()
		left, err = infix(left)
		if err != nil {
			return left, err
		}
	}
	return left, nil
}

// parseVariable
func (p *Parser) parseVariable() (ast.Expression, error) {
	variable := &ast.VariableExpression{Token: p.curToken}
	p.next() // eat $
	if err := p.assertTokenType(token.IDENT); err != nil {
		return variable, err
	}
	variable.Name = p.curToken.Literal
	// eat IDENT
	p.next()
	return variable, nil
}

func (p *Parser) parseConstant() (ast.Expression, error) {
	panic("implement me")
}

func (p *Parser) parseClassDeclaration() (ast.Expression, error) {
	panic("implement me")
}

func (p *Parser) parseTraitDeclaration() (ast.Expression, error) {
	panic("implement me")
}

func (p *Parser) parseAssignment(left ast.Expression) (ast.Expression, error) {
	as := &ast.AssignmentExpression{Token: p.curToken}
	p.next() // eat `=`

	switch left.(type) {
	case *ast.VariableExpression:
		// do noting, that's okay
	default:
		return nil, fmt.Errorf("can not assign to %s", left.String())
	}
	as.Left = left
	right, err := p.parseExpression(-1)
	if err != nil {
		return left, err
	}
	as.Right = right

	return as, nil
}

func (p *Parser) parseBinaryExpression(left ast.Expression) (ast.Expression, error) {
	be := new(ast.BinaryExpression)
	be.Left = left

	// operator
	be.Op = p.curToken.Literal
	be.Token = p.curToken
	precedence := precedences[p.curToken.Type]
	p.next() // eat operator

	right, err := p.parseExpression(precedence)

	if err != nil {
		return nil, err
	}

	be.Right = right

	return be, nil
}

func (p *Parser) parseArrayInitialization() (ast.Expression, error) {
	panic("implement me")
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	indexExpression := &ast.IndexExpression{Left: left, Token: p.curToken}
	p.next() // eat `[`

	value, err := p.parseExpression(-1)
	if err != nil {
		return indexExpression, err
	}

	if err := p.assertTokenType(token.SQUARE_BRACKET_CLOSING); err != nil {
		return indexExpression, err
	}
	// eat `]`
	p.next()
	indexExpression.Index = value

	return indexExpression, nil
}

func (p *Parser) parseStringLiteral() (ast.Expression, error) {
	defer p.next() // eat string
	return &ast.StringLiteral{Value: p.curToken.Literal}, nil
}

func (p *Parser) parseConditionalExpression() (ast.Expression, error) {
	ce := &ast.ConditionalExpression{Token: p.curToken}
	p.next() // eat `if`

	conditionExpression, err := p.parseExpression(-1)
	if err != nil {
		return ce, err
	}
	ce.Condition = conditionExpression

	consequenceBlock, err := p.parseBlock()
	if err != nil {
		return ce, err
	}
	ce.Consequence = consequenceBlock
	if p.curToken.Type == token.ELSE {
		p.next() // eat `else`
		alternativeBlock, err := p.parseBlock()
		if err != nil {
			return ce, err
		}
		ce.Alternative = alternativeBlock
	}

	return ce, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	// eat `(
	p.next()
	exp, err := p.parseExpression(-1)
	if err != nil {
		return exp, err
	}
	if err := p.assertTokenType(token.PARENTHESIS_CLOSING); err != nil {
		return exp, err
	}
	// eat `)`
	p.next()
	return exp, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	defer p.next() // eat IDENT
	return &ast.Identifier{Value: p.curToken.Literal}, nil
}

func (p *Parser) parseInteger() (ast.Expression, error) {
	defer p.next() // eat NUMBER
	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	return &ast.IntegerLiteral{Value: value}, err
}

func (p *Parser) parseInstanceOfExpression(left ast.Expression) (ast.Expression, error) {
	iof := ast.InstanceOfExpression{Object: left, Token: p.curToken}
	p.next() // eat `instanceof`

	right, err := p.parseExpression(-1)
	if err != nil {
		return iof, err
	}
	iof.Type = right

	return iof, nil
}

func (p *Parser) parseFunctionCall(left ast.Expression) (ast.Expression, error) {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("expected token ident, %v given", left)
	}
	call := &ast.FunctionCall{Target: ident, Token: p.curToken}
	var err error
	// parse all (<args>)
	call.CallArgs, err = p.parseExpressionList()

	return call, err
}

func (p *Parser) parseExpressionList() ([]ast.Expression, error) {
	list := make([]ast.Expression, 0, 8)
	p.next() // eat `(`
	if p.curToken.Type == token.PARENTHESIS_CLOSING {
		p.next() // eat `)`
		return list, nil
	}

	for {
		arg, err := p.parseExpression(-1)
		if err != nil {
			return list, err
		}
		list = append(list, arg)
		if p.curToken.Type == token.PARENTHESIS_CLOSING {
			p.next() // eat `)`
			break
		} else if p.curToken.Type == token.COMMA {
			// eat ','
			p.next()
		}
	}

	return list, nil
}

func (p *Parser) parseObjectFetch(left ast.Expression) (ast.Expression, error) {
	// eat `->`
	p.next()
	target, err := p.parseIdentifier()
	if err != nil {
		return left, err
	}
	if p.curToken.Type != token.PARENTHESIS_OPENING {
		return ast.PropertyDereference{Object: left, PropertyName: target.String()}, nil
	}
	methodCall := ast.MethodCall{Object: left, FunctionCall: ast.FunctionCall{
		Target: &ast.Identifier{Value: target.String()},
	}}
	// eat `(`
	p.next()
	// if next is `)` there's no arg
	if p.peek().Type == token.PARENTHESIS_CLOSING {
		// eat `)`
		p.next()
		return methodCall, nil
	}

	for {
		arg, err := p.parseExpression(-1)
		if err != nil {
			return methodCall, nil
		}
		methodCall.CallArgs = append(methodCall.CallArgs, arg)
		if p.curToken.Type == token.PARENTHESIS_CLOSING {
			break
		} else if p.curToken.Type == token.COMMA {
			// eat ','
			p.next()
		}
	}

	return methodCall, nil
}
