package object

import (
	"strings"
	"strconv"
	"errors"
)

func stringConcat(this Object, args ...Object) (Object, error) {
	l := this.(*StringObject)
	if len(args) != 1 {
		panic("error")
	}
	arg, e := ToString(args[0])
	if e != nil {
		return Null, e
	}
	return &StringObject{Value: l.Value + arg.Value}, nil
}

func repeat(this Object, args ...Object) (Object, error) {
	l := this.(*StringObject)
	if len(args) != 1 {
		panic("error")
	}
	arg, e := ToInteger(args[0])
	if e != nil {
		return Null, e
	}
	return &StringObject{Value: strings.Repeat(l.Value, int(arg.Value))}, nil
}

func toInt(this Object, args ...Object) (Object, error) {
	l := this.(*StringObject)
	if len(args) != 0 {
		panic("error")
	}
	i, e := strconv.ParseInt(l.Value, 10, 64)
	if e != nil {
		return nil, e
	}
	return &IntegerObject{Value: i}, nil
}

func index(this Object, args ...Object) (Object, error) {
	l := this.(*StringObject)
	if len(args) != 1 {
		panic("error")
	}
	arg, e := ToInteger(args[0])
	if e != nil {
		return nil, e
	}
	r := []rune(l.Value)
	if len(r) <= int(arg.Value) {
		return Null, errors.New("index error")
	}

	return &StringObject{Value: string(r[arg.Value])}, nil
}

var (
	m = map[string]Method{
		"__add":   newMethod(stringConcat, VisibilityPublic),
		"__mul":   newMethod(repeat, VisibilityPublic),
		"__toInt": newMethod(toInt, VisibilityPublic),
		"__index": newMethod(index, VisibilityPublic),
	}

	stringClass = internalClass{
		name:      "String",
		final:     true,
		abstract:  false,
		methodSet: newMethodSet(m),
	}
)

// StringObject ...
type StringObject struct {
	Value string
}

func (StringObject) Class() Class {
	return stringClass
}

func (StringObject) Id() string { panic("implement me") }
