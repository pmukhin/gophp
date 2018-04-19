package object

import "fmt"

func infer(this Object, os ...Object) (*IntegerObject, *IntegerObject, error) {
	if len(os) != 1 {
		return nil, nil, fmt.Errorf("__add takes exactly one parameter, %d given", len(os))
	}
	i, ok := this.(*IntegerObject)
	if !ok {
		panic("?")
	}
	switch right := os[0].(type) {
	case *IntegerObject:
		return i, right, nil
	default:
		return nil, nil, fmt.Errorf("unsupported operand %v", os[0])
	}
}

func iEqual(this Object, os ...Object) (Object, error) {
	l, r, e := infer(this, os...)
	if e != nil {
		return Null, e
	}
	if l.Value == r.Value {
		return True, nil
	}
	return False, nil
}

func isGreater(this Object, os ...Object) (Object, error) {
	l, r, e := infer(this, os...)
	if e != nil {
		return Null, e
	}
	if l.Value > r.Value {
		return True, nil
	}
	return False, nil
}

func iAdd(this Object, os ...Object) (Object, error) {
	l, r, e := infer(this, os...)
	if e != nil {
		return Null, e
	}
	return &IntegerObject{Value: l.Value + r.Value}, nil
}

func iSub(this Object, os ...Object) (Object, error) {
	l, r, e := infer(this, os...)
	if e != nil {
		return Null, e
	}
	return &IntegerObject{Value: l.Value - r.Value}, nil
}

func iMul(this Object, os ...Object) (Object, error) {
	l, r, e := infer(this, os...)
	if e != nil {
		return Null, e
	}
	return &IntegerObject{Value: l.Value * r.Value}, nil
}

func iDiv(this Object, os ...Object) (Object, error) {
	l, r, e := infer(this, os...)
	if e != nil {
		return Null, e
	}
	if r.Value == 0 {
		return Null, fmt.Errorf("division by zero is forbidden")
	}
	return &IntegerObject{Value: l.Value / r.Value}, nil
}

func iToBoolean(this Object, os ...Object) (Object, error) {
	if this.(*IntegerObject).Value == 0 {
		return False, nil
	}
	return True, nil
}

var (
	ic = func(value interface{}) (Object, error) {
		v, ok := value.(int64)
		if !ok {
			return nil, fmt.Errorf("%v is not an integer", value)
		}
		return &IntegerObject{Value: v}, nil
	}

	integerMethodsMap = map[string]Method{
		"__add":       newMethod(iAdd, VisibilityPublic),
		"__sub":       newMethod(iSub, VisibilityPublic),
		"__mul":       newMethod(iMul, VisibilityPublic),
		"__div":       newMethod(iDiv, VisibilityPublic),
		"__equal":     newMethod(iEqual, VisibilityPublic),
		"__gt":        newMethod(isGreater, VisibilityPublic),
		"__toBoolean": newMethod(iToBoolean, VisibilityPublic),
	}

	//IntegerClass = newInternalClass("Integer", true, false, integerConstructor{}, InternalConstructor(ic))
	IntegerClass = internalClass{
		name:                "Integer",
		final:               true,
		abstract:            false,
		constructor:         integerConstructor{},
		internalConstructor: ic,
		methodSet:           newMethodSet(integerMethodsMap),
	}
)

type IntegerObject struct {
	Value int64
}

func (IntegerObject) Class() Class {
	return IntegerClass
}

func (IntegerObject) Id() string {
	panic("implement me")
}

type integerConstructor struct{}

func (integerConstructor) Call(this Object, object ...Object) (Object, error) {
	panic("implement me")
}

func (integerConstructor) Visibility() Visibility {
	return VisibilityPublic
}
