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
			return returnObject{value: object.Null}, e
		}
		// if it's return
		if ret, ok := r.(returnObject); ok {
			return ret, nil
		}
		// add
		ret = r
	}

	if ret == nil {
		ret = returnObject{value: object.Null}
	}

	return ret, nil
}
