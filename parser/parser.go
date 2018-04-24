package parser

import (
	"github.com/pmukhin/gophp/scanner"
	"github.com/pmukhin/gophp/ast"
	"github.com/pmukhin/gophp/token"
	"fmt"
	"strings"
	"strconv"
	phperror "github.com/pmukhin/gophp/error"
	"reflect"
)

const (
	pLowest     = iota
	pBraces      // {
	pAssignment  // $y = 0
	pSum         // + or -
	pProduct     // *, /, %
	pPrefix      // -$x or --$x
)

var precedences = map[token.TokenType]int{
	token.CURLY_OPEN: pBraces,

	token.EQUAL: pAssignment,

	token.DOUBLE_DOT: pAssignment,

	token.IS_EQUAL:            1,
	token.IS_IDENTICAL:        1,
	token.IS_SMALLER:          1,
	token.IS_SMALLER_OR_EQUAL: 1,
	token.IS_GREATER:          1,
	token.IS_GREATER_OR_EQUAL: 1,


	token.PARENTHESIS_OPENING: 10,
	//
	token.PLUS:  pSum,
	token.MINUS: pSum,

	token.MOD: pProduct,
	token.DIV: pProduct,
	token.MUL: pProduct,

	token.OBJECT_OPERATOR: 5,

	token.INC: pPrefix,
	token.DEC: pPrefix,
}

type (
	prefixParser func() ast.Expression
	infixParser func(ast.Expression) ast.Expression
)

type Parser struct {
	nextToken *token.Token
	curToken  token.Token

	errorFormatter *phperror.Formatter

	prefixExpressionParsers map[token.TokenType]prefixParser
	infixExpressionParsers  map[token.TokenType]infixParser

	err error

	scn *scanner.Scanner
}

func New(s *scanner.Scanner, errorFormatter *phperror.Formatter) *Parser {
	p := new(Parser)
	p.scn = s
	p.errorFormatter = errorFormatter
	p.init()

	return p
}

func (p *Parser) emitError(fmtStr string, args ...interface{}) {
	p.err = p.errorFormatter.Format(fmt.Sprintf(fmtStr, args...), p.curToken.Pos)
}

func (p *Parser) emitErrorInPos(pos int, fmtStr string, args ...interface{}) {
	p.err = p.errorFormatter.Format(fmt.Sprintf(fmtStr, args...), pos)
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
	p.prefixExpressionParsers[token.FOREACH] = p.parseForeach
	p.prefixExpressionParsers[token.PUBLIC] = p.parseMethodDeclaration
	p.prefixExpressionParsers[token.NEW] = p.parseNewExpression

	p.infixExpressionParsers = make(map[token.TokenType]infixParser)
	// infix parsers
	p.infixExpressionParsers[token.EQUAL] = p.parseAssignment
	p.infixExpressionParsers[token.DOUBLE_DOT] = p.parseRangeExpression

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
	p.infixExpressionParsers[token.OBJECT_OPERATOR] = p.parseFetchExpression
	p.infixExpressionParsers[token.PARENTHESIS_OPENING] = p.parseFunctionCall

	p.next()
}

// Parse
func (p *Parser) Parse() (*ast.Module, error) {
	program := &ast.Module{Statements: make([]ast.Statement, 0, 256)}
	for p.curToken.Type != token.EOF {
		st := p.parseStatement()
		if st != nil {
			program.Statements = append(program.Statements, st)
		}
		if p.err != nil {
			return nil, p.err
		}
	}
	return program, nil
}

func (p *Parser) parseForeach() ast.Expression {
	foreach := &ast.ForEachExpression{Token: p.curToken}
	p.next() // eat `foreach`
	endWithParen := false

	if p.oneOf(token.PARENTHESIS_OPENING) {
		endWithParen = true
		p.next() // eat `(`
	}

	// parse array
	foreach.Array = p.parseExpression(pLowest)
	p.eatOfType(token.AS) // eat `as`

	first := p.parseVariable()

	if p.oneOf(token.DOUBLE_ARROW) {
		// we have both keys and values
		p.next() // eat `=>`
		value := p.parseVariable()
		foreach.Key = first.(*ast.VariableExpression)
		foreach.Value = value.(*ast.VariableExpression)
	} else {
		foreach.Value = first.(*ast.VariableExpression)
	}
	if endWithParen {
		p.eatOfType(token.PARENTHESIS_CLOSING)
	}
	foreach.Block = p.parseBlock()

	return foreach
}

// parseFor ... will there be for loop?
func (p *Parser) parseFor() ast.Expression {
	panic("not implemented")
}

