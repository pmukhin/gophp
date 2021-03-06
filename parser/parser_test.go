package parser

import (
	"testing"
	"github.com/pmukhin/gophp/scanner"
	"github.com/pmukhin/gophp/error"
	"github.com/pmukhin/gophp/ast"
	"reflect"
	"strings"
	"github.com/pmukhin/gophp/token"
	"fmt"
)

func run(t *testing.T, input string, expectations []ast.Statement) {
	scn := scanner.New([]rune(input))
	parser := New(scn, error.NewFormatter("<test>", []rune(input)))

	program, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(program.Statements)
	return
	if !reflect.DeepEqual(expectations, program.Statements) {
		t.Errorf("expectations: \n%v, \ngot: \n%v", expectations, program.Statements)
	}
	if scn.HasNext() {
		r := make([]string, 32)
		for {
			t := scn.Next()
			if t.Type == token.EOF {
				break
			}
			r = append(r, t.Literal)
		}

		t.Errorf("scanner still has tokens: %s", strings.Join(r, "; "))
	}
}

func TestParser_Parse_ArithmeticPrecedence(t *testing.T) {
	input := []rune(`println(5 + 5 * 3)`)
	scn := scanner.New(input)
	parser := New(scn, error.NewFormatter("<test>", input))

	_, e := parser.Parse()
	if e != nil {
		t.Error(e)
	}
}

func TestParser_Parse_Incomplete(t *testing.T) {
	input := []rune(`function fib($n) { if $n < 2 { $n } else { fib($n-1) + fib($n-2) }`)
	scn := scanner.New(input)
	parser := New(scn, error.NewFormatter("<test>", input))

	_, err := parser.Parse()
	if err == nil {
		t.Errorf("must return emitError")
	}
	if err != nil && !strings.Contains(err.Error(), "unexpected token EOF") {
		t.Errorf("expected EOF error got: %s", err.Error())
	}
}

func TestParser_ParseSimpleConditional_FunctionDeclaration(t *testing.T) {
	input := `function first() {}`
	run(t, input, []ast.Statement{
		&ast.ExpressionStatement{
			Expression: &ast.FunctionDeclarationExpression{
				Name:  &ast.Identifier{Value: "first"},
				Args:  []*ast.Arg{},
				Block: &ast.BlockStatement{},
			},
		},
	})
}

func TestParser_ParseSimpleConditional_NewStyle(t *testing.T) {
	input := `if (true) { 5 } else { 0 }`
	run(t, input, []ast.Statement{
		&ast.ExpressionStatement{
			Expression: &ast.ConditionalExpression{
				Condition: &ast.BooleanExpression{Value: true},
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
	})
}

func TestParse_ParseFunctionCall(t *testing.T) {
	input := `println($var);`
	run(t, input, []ast.Statement{
		&ast.ExpressionStatement{
			Expression: &ast.FunctionCall{
				Target: &ast.Identifier{Value: "println"},
				CallArgs: []ast.Expression{
					&ast.VariableExpression{Name: "var"},
				},
			},
		},
	})
}

func TestParser_ParseAssignment(t *testing.T) {
	input := `$var = 5;`
	run(t, input, []ast.Statement{
		&ast.ExpressionStatement{
			Expression: &ast.AssignmentExpression{
				Left:  &ast.VariableExpression{Name: "var"},
				Right: &ast.IntegerLiteral{Value: 5},
			},
		},
	})
}

func TestParser_ParseNamespaceStatement(t *testing.T) {
	run(t, `namespace Silex;`, []ast.Statement{
		ast.NamespaceStatement{Namespace: "Silex"},
	})
}

func TestParser_ParseUseStatement(t *testing.T) {
	run(t, `use Symfony\Component\Debug\Exception\FlattenException;`, []ast.Statement{
		ast.UseStatement{Namespace: "Symfony\\Component\\Debug\\Exception", Classes: []string{"FlattenException"}},
	})
}
