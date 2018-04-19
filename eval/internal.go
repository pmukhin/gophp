package eval

import (
	"github.com/pmukhin/gophp/object"
	"fmt"
	"os"
	"errors"
)

var InternalFunctionTable *FunctionTable

func init() {
	InternalFunctionTable = NewFunctionTable()
	// println(...$args)
	InternalFunctionTable.RegisterFunc("println", object.NewInternalFunc("println", func(args ...object.Object) (object.Object, error) {
		if len(args) < 0 {
			return object.Null, errors.New("println expects at least 1 argument")
		}
		for _, a := range args {
			toString := a.Class().Methods().Find("__toString")
			if toString == nil {
				fmt.Printf("%v", a)
			} else {
				s, e := toString.Call(a)
				if e != nil {
					return object.Null, e
				}
				// @todo print String object
				fmt.Printf("%v", s)
			}
		}
		fmt.Println()
		return object.Null, nil
	}))

	// exit()
	InternalFunctionTable.RegisterFunc("exit", object.NewInternalFunc("exit", func(args ...object.Object) (object.Object, error) {
		os.Exit(-0)
		// for compiler
		return object.Null, nil
	}))
}

type FunctionTable struct {
	table map[string]object.Callable
}

func (ft *FunctionTable) RegisterFunc(name string, fun object.Callable) {
	ft.table[name] = fun
}

func (ft *FunctionTable) Find(name string) (object.Callable, error) {
	if f, ok := ft.table[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("function %s is not defined", name)
}

func NewFunctionTable() *FunctionTable {
	return &FunctionTable{table: make(map[string]object.Callable)}
}
