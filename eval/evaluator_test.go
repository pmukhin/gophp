package eval

import (
	"github.com/pmukhin/gophp/ast"
	"github.com/pmukhin/gophp/object"
	"github.com/pmukhin/gophp/parser"
	"github.com/pmukhin/gophp/scanner"
	"testing"
)

func checkContextVariableInt(t *testing.T, ctx object.Context, name string, value int64) {
	if el, err := ctx.GetContextVar(name); err != nil {
		t.Error(err)
	} else {
		if e, ok := el.(*object.IntegerObject); !ok {
			t.Errorf("%v is not Integer, it's %s", e.Class().Name())
		} else {
			if e.Value != value {
				t.Errorf("$%s is %d, not %d", name, e.Value, value)
			}
		}
	}
}

func TestEval(t *testing.T) {
	// tested code:
	// $variableInteger = 5;
	// $secondVariableInteger = 5;
	// $sum = $variableInteger + $secondVariableInteger;
	// $division = $variableInteger / $secondVariableInteger;
	// $multiplication = $variableInteger * $secondVariableInteger;
	// $sub = $variableInteger - $secondVariableInteger;
	t.Run("assignment & arithmetic", func(t *testing.T) {
		p := &ast.Module{}
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
		ctx := object.NewContext(nil)
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
		p := &ast.Module{}
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
		ctx := object.NewContext(nil)
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

	t.Run("print string var", func(t *testing.T) {
		p := parser.New(scanner.New([]rune(`
			println(5)
		`)))
		program, e := p.Parse()
		if e != nil {
			t.Error(e)
		}
		ctx := object.NewContext(nil)
		_, e = Eval(program, ctx)
		if e != nil {
			t.Error(e)
		}
	})

	t.Run("print string var", func(t *testing.T) {
		p := parser.New(scanner.New([]rune(`
			$str = 'hello'
			println($str)
		`)))
		program, e := p.Parse()
		if e != nil {
			t.Error(e)
		}
		ctx := object.NewContext(nil)
		_, e = Eval(program, ctx)
		if e != nil {
			t.Error(e)
		}
	})

	t.Run("evaluate user function", func(t *testing.T) {
		p := parser.New(scanner.New([]rune(`
			function test() {}
			$var = test();
			println($var)
		`)))
		program, e := p.Parse()
		if e != nil {
			t.Error(e)
		}
		ctx := object.NewContext(nil)
		_, e = Eval(program, ctx)
		if e != nil {
			t.Error(e)
		}
	})

	t.Run("evaluate bool expression and print", func(t *testing.T) {
		p := parser.New(scanner.New([]rune(`
			$is = 5 > 5
			println($is)
`)))
		program, e := p.Parse()
		if e != nil {
			t.Error(e)
		}
		ctx := object.NewContext(nil)
		_, e = Eval(program, ctx)
		if e != nil {
			t.Error(e)
		}
	})

	t.Run("assignment of conditional", func(t *testing.T) {
		p := &ast.Module{}
		p.Statements = []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "result"},
					Right: &ast.ConditionalExpression{
						Condition: &ast.Identifier{Value: "true"},
						Consequence: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Expression: &ast.IntegerLiteral{Value: 5},
								},
							},
						},
						Alternative: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Expression: &ast.IntegerLiteral{Value: 0},
								},
							},
						},
					},
				},
			},
		}
		ctx := object.NewContext(nil)
		_, e := Eval(p, ctx)
		if e != nil {
			t.Error(e)
		}
		checkContextVariableInt(t, ctx, "result", 5)
	})

	// tested code:
	// $result = println(7);
	t.Run("simplest function call", func(t *testing.T) {
		p := &ast.Module{}
		p.Statements = []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.AssignmentExpression{
					Left: &ast.VariableExpression{Name: "result"},
					Right: &ast.FunctionCall{
						Target: &ast.Identifier{Value: "println"},
						CallArgs: []ast.Expression{
							&ast.IntegerLiteral{Value: 7},
						},
					},
				},
			},
		}
		ctx := object.NewContext(nil)
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

func TestEval_Fibonacci_OldSyntax(t *testing.T) {
	p := parser.New(scanner.New([]rune(`
			function fib($n) { 
				if $n < 2 { 
					return 1; 
				} else { 
					return fib($n - 2) + fib($n - 1); 
				} 
			}
			$result1 = fib(1);
			$result2 = fib(2);
			$result3 = fib(3);
	`)))
	program, e := p.Parse()
	if e != nil {
		t.Error(e)
	}
	ctx := object.NewContext(nil)
	_, e = Eval(program, ctx)
	if e != nil {
		t.Error(e)
	}
	checkContextVariableInt(t, ctx, "result1", 1)
	checkContextVariableInt(t, ctx, "result2", 2)
	checkContextVariableInt(t, ctx, "result3", 3)
}

func TestEval_Fibonacci_NewSyntax(t *testing.T) {
	p := parser.New(scanner.New([]rune(`
			function fib($n) { if $n < 2 { $n } else { fib($n - 2) + fib($n - 1) } }
			$result1 = fib(0)
			$result2 = fib(1)
			$result3 = fib(2)
			$result4 = fib(3)
	`)))
	program, e := p.Parse()
	if e != nil {
		t.Error(e)
	}
	ctx := object.NewContext(nil)
	_, e = Eval(program, ctx)
	if e != nil {
		t.Error(e)
	}
	checkContextVariableInt(t, ctx, "result1", 0)
	checkContextVariableInt(t, ctx, "result2", 1)
	checkContextVariableInt(t, ctx, "result3", 1)
	checkContextVariableInt(t, ctx, "result4", 2)

}
