package object

import "math/rand"

func registerMathFunctions(ctx Context) {
	ctx.SetGlobal("math\\random", NewInternalFunc(func(args ...Object) (Object, error) {
		r := rand.Int63()
		return &IntegerObject{Value: r}, nil
	}))
}