func (p *Parser) parseStatement() ast.Statement {
	var st ast.Statement

	switch p.curToken.Type {
	case token.COMMENT:
		p.next() // eat comment
		st = nil
	case token.USE:
		// use NameSpace\\Class;
		st = p.parseUseStatement()
	case token.NAMESPACE:
		// namespace NameSpace;
		st = p.parseNamespaceStatement()
	case token.RETURN:
		// return $value;
		st = p.parseReturnStatement()
	default:
		st = p.parseExpressionStatement()
	}

	if p.oneOf(token.SEMICOLON) {
		p.next() // eat `;` or `\n`
	}
	return st
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	es := new(ast.ExpressionStatement) // alloc
	es.Expression = p.parseExpression(-1)

	if es.Expression == nil {
		return nil
	}

	return es
}

// eatOfType ...
func (p *Parser) eatOfType(tokenType token.TokenType) {
	defer p.next()
	p.assertTokenType(tokenType)
}

func (p *Parser) assertTokenType(tokenType token.TokenType) {
	if p.curToken.Type != tokenType {
		p.emitError("expected token %d, %s given", tokenType, p.curToken.Literal)
	}
}

// parseUseStatement parses statements like
// `use Symfony\Component\HttpFoundation\Response;`
func (p *Parser) parseNamespaceStatement() *ast.NamespaceStatement {
	ns := new(ast.NamespaceStatement)
	namespace := make([]string, 0, 8)
	p.next() // eat `namespace`

	for {
		p.assertTokenType(token.IDENT)
		namespace = append(namespace, p.curToken.Literal)
		p.next() // eat Ident
		if p.curToken.Type == token.BACKSLASH {
			p.next() // eat token.BACKSLASH
			continue
		}
		break
	}

	pathLen := len(namespace)
	if pathLen == 0 {
		p.emitError("empty path in namespace directive")
		return nil
	}
	ns.Namespace = strings.Join(namespace, "\\")

	return ns
}

// parseUseStatement parses statements like
// `use Symfony\Component\HttpFoundation\Response;`
func (p *Parser) parseUseStatement() *ast.UseStatement {
	us := &ast.UseStatement{Token: p.curToken}
	p.next() // eat `use`

	namespace := make([]string, 0, 8)
	for {
		p.assertTokenType(token.IDENT)
		namespace = append(namespace, p.curToken.Literal)
		p.next() // eat Ident
		if p.curToken.Type == token.BACKSLASH {
			p.next() // eat token.BACKSLASH
			continue
		}
		break
	}

	pathLen := len(namespace)
	if pathLen == 0 {
		p.emitError("empty classname in use directive")
	}

	us.Classes = []string{namespace[pathLen-1]}
	us.Namespace = strings.Join(namespace[0:pathLen-1], "\\")

	return us
}

// parseFunctionDeclaration
func (p *Parser) parseFunctionDeclaration() ast.Expression {
	fun := &ast.FunctionDeclarationExpression{Token: p.curToken}
	p.next() // eat `function`

	// function has a name
	if p.curToken.Type == token.IDENT {
		// give function a name
		fun.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.next() // eat IDENT
	} else {
		fun.Anonymous = true
	}

	p.assertTokenType(token.PARENTHESIS_OPENING) // must be `(`
	fun.Args = p.parseArgs()

	if p.curToken.Type == token.USE {
		p.next() // eat `use`
		//fun.BindVars = p.parse
	}

	// have a return type
	if p.curToken.Type == token.COLON {
		fun.ReturnType = p.parseReturnType()
	}

	if p.curToken.Type == token.CURLY_OPENING {
		fun.Block = p.parseBlock()
	} else {
		p.emitError("expected : or {, got %s", p.curToken.Literal)
		return nil
	}

	return fun
}

