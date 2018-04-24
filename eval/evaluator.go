package eval

import (
	"errors"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"github.com/pmukhin/gophp/ast"
	"github.com/pmukhin/gophp/object"
	"strings"
)

type exceptionError struct {
	class   object.Object
	message string
}

func (e exceptionError) Error() string {
	return e.message
}

type stateType struct {
	namespace    string
	namespaceSet bool
	uses         map[string]string
}

func (s *stateType) SetNamespace(ns string) error {
	if s.namespaceSet {
		return errors.New("namespace is already set")
	}
	s.namespace = ns
	s.namespaceSet = true

	return nil
}

func newState() *stateType {
	s := new(stateType)
	s.uses = make(map[string]string)

	return s
}

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

// Evaluator ...
type Evaluator interface {
	Eval(ast.Node, object.Context) (object.Object, error)
}

func New() Evaluator {
	ev := new(evaluator)
	ev.stack = stack.New()
	ev.state = newState()

	return ev
}

type evaluator struct {
	stack *stack.Stack
	state *stateType
}

//
func (ev *evaluator) evalBlock(ctx object.Context, node *ast.BlockStatement) (object.Object, error) {
	var ret object.Object = nil

	for _, st := range node.Statements {
		r, err := ev.Eval(st, ctx)
		if err != nil {
			return returnObject{value: object.Null}, err
		}
		// if it's return
		if ret, ok := r.(returnObject); ok {
			return ret, nil
		}
		// add
		ret = r
	}

	if ret == nil {
		ret = returnObject{value: object.Null}
	}

	return ret, nil
}

func (ev *evaluator) evalRegisteredFunc(name *ast.Identifier, callArgs []ast.Expression, ctx object.Context) (object.Object, error) {
	callName := name.Value
	if use, ok := ev.state.uses[callName]; ok {
		callName = use
	}

	fun, err := ctx.GetGlobal(callName)
	if err != nil {
		return nil, err
	}
	args := make([]object.Object, len(callArgs))
	for i, a := range callArgs {
		args[i], err = ev.Eval(a, ctx)
		if err != nil {
			return nil, err
		}
	}
	switch realFun := fun.(type) {
	case *object.InternalFunction:
		return realFun.Call(args...)
	case *object.UserFunction:
		funCtx, e := ev.injectArgs(ctx, callArgs, realFun)
		if e != nil {
			return object.Null, e
		}
		return ev.Eval(realFun.Block(), funCtx)
	default:
		panic("unexpected function type")
	}
}

func (ev *evaluator) injectArgs(ctx object.Context, callArgs []ast.Expression, fun object.FunctionObject) (object.Context, error) {
	var err error
	args := make([]object.Object, len(callArgs))
	for i, a := range callArgs {
		args[i], err = ev.Eval(a, ctx)
		if err != nil {
			return nil, err
		}
	}
	funCtx := object.CloneContext(ctx, nil)
	for i, definedArg := range fun.Args() {
		funCtx.SetContextVar(definedArg.Name.Name, args[i])
	}

	return funCtx, nil
}

// unpackReturnObject ...
func unpackReturnObject(o object.Object, err error) (object.Object, error) {
	if v, ok := o.(returnObject); ok {
		return v.value, err
	}
	return o, err
}

// evalForeach ...
func (ev *evaluator) evalForeach(foreach *ast.ForEachExpression, ctx object.Context) (object.Object, error) {
	array, err := ev.Eval(foreach.Array, ctx)
	if err != nil {
		return object.Null, err
	}
	for index, value := range array.(*object.ArrayObject).Values {
		if foreach.Key != nil {
			ctx.SetContextVar(foreach.Key.Name, &object.IntegerObject{Value: int64(index)})
		}
		ctx.SetContextVar(foreach.Value.Name, value)
		_, err := ev.Eval(foreach.Block, ctx)

		if err != nil {
			return object.Null, err
		}
	}

	return object.Null, nil
}

// registerFunc puts func into globals table
func registerFunc(ctx object.Context, name string, fun object.FunctionObject) error {
	return ctx.SetGlobal(name, fun)
}

func (ev *evaluator) wrap(err error, node ast.Node) error {
	stackTrace := make([]string, ev.stack.Len())
	for i := 0; i < ev.stack.Len(); i++ {
		stackTrace[i] = ev.stack.Pop().(string)
	}
	return fmt.Errorf("%s at %s", err.Error(), strings.Join(stackTrace, "\n"))
}

func (ev *evaluator) Eval(node ast.Node, ctx object.Context) (object.Object, error) {
	return ev.doEval(node, ctx)
}

func (ev *evaluator) evalFunctionCall(node *ast.FunctionCall, ctx object.Context) (object.Object, error) {
	ev.stack.Push(fmt.Sprintf("%s", node.Target.String()))
	if funcName, ok := node.Target.(*ast.Identifier); ok {
		return unpackReturnObject(ev.evalRegisteredFunc(funcName, node.CallArgs, ctx))
	}
	anonymous := node.Target.(ast.Expression)
	resolve, err := ev.Eval(anonymous, ctx)
	if err != nil {
		return object.Null, err
	}
	funCtx, err := ev.injectArgs(ctx, node.CallArgs, resolve.(object.FunctionObject))
	return unpackReturnObject(ev.Eval(resolve.(object.FunctionObject).Block(), funCtx))
}

func (ev *evaluator) evalArray(node *ast.ArrayLiteral, ctx object.Context) (object.Object, error) {
	values := make([]object.Object, len(node.Elements))
	for i, expression := range node.Elements {
		result, err := ev.Eval(expression, ctx)
		if err != nil {
			return object.Null, err
		}
		values[i] = result
	}
	return object.NewArray(values...)
}

