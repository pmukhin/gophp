package eval

import (
	"github.com/pmukhin/gophp/object"
	"github.com/pmukhin/gophp/ast"
)

func evalBlock(ctx object.Context, node *ast.BlockStatement) (object.Object, error) {
	var ret object.Object = nil
	for _, st := range node.Statements {
		r, e := Eval(st, ctx)
		if e != nil {
			return object.Null, e
		}
		// if it's return
		if ret, ok := r.(returnObject); ok {
			// unpack return object first...
			return ret.value, nil
		}
		// add
		ret = r
	}

	if ret == nil {
		return object.Null, nil
	}
	return ret, nil
}
