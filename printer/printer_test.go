package printer

import (
	"testing"
	"github.com/pmukhin/gophp/ast"
	"reflect"
)

func TestPrinter_Result_ListOfStatements(t *testing.T) {
	program := &ast.Module{}
	program.Statements = []ast.Statement{
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
	p := New()
	p.Visit(ast.ProgramType, program)

	expectedResult := `$variableInteger = 5;
$secondVariableInteger = 5;
$sum = $variableInteger + $secondVariableInteger;
$division = $variableInteger / $secondVariableInteger;
$multiplication = $variableInteger * $secondVariableInteger;
$sub = $variableInteger - $secondVariableInteger;`

	if !reflect.DeepEqual(p.Result(), expectedResult) {
		t.Errorf("p.Result() = %v, expected %v", p.Result(), expectedResult)
	}
}
