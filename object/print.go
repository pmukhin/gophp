package object

import (
	"fmt"
	"os"
	"errors"
)

func doPrint(delimiter string) func(args ...Object) (Object, error) {
	return func(args ...Object) (Object, error) {
		if len(args) < 0 {
			return Null, errors.New("println expects at least 1 argument")
		}
		for _, a := range args {
			s, e := ToString(a)
			if e != nil {
				fmt.Printf("<%s> %v ", a.Class().Name(), s.Class())
				continue
			}
			fmt.Print(s.Value)
		}
		fmt.Print(delimiter)
		return Null, nil
	}
}

func registerPrintFunctions(ctx Context) {
	ctx.SetGlobal("print", NewInternalFunc(doPrint("")))
	ctx.SetGlobal("println", NewInternalFunc(doPrint("\n")))
	ctx.SetGlobal("exit", NewInternalFunc(func(args ...Object) (Object, error) {
		os.Exit(0)
		// for compiler
		return Null, nil
	}))
}
