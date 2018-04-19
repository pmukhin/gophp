package object

import "fmt"

func ToString(o Object) (*StringObject, error) {
	if o.Class().Name() == "String" {
		return o.(*StringObject), nil
	}
	toString := o.Class().Methods().Find("__toString")
	if toString == nil {
		return nil, fmt.Errorf("%v can not be converted to String", o)
	}
	argStr, e := toString.Call(o)
	if e != nil {
		return nil, e
	}
	return argStr.(*StringObject), nil
}

func ToInteger(o Object) (*IntegerObject, error) {
	if o.Class().Name() == "Integer" {
		return o.(*IntegerObject), nil
	}
	toInt := o.Class().Methods().Find("__toInt")
	if toInt == nil {
		return nil, fmt.Errorf("%v can not be converted to Int", o)
	}
	argStr, e := toInt.Call(o)
	if e != nil {
		return nil, e
	}
	return argStr.(*IntegerObject), nil
}
