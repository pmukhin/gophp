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

	"==":  "__equal",
	"===": "__identical",
	">":   "__gt",
	"<":   "__lt",
	">=":  "__gte",
	"=<":  "__lte",


	"&": "__and",
	"|": "__or",
}

func classExists(class string) bool {
	return false
}

func getClass(class string) object.Class {
	return nil
}

func Eval(node ast.Node, ctx object.Context) (object.Object, error) {
	switch node := node.(type) {
	case *ast.FunctionDeclarationExpression:
		o := object.NewUserFunc(node.Name.Value, node.Args, node.Block)
		e := ctx.GetFunctionTable().RegisterFunc(o)
		if e != nil {
			return object.Null, e
		}
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
		switch realFun := fun.(type) {
		case *object.InternalFunction:
			return realFun.Call(ctx, args...)
		case *object.UserFunction:
			return Eval(realFun.Block(), ctx)
		}
	case *ast.ConditionalExpression:
		condition, err := Eval(node.Condition, ctx)
		if err != nil {
			return object.Null, err
		}
		var boolean *object.BooleanObject

		b, isBool := condition.(*object.BooleanObject)
		if !isBool {
			toBoolean := condition.Class().Methods().Find("__toBoolean")
			if toBoolean == nil {
				return object.Null, fmt.Errorf("can not convert %v to Boolean", condition)
			}
			b, e := toBoolean.Call(condition)
			if e != nil {
				return object.Null, err
			}
			boolean = b.(*object.BooleanObject)
		} else {
			boolean = b
		}
		if boolean.Value {
			return Eval(node.Consequence, ctx)
		}
		return Eval(node.Alternative, ctx)
	case *ast.Identifier:
		if node.Value == "true" {
			return object.True, nil
		}
		if node.Value == "false" {
			return object.False, nil
		}
		return object.Null, nil
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
		return nil, fmt.Errorf("operator %s (method %s) is not defined on type %s",
			node.Op, opMethods[node.Op], l.Class().Name())
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
	case *ast.BlockStatement:
		var r object.Object
		for _, st := range node.Statements {
			var e error
			r, e = Eval(st, ctx)
			if e != nil {
				return r, e
			}
		}
		if r == nil {
			return object.Null, nil
		}
		return r, nil
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
