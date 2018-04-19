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

var precedences = map[scanner.TokenType]int{
	token.EQUAL: 0,
}

type (
	prefixParser func() (ast.Expression, error)
	infixParser func(ast.Expression) (ast.Expression, error)
)

type Parser struct {
	nextToken *scanner.Token
	curToken  scanner.Token

	prefixExpressionParsers map[scanner.TokenType]prefixParser
	infixExpressionParsers  map[scanner.TokenType]infixParser

	__scn *scanner.Scanner
}

func New(s *scanner.Scanner) *Parser {
	p := new(Parser)
	p.__scn = s
	p.init()

	return p
}

func (p *Parser) next() {
	if p.nextToken == nil {
		tok := p.__scn.Next()
		p.curToken = tok
	} else {
		p.curToken = *p.nextToken
		p.nextToken = nil
	}
}

func (p *Parser) peek() scanner.Token {
	if p.nextToken == nil {
		tok := p.__scn.Next()
		p.nextToken = &tok

		return tok
	}
	return *p.nextToken
}

func (p *Parser) init() {
	p.prefixExpressionParsers = make(map[scanner.TokenType]prefixParser)
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

	p.infixExpressionParsers = make(map[scanner.TokenType]infixParser)
	// infix parsers
	p.infixExpressionParsers[token.EQUAL] = p.parseAssignment

	p.infixExpressionParsers[token.PLUS] = p.parseBinaryExpression
	p.infixExpressionParsers[token.MINUS] = p.parseBinaryExpression
	p.infixExpressionParsers[token.MUL] = p.parseBinaryExpression
	p.infixExpressionParsers[token.DIV] = p.parseBinaryExpression

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
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{Statements: make([]ast.Statement, 0, 256)}
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

func (p *Parser) parseStatement() (st ast.Statement, err error) {
	switch p.curToken.Type {
	case token.EOF:
		err = io.EOF
		return
		// top
	case token.USE:
		st, err = p.parseUseStatement()
	case token.NAMESPACE:
		st, err = p.parseNamespaceStatement()
	default:
		st, err = p.parseExpressionStatement()
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
	if p.curToken.Type == token.SEMICOLON {
		p.next() // eat `;` if it's there
	}

	return
}

func (p *Parser) eatOfType(tokenType scanner.TokenType) error {
	p.next()
	return p.assertTokenType(tokenType)
}

func (p *Parser) assertTokenType(tokenType scanner.TokenType) error {
	if p.curToken.Type != tokenType {
		return fmt.Errorf("expected token %d, %s given", tokenType, p.curToken.Literal)
	}
	return nil
}

// parseUseStatement parses statements like
// `use Symfony\Component\HttpFoundation\Response;`
func (p *Parser) parseNamespaceStatement() (ns ast.NamespaceStatement, err error) {
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
func (p *Parser) parseUseStatement() (us ast.UseStatement, err error) {
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
	p.next() // eat `function`
	fun := ast.FunctionDeclarationExpression{}

	// function has a name
	if p.curToken.Type == token.IDENT {
		// give function a name
		fun.Name = ast.Identifier{Value: p.curToken.Literal}
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
	return &ast.Identifier{Value: returnType}, nil
}

func (p *Parser) parseBlock() (*ast.BlockStatement, error) {
	p.next() // eat `{`
	block := new(ast.BlockStatement)
	for {
		st, err := p.parseExpressionStatement()
		if err != nil {
			return block, err
		}
		block.Statements = append(block.Statements, st)
		if p.curToken.Type == token.CURLY_CLOSING {
			break
		}
	}
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

	arg.Type = &ast.Identifier{Value: typeName}
	return
}

// parseArg parses untyped arg like `$value = "someDefaultString"`
func (p *Parser) parseArg() (ast.Arg, error) {
	arg := ast.Arg{}
	p.next()
	// should get $`string_` here
	if err := p.assertTokenType(token.IDENT); err != nil {
		return arg, err
	}
	// assign name
	arg.Name = ast.VariableExpression{Name: p.curToken.Literal}

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

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	// first get prefix to guess the expression kind
	prefix, ok := p.prefixExpressionParsers[p.curToken.Type]
	// not a prefix
	if !ok {
		return ast.Null{}, fmt.Errorf("%s is not a prefix operator", p.curToken.Literal)
	}
	// we got the left
	// e.g. for variable assignment it's $var
	left, err := prefix()
	if err != nil {
		return left, err
	}
	for precedence < precedences[p.curToken.Type] {
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
	// eat $
	p.next()
	variable := &ast.VariableExpression{}
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
	// eat `=`
	p.next()

	switch left.(type) {
	case *ast.VariableExpression:
		// do noting, that's okay
	default:
		return ast.Null{}, fmt.Errorf("can not assign to %s", left.String())
	}

	as := &ast.AssignmentExpression{
		Left: left,
	}
	right, err := p.parseExpression(-1)
	if err != nil {
		return left, err
	}
	as.Right = right

	return as, nil
}

func (p *Parser) parseBinaryExpression(left ast.Expression) (ast.Expression, error) {
	// operator
	op := p.curToken.Literal
	p.next() // eat operator
	right, err := p.parseExpression(-1)

	if err != nil {
		return nil, err
	}
	return &ast.BinaryExpression{Left: left, Op: op, Right: right}, nil
}

func (p *Parser) parseArrayInitialization() (ast.Expression, error) {
	panic("implement me")
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	// eat `[`
	p.next()

	indexExpression := ast.IndexExpression{ /*Left: left*/ }
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
	defer p.next()
	return ast.StringLiteral{Value: p.curToken.Literal}, nil
}

func (p *Parser) parseConditionalExpression() (ast.Expression, error) {
	ce := &ast.ConditionalExpression{}
	// eat `if`
	p.next()

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
	// eat `instanceof`
	p.next()
	iof := ast.InstanceOfExpression{Object: left}
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
	call := &ast.FunctionCall{Target: ident}
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
