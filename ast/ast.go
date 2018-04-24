package ast

import (
	"strconv"
	"bytes"
	"strings"
	"github.com/pmukhin/gophp/token"
)

var fourSpaces = strings.Repeat(" ", 4)

type NodeType uint8

const (
	IdentifierType NodeType = iota
	StatementType
	ProgramType
)

const (
	ModPublic    = iota
	ModPrivate
	ModProtected
)

type Node interface {
	// Pos returns the position of first character belonging to the node
	Pos() int
	// End returns the position of first character immediately after the node
	End() int
	// TokenLiteral returns the literal of the node
	TokenLiteral() string
	// String returns a string representation of the node
	String() string
	// Visitor pattern
	Accept(Visitor)
}

type Statement interface {
	Node
	statementNode()
}

// An Expression represents an expression within the AST
//
// All expression nodes implement the Expression interface.
type Expression interface {
	Node
	expressionNode()
}

// Identifier represents variable names, method names, constants
type Identifier struct {
	Token token.Token
	Value string
}

func (Identifier) Pos() int {
	panic("implement me")
}

func (Identifier) End() int {
	panic("implement me")
}

func (Identifier) TokenLiteral() string {
	panic("implement me")
}

func (i Identifier) String() string {
	return i.Value
}

func (i Identifier) Accept(visitor Visitor) {
	visitor.Visit(IdentifierType, i)
}

func (Identifier) expressionNode() {}

// ExpressionStatement is a container for an expression
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es ExpressionStatement) Accept(Visitor) {
	panic("implement me")
}

func (es ExpressionStatement) expressionNode() {
	panic("implement me")
}

func (ExpressionStatement) Pos() int {
	panic("implement me")
}

func (ExpressionStatement) End() int {
	panic("implement me")
}

func (ExpressionStatement) TokenLiteral() string {
	panic("implement me")
}

func (es ExpressionStatement) String() string { return es.Expression.String() + ";" }

func (ExpressionStatement) statementNode() {}

// BlockStatement represents a list of statements
type BlockStatement struct {
	Token token.Token
	// the { token or the first token from the first statement
	Statements []Statement
}

func (BlockStatement) Pos() int {
	panic("implement me")
}

func (BlockStatement) End() int {
	panic("implement me")
}

func (BlockStatement) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (be BlockStatement) String() string {
	out := bytes.Buffer{}
	out.WriteString("{\n")
	for _, k := range be.Statements {
		out.WriteString(fourSpaces + k.String() + "\n")
	}
	out.WriteString("}")

	return out.String()
}

func (BlockStatement) Accept(Visitor) { panic("implement me") }

func (BlockStatement) statementNode() {}

// InstanceOfExpression represents
// $object instanceof SomeType
type InstanceOfExpression struct {
	Token  token.Token
	Object Expression
	Type   Expression
}

func (ioe InstanceOfExpression) Accept(Visitor) {
	panic("implement me")
}

func (InstanceOfExpression) Pos() int {
	panic("implement me")
}

func (InstanceOfExpression) End() int {
	panic("implement me")
}

func (InstanceOfExpression) TokenLiteral() string {
	panic("implement me")
}

func (ioe InstanceOfExpression) String() string {
	return ioe.Object.String() + " instanceof " + ioe.Type.String()
}

// expressionNode ...
func (InstanceOfExpression) expressionNode() {}

// ConditionalExpression represents
// if ($expression) {} else if {} else {}
type ConditionalExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ConditionalExpression) Accept(Visitor) {
	panic("implement me")
}

func (ConditionalExpression) Pos() int {
	panic("implement me")
}

func (ConditionalExpression) End() int {
	panic("implement me")
}

func (ConditionalExpression) TokenLiteral() string { panic("implement me") }

// String ...
func (ce ConditionalExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("if ")
	out.WriteString(ce.Condition.String() + " ")
	out.WriteString(ce.Consequence.String())

	if ce.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ce.Alternative.String())
	}

	return out.String()
}

// expressionNode ...
func (ConditionalExpression) expressionNode() {}