func (p *Parser) oneOf(ts ...token.TokenType) bool {
	for _, t := range ts {
		if p.curToken.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) parseReturnStatement() ast.Statement {
	r := &ast.ReturnStatement{Token: p.curToken}
	p.next() // eat `return`
	if p.oneOf(token.SEMICOLON, token.CURLY_CLOSING) {
		r.Value = &ast.Null{}
		return r
	}
	r.Value = p.parseExpression(pLowest)
	return r
}

func (p *Parser) parseArgs() []*ast.Arg {
	args := make([]*ast.Arg, 0, 4)
	p.next() // eat `(`
	if p.curToken.Type == token.PARENTHESIS_CLOSING { // if () empty args
		p.next() // eat `)`
		return args
	}

	// if not
	for {
		// okay, we have an arg
		var arg *ast.Arg
		// we have a type!
		if p.curToken.Type == token.IDENT {
			arg = p.parseTypedArg()
		} /* no type, just var def */ else if p.curToken.Type == token.VAR {
			arg = p.parseArg()
		} else {
			p.emitError("unexpected token %s", p.curToken.Literal)
			return nil
		}
		args = append(args, arg)

		if p.curToken.Type == token.PARENTHESIS_CLOSING {
			break
		} else if p.curToken.Type == token.COMMA {
			p.next() // eat `,`
		} else {
			p.emitError("expected , or ), got %s", p.curToken.Literal)
			return nil
		}
	}
	p.next() // eat `)`
	return args
}

// parseReturnType parses return type declarations like
// function()`: ReturnTypeClass` {
func (p *Parser) parseReturnType() *ast.Identifier {
	p.next() // eat `:`
	p.assertTokenType(token.IDENT)
	returnType := p.curToken.Literal
	p.next() // eat IDENT
	// token.IDENT is return type
	return &ast.Identifier{Token: p.curToken, Value: returnType}
}

func (p *Parser) parseBlock() *ast.BlockStatement {
	p.next() // eat `{`
	block := new(ast.BlockStatement)
	for p.curToken.Type != token.CURLY_CLOSING {
		st := p.parseStatement()
		if st == nil || p.err != nil {
			return nil
		}
		block.Statements = append(block.Statements, st)
	}
	p.next() // eat `}`
	return block
}

// parseTypedArg parses typed arg like `array $values = []`
func (p *Parser) parseTypedArg() *ast.Arg {
	arg := new(ast.Arg)
	typeName := p.curToken.Literal // eat `Type`
	p.next()                       // eat ident
	t := &ast.Identifier{Value: typeName, Token: p.curToken}

	arg = p.parseArg()
	arg.Type = t

	return arg
}

// parseArg parses untyped arg like `$value = "someDefaultString"`
func (p *Parser) parseArg() *ast.Arg {
	arg := &ast.Arg{Token: p.curToken}
	p.next() // eat `$`
	// should get $`string_` here
	p.assertTokenType(token.IDENT)
	// assign name
	arg.Name = ast.VariableExpression{Name: p.curToken.Literal, Token: p.curToken}
	p.next() // eat name

	if p.peek().Type == token.EQUAL {
		// we have assign
		p.eatOfType(token.EQUAL)
		arg.DefaultValue = p.parseExpression(pLowest)
	}
	return arg
}

func (p *Parser) getPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return pLowest
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// first get prefix to guess the expression kind
	prefix, ok := p.prefixExpressionParsers[p.curToken.Type]
	// not a prefix
	if !ok {
		p.emitError("unexpected token %s", p.curToken.Literal)
		return nil
	}
	// we got the left
	// e.g. for variable assignment it's $
	left := prefix()
	for precedence < p.getPrecedence() {
		if p.curToken.Type == token.SEMICOLON {
			return left
		}
		infix, ok := p.infixExpressionParsers[p.curToken.Type]
		if !ok {
			return left
		}
		//p.next()
		left = infix(left)
	}
	return left
}

// parseVariable
func (p *Parser) parseVariable() ast.Expression {
	variable := &ast.VariableExpression{Token: p.curToken}
	p.next() // eat $
	p.assertTokenType(token.IDENT)
	variable.Name = p.curToken.Literal
	// eat IDENT
	p.next()
	return variable
}

// parseConstant ...
func (p *Parser) parseConstant() ast.Expression {
	ce := &ast.ConstantExpression{Token: p.curToken}
	p.next() // eat `const`
	p.assertTokenType(token.IDENT)
	ce.Name = p.parseIdentifier().(*ast.Identifier)

	return ce
}

func (p *Parser) parseNewExpression() ast.Expression {
	cle := &ast.NewExpression{Token: p.curToken}
	p.next() // eat `new`

	cle.ClassName = p.parseIdentifier().(*ast.Identifier)
	cle.Args = p.parseExpressionList()

	return cle
}

func (p *Parser) parseMethodDeclaration() ast.Expression {
	mde := &ast.MethodDeclarationExpression{Token: p.curToken}
	mde.Access = modifier(p.curToken.Type)
	p.next() // eat `public|protected|private`

	fun := p.parseFunctionDeclaration()
	mde.FunctionDeclarationExpression = *(fun.(*ast.FunctionDeclarationExpression))

	return mde
}

func (p *Parser) parseClassDeclaration() ast.Expression {
	cde := ast.ClassDeclarationExpression{Token: p.curToken}
	p.next() // eat `class`

	cde.Name = p.parseIdentifier().(*ast.Identifier)
	cde.Block = p.parseBlock()

	return cde
}

func (p *Parser) parseTraitDeclaration() ast.Expression {
	panic("implement me")
}

