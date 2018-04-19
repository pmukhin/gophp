package eval

import (
	"github.com/pmukhin/gophp/ast"
	"github.com/pmukhin/gophp/object"
	"errors"
	"fmt"
)

var opMethods = map[string]string{
	"+": "__add",
	"-": "__sub",
	"/": "__div",
	"*": "__mul",

	"==": "__equal",
	">":  "__gt",
	"<":  "__lt",
	">=": "__gte",
	"=<": "__lte",


	"&": "__and",
	"|": "__or",
}

func classExists(class string) bool {
	return false
}

func getClass(class string) object.Class {
	return nil
}

func Eval(node ast.Node, ctx Context) (object.Object, error) {
	switch node := node.(type) {
	case *ast.FunctionCall:
		name := node.Target.Value
		fun, e := ctx.GetFunctionTable().Find(name)
		if e != nil {
			return nil, e
		}
		args := make([]object.Object, len(node.CallArgs))
		for i, a := range node.CallArgs {
			args[i], e = Eval(a, ctx)
			if e != nil {
				return nil, e
			}
		}
		return fun.Call(args...)

	case *ast.BinaryExpression:
		l, e := Eval(node.Left, ctx)
		if e != nil {
			return nil, e
		}
		r, e := Eval(node.Right, ctx)
		if e != nil {
			return nil, e
		}
		if m := l.Class().Methods().Find(opMethods[node.Op]); m != nil {
			return m.Call(l, r)
		}
		return nil, fmt.Errorf("undefined operator %s", node.Op)
	case *ast.IntegerLiteral:
		return object.IntegerClass.InternalConstructor(node.Value)
	case *ast.VariableExpression:
		v, e := ctx.GetContextVar(node.Name)
		if e != nil {
			return object.Null, nil
		}
		return v, nil
	case *ast.Null:
		return object.Null, nil
	case *ast.AssignmentExpression:
		right, e := Eval(node.Right, ctx)
		if e != nil {
			return nil, e
		}
		switch left := node.Left.(type) {
		case *ast.VariableExpression:
			ctx.SetContextVar(left.Name, right)
			return right, nil
		}
	case *ast.ClassInstantiationExpression:
		name := node.ClassName.String()
		if !classExists(name) {
			return nil, errors.New("class " + name + " does not exist")
		}
		args := make([]object.Object, len(node.Args))
		for i, argExpression := range node.Args {
			a, e := Eval(argExpression, ctx)
			if e != nil {
				return nil, e
			}
			args[i] = a
		}
		cls := getClass(name)

		return cls.Constructor().Call(nil, args...)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, ctx)
	case *ast.Program:
		var (
			r object.Object
			e error
		)
		for _, st := range node.Statements {
			r, e = Eval(st, ctx)
			if e != nil {
				return nil, e
			}
		}
		return r, nil
	default:
		return nil, fmt.Errorf("unexpected node %s", node.String())
	}
	return object.Null, nil
}
