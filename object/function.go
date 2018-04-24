package object

import (
	"github.com/pmukhin/gophp/ast"
	"fmt"
	"strings"
	"encoding/base64"
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
	representation := fmt.Sprintf("<object of type %s, %s(%s): [NOT IMPLEMENTED]>", this.Class().Name(),
		f.Name(), strings.Join(argsString, ", "))

	return &StringObject{Value: representation}, nil
}

var (
	methods = map[string]Method{
		"__toString": newMethod(funToString, VisibilityPublic),
	}

	functionClass = internalClass{
		name:      "Function",
		final:     true,
		abstract:  false,
		methodSet: newMethodSet(methods),
	}
)

func NewInternalFunc(name string, f func(ctx Context, args ...Object) (Object, error)) FunctionObject {
	return &InternalFunction{name: name, f: f}
}

type FunctionObject interface {
	Object
	Name() string
	Args() []*ast.Arg
	Block() *ast.BlockStatement
}

type InternalFunction struct {
	name string
	args []*ast.Arg
	f    func(ctx Context, args ...Object) (Object, error)
}

func (inf InternalFunction) Class() Class {
	return functionClass
}

func (inf InternalFunction) Id() string { panic("implement me") }

func (inf InternalFunction) Name() string {
	return inf.name
}

func (inf InternalFunction) Args() []*ast.Arg {
	return inf.args
}

func (InternalFunction) Block() *ast.BlockStatement { return nil }

func (inf InternalFunction) Call(ctx Context, args ...Object) (Object, error) {
	return inf.f(ctx, args...)
}

type UserFunction struct {
	name  string
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
		name:  base64.StdEncoding.EncodeToString(b),
		args:  args,
		block: block,
	}
}

func NewUserFunc(ns string, name string, args []*ast.Arg, block *ast.BlockStatement) FunctionObject {
	var uf string

	if ns == "" {
		uf = name
	} else {
		uf = ns + "\\" + name
	}

	return &UserFunction{
		name:  uf,
		args:  args,
		block: block,
	}
}

func (UserFunction) Class() Class { return functionClass }

func (UserFunction) Id() string { panic("implement me") }

func (uf UserFunction) Name() string { return uf.name }

func (uf UserFunction) Args() []*ast.Arg { return uf.args }

func (uf UserFunction) Block() *ast.BlockStatement { return uf.block }