func (p *Parser) parseAssignment(left ast.Expression) ast.Expression {
	as := &ast.AssignmentExpression{Token: p.curToken}
	p.next() // eat `=`

	switch left.(type) {
	case *ast.VariableExpression:
	case *ast.ConstantExpression:
		// do noting, that's okay
	default:
		p.emitError("can not assign to %s", left.String())
		return nil
	}
	as.Left = left
	as.Right = p.parseExpression(pLowest)

	return as
}

func (p *Parser) parseRangeExpression(left ast.Expression) ast.Expression {
	re := &ast.RangeExpression{}
	p.next() // eat `..`

	re.Left = left
	re.Right = p.parseExpression(pLowest)

	return re
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	be := new(ast.BinaryExpression)
	be.Left = left

	// operator
	be.Op = p.curToken.Literal
	be.Token = p.curToken
	precedence := precedences[p.curToken.Type]
	p.next() // eat operator

	be.Right = p.parseExpression(precedence)

	return be
}

func (p *Parser) parseArrayInitialization() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	p.next() // eat `[`
	if p.oneOf(token.SQUARE_BRACKET_CLOSING) {
		goto exit
	}

	for {
		array.Elements = append(array.Elements, p.parseExpression(pLowest))
		if p.oneOf(token.SQUARE_BRACKET_CLOSING) {
			break
		}
		p.eatOfType(token.COMMA)
	}
exit:
	p.next() // eat `]`
	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	indexExpression := &ast.IndexExpression{Left: left, Token: p.curToken}
	p.next() // eat `[`

	value := p.parseExpression(pLowest)
	p.assertTokenType(token.SQUARE_BRACKET_CLOSING)

	// eat `]`
	p.next()
	indexExpression.Index = value

	return indexExpression
}

func (p *Parser) parseStringLiteral() ast.Expression {
	defer p.next() // eat string
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

func (p *Parser) parseConditionalExpression() ast.Expression {
	ce := &ast.ConditionalExpression{Token: p.curToken}
	p.next() // eat `if`

	ce.Condition = p.parseExpression(pLowest)
	ce.Consequence = p.parseBlock()

	if p.curToken.Type == token.ELSE {
		p.next() // eat `else`
		ce.Alternative = p.parseBlock()
	}

	return ce
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// eat `(
	p.next()
	exp := p.parseExpression(pLowest)
	p.assertTokenType(token.PARENTHESIS_CLOSING)

	// eat `)`
	p.next()
	return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
	i := &ast.Identifier{Token: p.curToken}
	value := p.curToken.Literal
	for {
		p.next() // eat current Ident
		if !p.oneOf(token.BACKSLASH) {
			break
		}
		p.next() // eat `\`
		value += "\\"
		p.assertTokenType(token.IDENT)
		value += p.curToken.Literal
	}
	i.Value = value

	return i
}

func (p *Parser) parseInteger() ast.Expression {
	defer p.next() // eat NUMBER

	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.emitError(err.Error())
		return nil
	}

	return &ast.IntegerLiteral{Value: value}
}

func (p *Parser) parseInstanceOfExpression(left ast.Expression) ast.Expression {
	iof := ast.InstanceOfExpression{Object: left, Token: p.curToken}
	p.next() // eat `instanceof`
	iof.Type = p.parseExpression(pLowest)

	return iof
}

func (p *Parser) parseFunctionCall(left ast.Expression) ast.Expression {
	call := &ast.FunctionCall{Token: p.curToken}
	switch def := left.(type) {
	case *ast.Identifier:
		call.Target = def
	case *ast.FunctionDeclarationExpression:
		call.Target = def
	case *ast.VariableExpression:
		call.Target = def
	default:
		p.emitErrorInPos(left.Pos(), "expected either ident or variable, %v given", def)
	}
	call.CallArgs = p.parseExpressionList()

	return call
}

func (p *Parser) parseExpressionList() []ast.Expression {
	list := make([]ast.Expression, 0, 8)
	p.next() // eat `(`
	if p.curToken.Type == token.PARENTHESIS_CLOSING {
		p.next() // eat `)`
		return list
	}

	for {
		list = append(list, p.parseExpression(pLowest))
		if p.curToken.Type == token.PARENTHESIS_CLOSING {
			p.next() // eat `)`
			break
		} else if p.curToken.Type == token.COMMA {
			// eat ','
			p.next()
		}
	}

	return list
}

func (p *Parser) parseFetchExpression(left ast.Expression) ast.Expression {
	fe := &ast.FetchExpression{Token: p.curToken}
	p.next() // eat `->`

	fe.Left = left
	fe.Right = p.parseExpression(precedences[fe.Token.Type])

	switch fe.Right.(type) {
	case *ast.FunctionCall:
	case *ast.Identifier:
	default:
		p.emitError("unexpected either an Identifier or a FunctionCall, %s given", reflect.TypeOf(fe.Right).String())
		return nil
	}

	return fe
}
