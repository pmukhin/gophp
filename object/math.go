package object

import "math/rand"

func init() {
	InternalFunctionTable.RegisterFunc(NewInternalFunc("math\\random", func(ctx Context, args ...Object) (Object, error) {
		r := rand.Int63()
		return &IntegerObject{Value: r}, nil
	}))
}
