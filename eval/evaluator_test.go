package eval

import (
	"testing"
	"github.com/pmukhin/gophp/ast"
	"github.com/pmukhin/gophp/object"
)

func TestEval(t *testing.T) {
	// tested code:
	// $variableInteger = 5;
	// $secondVariableInteger = 5;
	// $sum = $variableInteger + $secondVariableInteger;
	// $division = $variableInteger / $secondVariableInteger;
	// $multiplication = $variableInteger * $secondVariableInteger;
	// $sub = $variableInteger - $secondVariableInteger;
	t.Run("assignment & arithmetic", func(t *testing.T) {
		p := &ast.Program{}
		p.Statements = []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left:  &ast.VariableExpression{Name: "variableInteger"},
					Right: &ast.IntegerLiteral{Value: 5},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left:  &ast.VariableExpression{Name: "secondVariableInteger"},
					Right: &ast.IntegerLiteral{Value: 5},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "sum"},
					Right: &ast.BinaryExpression{
						Op:    "+",
						Left:  &ast.VariableExpression{Name: "variableInteger"},
						Right: &ast.VariableExpression{Name: "secondVariableInteger"},
					},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "division"},
					Right: &ast.BinaryExpression{
						Op:    "/",
						Left:  &ast.VariableExpression{Name: "variableInteger"},
						Right: &ast.VariableExpression{Name: "secondVariableInteger"},
					},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "multiplication"},
					Right: &ast.BinaryExpression{
						Op:    "*",
						Left:  &ast.VariableExpression{Name: "variableInteger"},
						Right: &ast.VariableExpression{Name: "secondVariableInteger"},
					},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "sub"},
					Right: &ast.BinaryExpression{
						Op:    "-",
						Left:  &ast.VariableExpression{Name: "variableInteger"},
						Right: &ast.VariableExpression{Name: "secondVariableInteger"},
					},
				},
			},
		}
		ctx := NewContext(nil, InternalFunctionTable)
		_, e := Eval(p, ctx)
		if e != nil {
			t.Error(e)
		}
		if v, e := ctx.GetContextVar("variableInteger"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.IntegerObject).Value != 5 {
				t.Errorf("expected $variableInteger to equal to %d", 5)
			}
		}
		if v, e := ctx.GetContextVar("secondVariableInteger"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.IntegerObject).Value != 5 {
				t.Errorf("expected $secondVariableInteger to equal to %d", 5)
			}
		}
		// sum
		if v, e := ctx.GetContextVar("sum"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.IntegerObject).Value != 10 {
				t.Errorf("expected $sum to equal to %d but got %d", 10, v.(*object.IntegerObject).Value)
			}
		}
		// division
		if v, e := ctx.GetContextVar("division"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.IntegerObject).Value != 1 {
				t.Errorf("expected $division to equal to %d but got %d", 10, v.(*object.IntegerObject).Value)
			}
		}
		// multiplication
		if v, e := ctx.GetContextVar("multiplication"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.IntegerObject).Value != 25 {
				t.Errorf("expected $multiplication to equal to %d but got %d", 10, v.(*object.IntegerObject).Value)
			}
		}
		// sub
		if v, e := ctx.GetContextVar("sub"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.IntegerObject).Value != 0 {
				t.Errorf("expected $sub to equal to %d but got %d", 0, v.(*object.IntegerObject).Value)
			}
		}
	})

	// tested code
	// $isEqual = 5 == 5;
	t.Run("is equal", func(t *testing.T) {
		p := &ast.Program{}
		p.Statements = []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "isEqual"},
					Right: &ast.BinaryExpression{
						Left:  &ast.IntegerLiteral{Value: 5},
						Op:    "==",
						Right: &ast.IntegerLiteral{Value: 5},
					},
				},
			},
		}
		ctx := NewContext(nil, InternalFunctionTable)
		_, e := Eval(p, ctx)
		if e != nil {
			t.Error(e)
		}
		if v, e := ctx.GetContextVar("isEqual"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.BooleanObject) != object.True {
				t.Errorf("expected $isEqual to be true")
			}
		}
	})

	// tested code:
	// $result = println(7);
	t.Run("simplest function call", func(t *testing.T) {
		p := &ast.Program{}
		p.Statements = []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "result"},
					Right: &ast.FunctionCall{
						Target: ast.Identifier{Value: "println"},
						CallArgs: []ast.Expression{
							&ast.IntegerLiteral{Value: 7},
						},
					},
				},
			},
		}
		ctx := NewContext(nil, InternalFunctionTable)
		_, e := Eval(p, ctx)
		if e != nil {
			t.Error(e)
		}
		if v, e := ctx.GetContextVar("result"); e != nil {
			t.Error(e)
		} else {
			if v.(*object.NullObject) != object.Null {
				t.Errorf("expected $result to be null")
			}
		}
	})
}
