package object

import (
	"github.com/pmukhin/gophp/ast"
	"fmt"
	"strings"
)

func funToString(this Object, args ...Object) (Object, error) {
	if len(args) != 0 {
		return Null, fmt.Errorf("expected 0 arguments, %d given", len(args))
	}
	f := this.(FunctionObject)
	argsString := make([]string, len(f.Args()))
	for i, arg := range f.Args() {
		str := fmt.Sprintf(arg.String())
		argsString[i] = str
	}
	representation := fmt.Sprintf("<object of type %s, (%s): [NOT IMPLEMENTED]>", this.Class().Name(),
		strings.Join(argsString, ", "))

	return &StringObject{Value: representation}, nil
}

var (
	methods = map[string]Method{
		"__toString": newMethod(funToString, VisibilityPublic),
	}

	functionClass = &InternalClass{
		name:      "Function",
		final:     true,
		abstract:  false,
		methodSet: newMethodSet(methods),
	}
)

func NewInternalFunc(f func(args ...Object) (Object, error)) FunctionObject {
	return &InternalFunction{f: f}
}

type FunctionObject interface {
	Object
	Args() []*ast.Arg
	Block() *ast.BlockStatement
}

type InternalFunction struct {
	args []*ast.Arg
	f    func(args ...Object) (Object, error)
}

func (inf InternalFunction) Class() Class {
	return functionClass
}

func (inf InternalFunction) Id() string { panic("implement me") }

func (inf InternalFunction) Args() []*ast.Arg {
	return inf.args
}

func (InternalFunction) Block() *ast.BlockStatement { return nil }

func (inf InternalFunction) Call(args ...Object) (Object, error) {
	return inf.f(args...)
}

type UserFunction struct {
	args  []*ast.Arg
	block *ast.BlockStatement
}

// NewAnonymousFunc ...
func NewAnonymousFunc(args []*ast.Arg, block *ast.BlockStatement) FunctionObject {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(i<<2*31 + i)
	}
	return &UserFunction{
		args:  args,
		block: block,
	}
}

func NewUserFunc(args []*ast.Arg, block *ast.BlockStatement) FunctionObject {
	return &UserFunction{
		args:  args,
		block: block,
	}
}

func (UserFunction) Class() Class { return functionClass }

func (UserFunction) Id() string { panic("implement me") }

func (uf UserFunction) Args() []*ast.Arg { return uf.args }

func (uf UserFunction) Block() *ast.BlockStatement { return uf.block }
