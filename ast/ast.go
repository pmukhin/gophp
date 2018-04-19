package ast

import (
	"strconv"
	"bytes"
	"strings"
)

var fourSpaces = strings.Repeat(" ", 4)

type NodeType uint8

const (
	IdentifierType NodeType = iota
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
func (ConditionalExpression) String() string { panic("implement me") }

// expressionNode ...
func (ConditionalExpression) expressionNode() {}

// AssignmentExpression represents
// $var = "someValue"
type AssignmentExpression struct {
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
	Name string
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
	return "[" + ie.Index.String() + "]"
}

// expressionNode ...
func (IndexExpression) expressionNode() {}

// BooleanExpression is either false or true
type BooleanExpression struct {
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
	IsPrefix bool
	Op       string
	Right    Expression
}

func (UnaryExpression) Pos() int {
	panic("implement me")
}

func (UnaryExpression) End() int {
	panic("implement me")
}

func (UnaryExpression) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (pe UnaryExpression) String() string { return pe.Op + pe.Right.String() }

// expressionNode ...
func (UnaryExpression) expressionNode() {}

// BinaryExpression is an expression like
// $var += 1 or $var >= $otherVar or $var != 2
type BinaryExpression struct {
	Left  Expression
	Op    string
	Right Expression
}

func (be BinaryExpression) Accept(Visitor) {
	panic("implement me")
}

func (BinaryExpression) Pos() int {
	panic("implement me")
}

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
	Condition Expression
	Body      *BlockStatement
}

func (WhileExpression) Pos() int {
	panic("implement me")
}

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
	Array Expression
	Key   *VariableExpression
	Value *VariableExpression
	Block *BlockStatement
}

func (ForEachExpression) Pos() int {
	panic("implement me")
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

// FunctionCall represents
// substr($str, 0)
type FunctionCall struct {
	Target   Identifier
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
	Anonymous  bool
	Name       Identifier
	Args       []Arg
	ReturnType *Identifier
	Block      *BlockStatement
}

func (FunctionDeclarationExpression) Pos() int {
	panic("implement me")
}

func (FunctionDeclarationExpression) End() int {
	panic("implement me")
}

func (FunctionDeclarationExpression) TokenLiteral() string {
	panic("implement me")
}

func (FunctionDeclarationExpression) Accept(Visitor) { panic("implement me") }

func (fde FunctionDeclarationExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("function " + fde.Name.String() + "(")

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

// Program is the whole program for file
type Program struct {
	Statements []Statement
}

func (p Program) Accept(Visitor) {
	panic("implement me")
}

func (p Program) expressionNode() {
	panic("implement me")
}

func (Program) Pos() int {
	panic("implement me")
}

func (Program) End() int {
	panic("implement me")
}

func (Program) TokenLiteral() string {
	panic("implement me")
}

// String ...
func (p Program) String() string {
	statements := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		if s != nil {
			statements[i] = s.String()
		}
	}
	return strings.Join(statements, "\n")
}

// statementNode ...
func (Program) statementNode() {
	panic("implement me")
}

// ClassInstantiationExpression is a statement like
// new DateTime('now')
type ClassInstantiationExpression struct {
	ClassName Identifier
	Args      []Expression
}

func (cis ClassInstantiationExpression) Accept(Visitor) {
	panic("implement me")
}

func (ClassInstantiationExpression) Pos() int {
	panic("implement me")
}

func (ClassInstantiationExpression) End() int {
	panic("implement me")
}

func (ClassInstantiationExpression) TokenLiteral() string {
	panic("implement me")
}

func (cis ClassInstantiationExpression) String() string {
	args := make([]string, len(cis.Args))
	for i, v := range cis.Args {
		args[i] = v.String()
	}
	return "new " + cis.ClassName.String() + "(" + strings.Join(args, ", ") + ")"
}

func (ClassInstantiationExpression) expressionNode() {}

// PropertyDereference
type PropertyDereference struct {
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

// statementNode ...
func (NamespaceStatement) statementNode() {}
