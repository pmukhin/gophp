package eval

import (
	"github.com/pmukhin/gophp/object"
	"fmt"
)

var InternalFunctionTable *FunctionTable

func init() {
	InternalFunctionTable = NewFunctionTable()
	InternalFunctionTable.RegisterFunc("println", object.NewInternalFunc("println", func(args ...object.Object) (object.Object, error) {
		for _, a := range args {
			toString := a.Class().Methods().Find("__toString")
			if toString == nil {
				fmt.Printf("%v", a)
			} else {
				fmt.Println(toString.Call(a))
			}
		}
		fmt.Println()
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
