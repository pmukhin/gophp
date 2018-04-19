package ast

type Visitor interface {
	Visit(NodeType, Node)
}