func (ev *evaluator) doEval(node ast.Node, ctx object.Context) (object.Object, error) {
	switch node := node.(type) {
	case *ast.RangeExpression:
		return ev.unpackRange(node, ctx)
	case *ast.ArrayLiteral:
		return ev.evalArray(node, ctx)
	case *ast.NamespaceStatement:
		ev.state.SetNamespace(node.Namespace)
	case *ast.UseStatement:
		for _, name := range node.Classes {
			ev.state.uses[name] = node.Namespace + "\\" + name
		}
	case *ast.ForEachExpression:
		return ev.evalForeach(node, ctx)
	case *ast.ReturnStatement:
		v, e := ev.Eval(node.Value, ctx)
		if e != nil {
			return object.Null, e
		}
		return returnObject{value: v}, nil
	case *ast.FunctionDeclarationExpression:
		if node.Anonymous == true {
			return object.NewAnonymousFunc(node.Args, node.Block), nil
		}
		name := object.FullyQ(ev.state.namespace, node.Name.Value)
		return object.Null, registerFunc(ctx, name, object.NewUserFunc(node.Args, node.Block))
	case *ast.FunctionCall:
		return ev.evalFunctionCall(node, ctx)
	case *ast.FetchExpression:
		return ev.evalFetchExpression(node, ctx)
	case *ast.ConditionalExpression:
		condition, err := ev.Eval(node.Condition, ctx)
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
			return ev.Eval(node.Consequence, ctx)
		}
		if node.Alternative != nil {
			return ev.Eval(node.Alternative, ctx)
		}
		return object.Null, err
	case *ast.Identifier:
		if node.Value == "true" {
			return object.True, nil
		}
		if node.Value == "false" {
			return object.False, nil
		}
		// other constant?
		return ctx.GetGlobal(node.Value)
	case *ast.BinaryExpression:
		l, err := ev.Eval(node.Left, ctx)
		if err != nil {
			return nil, err
		}
		r, err := ev.Eval(node.Right, ctx)
		if err != nil {
			return nil, err
		}
		if m := l.Class().Methods().Find(opMethods[node.Op]); m != nil {
			return m.Call(l, r)
		}
		return nil, fmt.Errorf("operator %s (method %s) is not defined on type %s",
			node.Op, opMethods[node.Op], l.Class().Name())
	case *ast.IndexExpression:
		l, err := ev.Eval(node.Left, ctx)
		if err != nil {
			return nil, err
		}
		index, err := ev.Eval(node.Index, ctx)
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
		right, e := ev.Eval(node.Right, ctx)
		if e != nil {
			return nil, e
		}
		switch left := node.Left.(type) {
		case *ast.ConstantExpression:
			e := ctx.SetGlobal(left.Name.Value, right)
			return object.Null, e
		case *ast.VariableExpression:
			e := ctx.SetContextVar(left.Name, right)
			return right, e
		}
	case *ast.ClassInstantiationExpression:
		panic("not implemented")
	case *ast.ExpressionStatement:
		return ev.Eval(node.Expression, ctx)
	case *ast.BlockStatement:
		return ev.evalBlock(ctx, node)
	case *ast.Module:
		var (
			r   object.Object
			err error
		)
		for _, st := range node.Statements {
			r, err = ev.Eval(st, ctx)
			if err != nil {
				return nil, err
			}
		}
		return r, nil
	default:
		return nil, fmt.Errorf("unexpected node %s", node.String())
	}
	return object.Null, nil
}

func (ev *evaluator) unpackRange(ex *ast.RangeExpression, ctx object.Context) (object.Object, error) {
	valLeft, err := ev.Eval(ex.Left, ctx)
	if err != nil {
		return object.Null, err
	}
	valRight, err := ev.Eval(ex.Right, ctx)
	if err != nil {
		return object.Null, err
	}

	intLeft, ok := valLeft.(*object.IntegerObject)
	if !ok {
		return object.Null, fmt.Errorf("can not iterate over %s", intLeft.Class().Name())
	}
	intRight, ok := valRight.(*object.IntegerObject)
	if !ok {
		return object.Null, fmt.Errorf("can not iterate over %s", intLeft.Class().Name())
	}

	l, r := intLeft.Value, intRight.Value
	array := &object.ArrayObject{}

	if l == r {
		return array, nil
	}
	if l < r {
		for l < r {
			array.Values = append(array.Values, &object.IntegerObject{Value: l})
			l++
		}
		return array, nil
	}
	for l > r {
		array.Values = append(array.Values, &object.IntegerObject{Value: l})
		l++
	}
	return array, nil
}

// evalFetchExpression ...
func (ev *evaluator) evalFetchExpression(ex *ast.FetchExpression, ctx object.Context) (object.Object, error) {
	obj, err := ev.Eval(ex.Left, ctx)
	if err != nil {
		return object.Null, err
	}

	switch r := ex.Right.(type) {
	case *ast.FunctionCall:
		return ev.evalMethodCall(obj, r, ctx)
	default:
		return object.Null, errors.New("unexpected right node")
	}
}

// evalMethodCall ...
func (ev *evaluator) evalMethodCall(obj object.Object, node *ast.FunctionCall, ctx object.Context) (object.Object, error) {
	methodName, ok := node.Target.(*ast.Identifier)
	if !ok {
		return object.Null, errors.New("method name must be an Identifier")
	}
	method := obj.Class().Methods().Find(methodName.Value)
	if method == nil {
		return object.Null, fmt.Errorf("method %s is not found in class %s", methodName.Value, obj.Class().Name())
	}
	args := make([]object.Object, len(node.CallArgs))
	for i, expr := range node.CallArgs {
		obj, err := ev.Eval(expr, ctx)
		if err != nil {
			return object.Null, err
		}
		args[i] = obj
	}
	return method.Call(obj, args...)
}
