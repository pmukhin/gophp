package parser

import (
	"testing"
	"github.com/pmukhin/gophp/scanner"
	"github.com/pmukhin/gophp/ast"
	"reflect"
	"strings"
	"github.com/pmukhin/gophp/token"
	"fmt"
)

func run(t *testing.T, input string, expectations []ast.Statement) {
	scn := scanner.New([]rune(input))
	parser := Parser{__scn: scn}
	parser.init()

	program, err := parser.Parse()
	fmt.Println(program.Statements)
	if err != nil {
		t.Fatal(err)
	}
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

func TestParser_ParseFunction(t *testing.T) {
	input := `
		$someVar = $response->statusCode;
		function render($view, array $parameters): Response {
        $twig = $this['twig'];
        if ($response instanceof StreamedResponse) {
            $response->setCallback(function () use ($twig, $view) {
                $twig->display($view);
            });
        } else {
            $response;
        }
        return $response;
    }`

	run(t, input, []ast.Statement{
		ast.ExpressionStatement{
			Expression: ast.FunctionDeclarationExpression{
				Anonymous:  false,
				Name:       ast.Identifier{Value: "render"},
				ReturnType: &ast.Identifier{Value: "Response"},
				Args: []ast.Arg{
					{
						Type: nil,
						Name: ast.VariableExpression{Name: "view"},
					},
					{
						Type: &ast.Identifier{Value: "array"},
						Name: ast.VariableExpression{Name: "parameters"},
					},
				},
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
