package object

import "math/rand"

func registerMathFunctions(ctx Context) {
	ctx.SetGlobal("math\\random", NewInternalFunc("math\\random", func(ctx Context, args ...Object) (Object, error) {
		r := rand.Int63()
		return &IntegerObject{Value: r}, nil
	}))
}