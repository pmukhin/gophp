package object

import (
	"github.com/pmukhin/gophp/ast"
)

var (
	functionClass = internalClass{
		name:     "Function",
		final:    true,
		abstract: false,
	}
)

type Callable interface {
	Call(...Object) (Object, error)
}

type internalFunc struct {
	name string
	f    func(args ...Object) (Object, error)
}

func (ifu internalFunc) Call(o ...Object) (Object, error) {
	return ifu.f(o...)
}

func (internalFunc) Class() Class {
	return functionClass
}

func (internalFunc) Id() string {
	panic("implement me")
}

func NewInternalFunc(name string, f func(args ...Object) (Object, error)) Callable {
	return &internalFunc{name: name, f: f}
}

type FunctionObject struct {
	Name  string
	Args  []ast.Arg
	Block *ast.BlockStatement
}

func (fo FunctionObject) Call(args ...Object) (Object, error) {
	panic("implement me")
}

func (FunctionObject) Class() Class {
	return functionClass
}

func (FunctionObject) Id() string {
	panic("implement me")
}
