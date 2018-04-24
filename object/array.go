package object

import (
	"fmt"
	"strings"
)

func arrayToString(this Object, args ...Object) (Object, error) {
	a := this.(*ArrayObject)
	strList := make([]string, len(a.Values))

	for i, v := range a.Values {
		s, e := ToString(v)
		if e == nil {
			strList[i] = s.Value
		} else {
			strList[i] = fmt.Sprintf("%v", v)
		}
	}

	return &StringObject{Value: "[" + strings.Join(strList, ", ") + "]"}, nil
}

func arrayLen(this Object, args ...Object) (Object, error) {
	return &IntegerObject{Value: int64(len(this.(*ArrayObject).Values))}, nil
}

func arrayAppend(this Object, args ...Object) (Object, error) {
	if len(args) == 0 {
		return Null, fmt.Errorf("at least 1 argument expected")
	}
	array := this.(*ArrayObject)
	for _, a := range args {
		array.Values = append(array.Values, a)
	}
	return Null, nil
}

var (
	arrayMethods = map[string]Method{
		"__toString": newMethod(arrayToString, VisibilityPublic),

		"length": newMethod(arrayLen, VisibilityPublic),
		"append": newMethod(arrayAppend, VisibilityPublic),
	}

	arrayClass = &InternalClass{
		name:      "Array",
		final:     false,
		abstract:  false,
		methodSet: newMethodSet(arrayMethods),
	}
)

func registerArrayConstants(ctx Context) {
	ctx.SetGlobal("Array", arrayClass)
}

// NewArray ...
func NewArray(os ...Object) (Object, error) {
	array := new(ArrayObject)
	array.Values = make([]Object, len(os), 32)

	for i, ob := range os {
		array.Values[i] = ob
	}

	return array, nil
}

type ArrayObject struct {
	Values []Object
}

func (ArrayObject) Class() Class { return arrayClass }

func (ArrayObject) Id() string {
	panic("implement me")
}
