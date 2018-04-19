package ast

import (
	"strings"
	"bytes"
)

// Node represents a node within the AST
//
// All node types implement the Node interface.
type Node interface {
	// TokenLiteral returns the literal of the node
	TokenLiteral() string
	// String returns a string representation of the node
	String() string
}

// A Statement represents a statement within the AST
//
// All statement nodes implement the Statement interface.
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

// literal
type literal interface {
	Node
	literalNode()
}

// IsLiteral returns true if n is a literal node, false otherwise
func IsLiteral(n Node) bool {
	_, ok := n.(literal)
	return ok
}

// UseStatement is a statement like
// `use Symfony\Component\HttpFoundation\Response;`
type UseStatement struct {
	Namespace string
	Classes   []string
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

type Arg struct {
	Type         string
	Name         string
	DefaultValue Expression
}

// FunctionDeclaration ...
type FunctionDeclaration struct {
	Name       string
	Anonymous  bool
	ReturnType string
	Args       []Arg
	Body       BlockStatement
}

func (FunctionDeclaration) TokenLiteral() string {
	return "function"
}

func (fd FunctionDeclaration) String() string {
	out := bytes.Buffer{}
	args := bytes.Buffer{}

	for _, arg := range fd.Args {
		if arg.Type != "" {
			args.WriteString(arg.Type)
			args.WriteString(" ")
		}
		args.WriteString("$" + arg.Name)
		//if arg.DefaultValue != "" {
		//	args.WriteString(" = " + arg.DefaultValue)
		//}
	}

	out.WriteString("function ")
	out.WriteString(fd.Name)
	out.WriteString("(")
	out.WriteString(args.String())
	out.WriteString(")")

	return out.String()
}

// expressionNode ...
func (FunctionDeclaration) expressionNode() {}

type BlockStatement struct {
	Statements []Statement
}

func (BlockStatement) TokenLiteral() string {
	return "{"
}

// String ...
func (BlockStatement) String() string {
	panic("implement me")
}

// statementNode ...
func (BlockStatement) statementNode() {}

// Variable represents single var like
// $someVar
type Variable struct {
	Name string
}

func (Variable) TokenLiteral() string {
	return "$"
}

func (v Variable) String() string {
	return "$" + v.Name
}

// expressionNode ...
func (Variable) expressionNode() {}

// Assignment is an expression like
// $var = $this->callMethod();
type Assignment struct {
	Left  Expression
	Right Expression
}

func (a Assignment) TokenLiteral() string {
	return "="
}

func (a Assignment) String() string {
	return "$" + a.Left.(Variable).Name + " = " + a.Right.String()
}

// expressionNode ...
func (Assignment) expressionNode() {}

// ExpressionStatement is a Statement wrapper an Expression
type ExpressionStatement struct {
	Expression Expression
}

func (ExpressionStatement) TokenLiteral() string {
	panic("implement me")
}

func (es ExpressionStatement) String() string {
	return es.Expression.String() + ";"
}

// statementNode ...
func (ExpressionStatement) statementNode() {}

type Null struct{}

func (Null) TokenLiteral() string {
	return "null"
}

func (Null) String() string {
	return "null"
}

// expressionNode ...
func (Null) expressionNode() {}

type Index struct {
	Left  Expression
	Value Expression
}

func (Index) TokenLiteral() string {
	return "["
}

func (i Index) String() string {
	return i.Left.String() + "[" + i.Value.String() + "]"
}

// expressionNode ...
func (Index) expressionNode() {}

type StringLiteral struct {
	Value string
}

func (s StringLiteral) TokenLiteral() string {
	return "'"
}

func (s StringLiteral) String() string {
	return s.Value
}

// expressionNode
func (StringLiteral) expressionNode() {}

type Condition struct {
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ce Condition) TokenLiteral() string {
	return "if"
}

func (ce Condition) String() string {
	return "if" + ce.Condition.String() + "{" + ce.Consequence.String() + "}"
}

// expressionNode ...
func (Condition) expressionNode() {}

// InstanceOf is an expression like
// $object instanceof SomeType
type InstanceOf struct {
	Object Expression
	Type   Expression
}

func (InstanceOf) TokenLiteral() string {
	return "instanceof"
}

func (iof InstanceOf) String() string {
	return iof.Object.String() + "instanceof" + iof.Type.String()
}

// expressionNode ...
func (InstanceOf) expressionNode() {}

// Identifier
type Identifier struct {
	Value string
}

func (i Identifier) TokenLiteral() string {
	return i.String()
}

func (i Identifier) String() string {
	return i.Value
}

// expressionNode ...
func (Identifier) expressionNode() {}

type PropertyDereference struct {
	Object       Expression
	PropertyName string
}

func (PropertyDereference) TokenLiteral() string {
	return "->"
}

func (pd PropertyDereference) String() string {
	return pd.Object.String() + "->" + pd.PropertyName
}

func (PropertyDereference) expressionNode() {}

type MethodCall struct {
	Object   Expression
	Name     string
	CallArgs []Expression
}

func (MethodCall) TokenLiteral() string {
	panic("implement me")
}

func (MethodCall) String() string {
	panic("implement me")
}

func (MethodCall) expressionNode() {
	panic("implement me")
}

// A Program node is the root node within the AST.
type Program struct {
	Statements []Statement
}