// AssignmentExpression represents
// $var = "someValue"
type AssignmentExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ae AssignmentExpression) Accept(Visitor) {
	panic("implement me")
}

func (AssignmentExpression) Pos() int {
	panic("implement me")
}

func (AssignmentExpression) End() int {
	panic("implement me")
}

func (AssignmentExpression) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (ae AssignmentExpression) String() string { return ae.Left.String() + " = " + ae.Right.String() }

// expressionNode ...
func (AssignmentExpression) expressionNode() {}

// Null represents null
type Null struct{}

func (Null) Accept(Visitor) { panic("implement me") }

func (Null) Pos() int { panic("implement me") }

func (Null) End() int { panic("implement me") }

// TokenLiteral ...
func (Null) TokenLiteral() string { return "null" }

// String
func (Null) String() string { return "null" }

func (Null) expressionNode() {}

// VariableExpression is $variable
type VariableExpression struct {
	Token token.Token
	Name  string
}

func (ve VariableExpression) Accept(Visitor) {
	panic("implement me")
}

func (VariableExpression) Pos() int {
	panic("implement me")
}

func (VariableExpression) End() int {
	panic("implement me")
}

func (VariableExpression) TokenLiteral() string {
	return "$"
}

func (ve VariableExpression) String() string {
	return "$" + ve.Name
}

func (VariableExpression) expressionNode() {}

// IntegerLiteral represents an integer in the AST
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) Accept(Visitor) {
	panic("implement me")
}

func (IntegerLiteral) Pos() int {
	panic("implement me")
}

func (IntegerLiteral) End() int {
	panic("implement me")
}

func (IntegerLiteral) TokenLiteral() string { return "" }

// String ...
func (il IntegerLiteral) String() string { return strconv.FormatInt(il.Value, 10) }

// expressionNode ...
func (IntegerLiteral) expressionNode() {}

// StringLiteral string in the AST
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl StringLiteral) Accept(Visitor) { panic("implement me") }

func (StringLiteral) Pos() int {
	panic("implement me")
}

func (StringLiteral) End() int {
	panic("implement me")
}

// TokenLiteral ...
func (StringLiteral) TokenLiteral() string { return "'" }

// String ...
func (sl StringLiteral) String() string { return "'" + sl.Value + "'" }

// expressionNode
func (StringLiteral) expressionNode() {}

// ArrayLiteral represents expressions like
// ['one', 'two', 3, $obj->getElement()]
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al ArrayLiteral) Accept(Visitor) {
	panic("implement me")
}

func (ArrayLiteral) Pos() int {
	panic("implement me")
}

func (ArrayLiteral) End() int {
	panic("implement me")
}

// TokenLiteral
func (ArrayLiteral) TokenLiteral() string {
	return "["
}

// String ...
func (al ArrayLiteral) String() string {
	out := bytes.Buffer{}
	out.WriteString("[")

	literals := make([]string, len(al.Elements))
	for i, v := range al.Elements {
		literals[i] = v.String()
	}

	out.WriteString(strings.Join(literals, ", "))
	out.WriteString("]")

	return out.String()
}

// expressionNode ...
func (ArrayLiteral) expressionNode() {}

// IndexExpression represents expressions like
// [$someVal] or [0]
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie IndexExpression) Accept(Visitor) {
	panic("implement me")
}

func (IndexExpression) Pos() int {
	panic("implement me")
}

func (IndexExpression) End() int {
	panic("implement me")
}

func (IndexExpression) TokenLiteral() string {
	return "["
}

func (ie IndexExpression) String() string {
	return ie.Left.String() + "[" + ie.Index.String() + "]"
}

// expressionNode ...
func (IndexExpression) expressionNode() {}

// BooleanExpression is either false or true
type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (be BooleanExpression) Accept(Visitor) {
	panic("implement me")
}

func (BooleanExpression) Pos() int {
	panic("implement me")
}

func (BooleanExpression) End() int {
	panic("implement me")
}

func (BooleanExpression) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (be BooleanExpression) String() string {
	if be.Value {
		return "true"
	}

	return "false"
}

