package object

import (
	"fmt"
	"os"
	"errors"
)

var InternalFunctionTable *FunctionTable

func doPrint(delim string) func(ctx Context, args ...Object) (Object, error) {
	return func(ctx Context, args ...Object) (Object, error) {
		if len(args) < 0 {
			return Null, errors.New("println expects at least 1 argument")
		}
		for _, a := range args {
			s, e := ToString(a)
			if e != nil {
				fmt.Printf("<%s> %v ", s.Class().Name(), s.Class())
				continue
			}
			fmt.Print(s.Value + " ")
		}
		fmt.Print(delim)
		return Null, nil
	}
}

func init() {
	InternalFunctionTable = NewFunctionTable()
	// print(...$args)
	InternalFunctionTable.RegisterFunc(NewInternalFunc("print", doPrint("")))
	// println(...$args)
	InternalFunctionTable.RegisterFunc(NewInternalFunc("println", doPrint("\n")))

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

func (ft *FunctionTable) Find(ns, name string) (FunctionObject, error) {
	var uf string

	if ns == "" {
		uf = name
	} else {
		uf = ns + "\\" + name
	}
	// first let's look in the same namespace
	if f, ok := ft.table[uf]; ok {
		return f, nil
	}
	// then global namespace
	if f, ok := ft.table[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("function %s is not defined", uf)
}

func NewFunctionTable() *FunctionTable {
	return &FunctionTable{table: make(map[string]FunctionObject)}
}
