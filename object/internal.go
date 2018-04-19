package object

import (
	"fmt"
	"os"
	"errors"
)

var InternalFunctionTable *FunctionTable

func init() {
	InternalFunctionTable = NewFunctionTable()
	// println(...$args)
	InternalFunctionTable.RegisterFunc(NewInternalFunc("println", func(ctx Context, args ...Object) (Object, error) {
		if len(args) < 0 {
			return Null, errors.New("println expects at least 1 argument")
		}
		for _, a := range args {
			toString := a.Class().Methods().Find("__toString")
			if toString == nil {
				fmt.Printf("%v", a)
			} else {
				s, e := toString.Call(a)
				if e != nil {
					return Null, e
				}
				// @todo print String object
				fmt.Printf("%v", s)
			}
		}
		fmt.Println()
		return Null, nil
	}))

	// exit()
	InternalFunctionTable.RegisterFunc(NewInternalFunc("exit", func(ctx Context, args ...Object) (Object, error) {
		os.Exit(-0)
		// for compiler
		return Null, nil
	}))
}

type FunctionTable struct {
	table map[string]FunctionObject
}

func (ft *FunctionTable) RegisterFunc(fun FunctionObject) error {
	if _, ok := ft.table[fun.Name()]; ok {
		return fmt.Errorf("function %s is already defined", fun.Name())
	}
	ft.table[fun.Name()] = fun
	return nil
}

func (ft *FunctionTable) Find(name string) (FunctionObject, error) {
	if f, ok := ft.table[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("function %s is not defined", name)
}

func NewFunctionTable() *FunctionTable {
	return &FunctionTable{table: make(map[string]FunctionObject)}
}