// expressionNode ...
func (BooleanExpression) expressionNode() {}

// UnaryExpression is an expression like
// -$var or !$var or ++$var
type UnaryExpression struct {
	Token    token.Token
	IsPrefix bool
	Op       string
	Right    Expression
}

func (ue UnaryExpression) Pos() int { return ue.Token.Pos }

func (UnaryExpression) End() int {
	panic("implement me")
}

func (UnaryExpression) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (ue UnaryExpression) String() string { return ue.Op + ue.Right.String() }

// expressionNode ...
func (UnaryExpression) expressionNode() {}

// BinaryExpression is an expression like
// $var += 1 or $var >= $otherVar or $var != 2
type BinaryExpression struct {
	Token token.Token
	Left  Expression
	Op    string
	Right Expression
}

func (be BinaryExpression) Accept(Visitor) {
	panic("implement me")
}

func (be BinaryExpression) Pos() int { return be.Token.Pos }

func (BinaryExpression) End() int {
	panic("implement me")
}

func (BinaryExpression) TokenLiteral() string {
	panic("implement me")
}

func (be BinaryExpression) String() string {
	return be.Left.String() + " " + be.Op + " " + be.Right.String()
}

// expressionNode
func (BinaryExpression) expressionNode() {}

// WhileExpression represents expressions like
// while (true) { code... }
type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (we WhileExpression) Pos() int { return we.Token.Pos }

func (WhileExpression) End() int {
	panic("implement me")
}

func (WhileExpression) TokenLiteral() string {
	return "while"
}

func (we WhileExpression) String() string {
	return "while (" + we.Body.String() + ")" + we.Body.String()
}

// expressionNode ...
func (WhileExpression) expressionNode() {}

// ForEachExpression represents
// foreach($array as $key => $value) {}
type ForEachExpression struct {
	Token token.Token
	Array Expression
	Key   *VariableExpression
	Value *VariableExpression
	Block *BlockStatement
}

func (fee ForEachExpression) Accept(Visitor) {
	panic("implement me")
}

func (fee ForEachExpression) Pos() int {
	return fee.Token.Pos
}

func (ForEachExpression) End() int {
	panic("implement me")
}

func (ForEachExpression) TokenLiteral() string {
	return "foreach"
}

// String ...
func (fee ForEachExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("foreach(" + fee.Array.String() + " as ")
	if fee.Key != nil {
		out.WriteString(fee.Key.String() + " => ")
	}
	out.WriteString(fee.Value.String() + ") ")
	out.WriteString(fee.Block.String())

	return out.String()
}

// expressionNode ...
func (ForEachExpression) expressionNode() {}

// ForExpression ...
type ForExpression struct {
	Token     token.Token
	Init      Expression
	Condition Expression
	Inc       Expression
}

func (fe ForExpression) Pos() int {
	return fe.Token.Pos
}

func (ForExpression) End() int {
	panic("implement me")
}

func (ForExpression) TokenLiteral() string {
	panic("implement me")
}

func (ForExpression) String() string {
	panic("implement me")
}

func (ForExpression) Accept(Visitor) {
	panic("implement me")
}

// expressionNode ...
func (ForExpression) expressionNode() {}

// RangeExpression ...
type RangeExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (re RangeExpression) Pos() int {
	return re.Token.Pos
}

func (RangeExpression) End() int {
	panic("implement me")
}

func (RangeExpression) TokenLiteral() string {
	panic("implement me")
}

func (re RangeExpression) String() string {
	return re.Left.String() + ".." + re.Right.String()
}

func (RangeExpression) Accept(Visitor) {
	panic("implement me")
}

// expressionNode ...
func (RangeExpression) expressionNode() {}

// FunctionCall represents
// substr($str, 0)
type FunctionCall struct {
	Token    token.Token
	Target   Expression
	CallArgs []Expression
}

func (fc FunctionCall) Accept(Visitor) {
	panic("implement me")
}

func (FunctionCall) Pos() int {
	panic("implement me")
}

func (FunctionCall) End() int {
	panic("implement me")
}

