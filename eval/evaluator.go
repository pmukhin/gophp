package eval

import (
	"github.com/pmukhin/gophp/ast"
	"github.com/pmukhin/gophp/object"
	"errors"
	"fmt"
)

type stateType struct {
	namespace string
}

var state = stateType{}

var opMethods = map[string]string{
	"+": "__add",
	"-": "__sub",
	"/": "__div",
	"*": "__mul",
	"%": "__mod",

	"==":  "__equal",
	"===": "__identical",
	">":   "__gt",
	"<":   "__lt",
	">=":  "__gte",
	"=<":  "__lte",


	"&": "__and",
	"|": "__or",
}

type returnObject struct {
	value object.Object
}

func (returnObject) Class() object.Class { panic("implement me") }

func (returnObject) Id() string { panic("implement me") }

func classExists(class string) bool {
	return false
}

func getClass(class string) object.Class {
	return nil
}

func Eval(node ast.Node, ctx object.Context) (object.Object, error) {
	switch node := node.(type) {
	case *ast.NamespaceStatement:
		state.namespace = node.Namespace
	case *ast.ReturnStatement:
		v, e := Eval(node.Value, ctx)
		if e != nil {
			return object.Null, e
		}
		return returnObject{value: v}, nil
	case *ast.FunctionDeclarationExpression:
		return object.Null, ctx.GetFunctionTable().
			RegisterFunc(object.NewUserFunc(state.namespace, node.Name.Value, node.Args, node.Block))
	case *ast.FunctionCall:
		callName := node.Target.Value
		fun, e := ctx.GetFunctionTable().Find(state.namespace, callName)
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
			funCtx := object.NewContext(ctx, ctx.GetFunctionTable())
			for i, definedArg := range realFun.Args() {
				funCtx.SetContextVar(definedArg.Name.Name, args[i])
			}
			return Eval(realFun.Block(), funCtx)
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
		if node.Alternative != nil {
			return Eval(node.Alternative, ctx)
		}
		return object.Null, err
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
	case *ast.IndexExpression:
		l, e := Eval(node.Left, ctx)
		if e != nil {
			return nil, e
		}
		index, e := Eval(node.Index, ctx)
		if i := l.Class().Methods().Find("__index"); i != nil {
			return i.Call(l, index)
		}
		return object.Null, fmt.Errorf("%v does not support indexing", l.Class().Name())
	case *ast.IntegerLiteral:
		return object.IntegerClass.InternalConstructor(node.Value)
	case *ast.StringLiteral:
		return &object.StringObject{Value: node.Value}, nil
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
		return evalBlock(ctx, node)
	case *ast.Module:
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
