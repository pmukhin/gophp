package object

import "os"

func registerOsConstants(ctx Context) {
	ctx.SetGlobal("os\\args", NewInternalFunc(func(args ...Object) (Object, error) {
		osArgs := os.Args[1:] // eat first arg
		vars := make([]Object, len(osArgs))
		for i := range osArgs {
			vars[i] = &StringObject{Value: osArgs[i]}
		}
		return &ArrayObject{Values: vars}, nil
	}))
}