func (FunctionCall) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (fc FunctionCall) String() string {
	out := bytes.Buffer{}
	out.WriteString(fc.Target.String() + "(")

	literals := make([]string, len(fc.CallArgs))
	for i, e := range fc.CallArgs {
		literals[i] = e.String()
	}
	out.WriteString(strings.Join(literals, ", ") + ")")

	return out.String()
}

// expressionNode ...
func (FunctionCall) expressionNode() {}

// MethodCall represents
// $object->method($args)
type MethodCall struct {
	Token  token.Token
	FunctionCall
	Object Expression
}

func (mc MethodCall) Accept(Visitor) {
	panic("implement me")
}

// String ...
func (mc MethodCall) String() string {
	return mc.Object.String() + "->" + mc.FunctionCall.String()
}

// Arg represents function declared argument
type Arg struct {
	Token        token.Token
	Type         *Identifier
	Name         VariableExpression
	DefaultValue Expression
	Variadic     bool
	IsReference  bool
}

// String ...
func (a Arg) String() string {
	out := bytes.Buffer{}
	if a.Type != nil {
		out.WriteString(a.Type.String() + " ")
	}
	if a.IsReference {
		out.WriteString("&")
	}
	if a.Variadic {
		out.WriteString("...")
	}
	out.WriteString(a.Name.String())
	if a.DefaultValue != nil {
		out.WriteString(" = " + a.DefaultValue.String())
	}
	return out.String()
}

// FunctionDeclarationExpression is an expression like
// `function <Name> (<Args> <Variadic>): <ReturnType> { <Block> }`
type FunctionDeclarationExpression struct {
	Token      token.Token
	Anonymous  bool
	Name       *Identifier
	Args       []*Arg
	ReturnType *Identifier
	Block      *BlockStatement
}

func (fde FunctionDeclarationExpression) Pos() int { return fde.Token.Pos }

func (FunctionDeclarationExpression) End() int {
	panic("implement me")
}

func (FunctionDeclarationExpression) TokenLiteral() string {
	panic("implement me")
}

func (FunctionDeclarationExpression) Accept(Visitor) { panic("implement me") }

func (fde FunctionDeclarationExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("function")
	if !fde.Anonymous {
		out.WriteString(" " + fde.Name.String() + "(")
	}

	args := make([]string, len(fde.Args))
	for i, a := range fde.Args {
		args[i] = a.String()
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	if fde.ReturnType != nil {
		out.WriteString(": " + fde.ReturnType.String())
	}

	out.WriteString(" " + fde.Block.String())

	return out.String()
}

func (FunctionDeclarationExpression) expressionNode() {
	panic("implement me")
}

type MethodDeclarationExpression struct {
	Token   token.Token
	Access  int32
	IsFinal bool
	FunctionDeclarationExpression
}

// Module is the whole program for file
type Module struct {
	Token      token.Token
	Statements []Statement
}

// Accept ...
func (p Module) Accept(v Visitor) { v.Visit(ProgramType, p) }

// expressionNode ...
func (p Module) expressionNode() {}

func (Module) Pos() int {
	panic("implement me")
}

func (Module) End() int {
	panic("implement me")
}

func (Module) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (p Module) String() string {
	statements := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		if s != nil {
			statements[i] = s.String()
		}
	}
	return strings.Join(statements, "\n")
}

// statementNode ...
func (Module) statementNode() {
	panic("implement me")
}

// NewExpression is a statement like
// new DateTime('now')
type NewExpression struct {
	Token     token.Token
	ClassName *Identifier
	Args      []Expression
}

func (cis NewExpression) Accept(Visitor) {
	panic("implement me")
}

func (NewExpression) Pos() int {
	panic("implement me")
}

func (NewExpression) End() int {
	panic("implement me")
}

func (NewExpression) TokenLiteral() string {
	panic("implement me")
}

func (cis NewExpression) String() string {
	args := make([]string, len(cis.Args))
	for i, v := range cis.Args {
		args[i] = v.String()
	}
	return "new " + cis.ClassName.String() + "(" + strings.Join(args, ", ") + ")"
}

func (NewExpression) expressionNode() {}

// PropertyDereference
type PropertyDereference struct {
	Token        token.Token
	Object       Expression
	PropertyName string
}

func (pd PropertyDereference) Pos() int {
	panic("implement me")
}

func (pd PropertyDereference) End() int {
	panic("implement me")
}

func (pd PropertyDereference) Accept(Visitor) {
	panic("implement me")
}

func (PropertyDereference) TokenLiteral() string {
	return "->"
}

func (pd PropertyDereference) String() string {
	return pd.Object.String() + "->" + pd.PropertyName
}

func (PropertyDereference) expressionNode() {}

// UseStatement is a statement like
// `use Symfony\Component\HttpFoundation\Response;`
type UseStatement struct {
	Token     token.Token
	Namespace string
	Classes   []string
}

func (us UseStatement) Pos() int {
	panic("implement me")
}

func (us UseStatement) End() int {
	panic("implement me")
}

func (us UseStatement) Accept(Visitor) {
	panic("implement me")
}

// TokenLiteral ...
func (UseStatement) TokenLiteral() string {
	return "use"
}

// String ...
func (us UseStatement) String() string {
	return "use " + us.Namespace + "{" + strings.Join(us.Classes, ",") + "};"
}

// statementNode ...
func (UseStatement) statementNode() {}

// NamespaceStatement is a statement like
// namespace Silex\Application;
type NamespaceStatement struct {
	Token     token.Token
	Namespace string
}

func (ns NamespaceStatement) Pos() int {
	panic("implement me")
}

func (ns NamespaceStatement) End() int {
	panic("implement me")
}

func (ns NamespaceStatement) Accept(Visitor) {
	panic("implement me")
}

// TokenLiteral ...
func (NamespaceStatement) TokenLiteral() string {
	return "namespace"
}

// String ...
func (ns NamespaceStatement) String() string {
	return "namespace " + ns.Namespace + ";"
}

// ConstantExpression ...
type ConstantExpression struct {
	Token token.Token
	Name  *Identifier
}

func (ConstantExpression) Pos() int {
	panic("implement me")
}

func (ConstantExpression) End() int {
	panic("implement me")
}

func (ConstantExpression) TokenLiteral() string {
	panic("implement me")
}

func (ce ConstantExpression) String() string {
	return "const " + ce.Name.String()
}

func (ConstantExpression) Accept(Visitor) {
	panic("implement me")
}

func (ConstantExpression) expressionNode() {
	panic("implement me")
}

// statementNode ...
func (NamespaceStatement) statementNode() {}

type ClassDeclarationExpression struct {
	Token token.Token
	Name  *Identifier
	Block *BlockStatement
}

func (cds ClassDeclarationExpression) Pos() int {
	return cds.Token.Pos
}

func (ClassDeclarationExpression) End() int {
	panic("implement me")
}

func (ClassDeclarationExpression) TokenLiteral() string {
	panic("implement me")
}

func (cde ClassDeclarationExpression) String() string {
	return "class " + cde.Name.String() + " " + cde.Block.String()
}

func (ClassDeclarationExpression) Accept(Visitor) {
	panic("implement me")
}

func (ClassDeclarationExpression) expressionNode() {}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (ReturnStatement) Pos() int {
	panic("implement me")
}

func (ReturnStatement) End() int {
	panic("implement me")
}

func (ReturnStatement) TokenLiteral() string {
	panic("implement me")
}

func (re ReturnStatement) String() string {
	return "return " + re.Value.String()
}

func (ReturnStatement) Accept(Visitor) {
	panic("implement me")
}

func (ReturnStatement) statementNode() {}

// FetchExpression ...
type FetchExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (fe FetchExpression) Pos() int {
	return fe.Token.Pos
}

func (FetchExpression) End() int {
	panic("implement me")
}

func (FetchExpression) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (fe FetchExpression) String() string {
	return fe.Left.String() + "->" + fe.Right.String()
}

func (FetchExpression) Accept(Visitor) {
	panic("implement me")
}

// expressionNode ...
func (FetchExpression) expressionNode() {}
