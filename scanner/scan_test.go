package scanner

import (
	"testing"
	"lang/token"
	"reflect"
)

func run(t *testing.T, input string, expectations []Token) {
	s := New([]rune(input))
	tokens := make([]Token, 0, 2048)

	for {
		tok := s.Next()
		if tok == tokenEof {
			break
		}
		tokens = append(tokens, tok)
	}

	if !reflect.DeepEqual(expectations, tokens) {
		t.Errorf("expected: \n%v, \ngot: \n%v", expectations, tokens)
	}
}

func TestScanner_Next_If(t *testing.T) {
	run(t, "if () {} else if () {} else {}", []Token{
		{Type: token.IF, Literal: "if"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		{Type: token.CURLY_CLOSING, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.IF, Literal: "if"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		{Type: token.CURLY_CLOSING, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		{Type: token.CURLY_CLOSING, Literal: "}"},
	})
}

func TestScanner_Next_Throw(t *testing.T) {
	run(t, "throw new HttpException", []Token{
		{Type: token.THROW, Literal: "throw"},
		{Type: token.NEW, Literal: "new"},
		{Type: token.IDENT, Literal: "HttpException"},
	})
}

func TestScanner_Next_Arithmetic(t *testing.T) {
	run(t, "$i % 2", []Token{
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.MOD, Literal: "%"},
		{Type: token.NUMBER, Literal: "2"},
	})

	run(t, "$i += 2", []Token{
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.PLUS_EQUAL, Literal: "+="},
		{Type: token.NUMBER, Literal: "2"},
	})

	run(t, "$i -= 2", []Token{
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.MINUS_EQUAL, Literal: "-="},
		{Type: token.NUMBER, Literal: "2"},
	})

	run(t, "$i /= 2", []Token{
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.DIV_EQUAL, Literal: "/="},
		{Type: token.NUMBER, Literal: "2"},
	})

	run(t, "$i *= 2", []Token{
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.MUL_EQUAL, Literal: "*="},
		{Type: token.NUMBER, Literal: "2"},
	})
}

func TestScanner_Next_LoopIncDec(t *testing.T) {
	var input = `
		++$i
		$i++
		--$i
		$i--`
	run(t, input, []Token{
		{Type: token.INC, Literal: "++"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.INC, Literal: "++"},
		{Type: token.DEC, Literal: "--"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.DEC, Literal: "--"},
	})
}

func TestScanner_Next_LoopFor(t *testing.T) {
	var input = `
for($i = 0; $i <= 5; ++$i) {}
`
	run(t, input, []Token{
		{Type: token.FOR, Literal: "for"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.NUMBER, Literal: "0"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.IS_SMALLER_OR_EQUAL, Literal: "<="},
		{Type: token.NUMBER, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INC, Literal: "++"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "i"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		{Type: token.CURLY_CLOSING, Literal: "}"},
	})
}

func TestScanner_Next_LoopForeach(t *testing.T) {
	var input = `
foreach ($array as $key => $value) {}
`
	run(t, input, []Token{
		{Type: token.FOREACH, Literal: "foreach"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "array"},
		{Type: token.AS, Literal: "as"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "key"},
		{Type: token.DOUBLE_ARROW, Literal: "=>"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "value"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		{Type: token.CURLY_CLOSING, Literal: "}"},
	})
}

func TestScanner_Next_Class(t *testing.T) {
	var input = `
namespace Silex;

class Application extends Container implements HttpKernelInterface {
    const LATE_EVENT = -512;
    protected $booted = false;

	public function register(array $values = []) {
        parent::register($values);
    }
}`
	expectations := []Token{
		// namespace Silex;
		{Type: token.NAMESPACE, Literal: "namespace"},
		{Type: token.IDENT, Literal: "Silex"},
		{Type: token.SEMICOLON, Literal: ";"},
		// class Application extends Container implements HttpKernelInterface {
		{Type: token.CLASS, Literal: "class"},
		{Type: token.IDENT, Literal: "Application"},
		{Type: token.EXTENDS, Literal: "extends"},
		{Type: token.IDENT, Literal: "Container"},
		{Type: token.IMPLEMENTS, Literal: "implements"},
		{Type: token.IDENT, Literal: "HttpKernelInterface"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		// const LATE_EVENT = -512;
		{Type: token.CONST, Literal: "const"},
		{Type: token.IDENT, Literal: "LATE_EVENT"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.NUMBER, Literal: "-512"},
		{Type: token.SEMICOLON, Literal: ";"},
		// protected $booted = false;
		{Type: token.PROTECTED, Literal: "protected"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "booted"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.IDENT, Literal: "false"},
		{Type: token.SEMICOLON, Literal: ";"},
		// public function register(array $values = []) {
		{Type: token.PUBLIC, Literal: "public"},
		{Type: token.FUNCTION, Literal: "function"},
		{Type: token.IDENT, Literal: "register"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.IDENT, Literal: "array"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "values"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.SQUARE_BRACKET_OPENING, Literal: "["},
		{Type: token.SQUARE_BRACKET_CLOSING, Literal: "]"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		// parent::register($values);
		{Type: token.IDENT, Literal: "parent"},
		{Type: token.PAAMAYIM_NEKUDOTAYIM, Literal: "::"},
		{Type: token.IDENT, Literal: "register"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "values"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		// }
		{Type: token.CURLY_CLOSING, Literal: "}"},
		// }
		{Type: token.CURLY_CLOSING, Literal: "}"},
	}
	run(t, input, expectations)
}

func TestScanner_Next_Function(t *testing.T) {
	var input = `function render($view, array $parameters = []): Response
    {
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

	expectations := []Token{
		// function render($view, array $parameters = []): Response
		{Type: token.FUNCTION, Literal: "function"},
		{Type: token.IDENT, Literal: "render"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "view"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "array"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "parameters"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.SQUARE_BRACKET_OPENING, Literal: "["},
		{Type: token.SQUARE_BRACKET_CLOSING, Literal: "]"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.COLON, Literal: ":"},
		{Type: token.IDENT, Literal: "Response"},
		// {
		{Type: token.CURLY_OPENING, Literal: "{"},
		// $twig = $this['twig'];
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "twig"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "this"},
		{Type: token.SQUARE_BRACKET_OPENING, Literal: "["},
		{Type: token.STRING, Literal: "twig"},
		{Type: token.SQUARE_BRACKET_CLOSING, Literal: "]"},
		{Type: token.SEMICOLON, Literal: ";"},
		// if ($response instanceof StreamedResponse) {
		{Type: token.IF, Literal: "if"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "response"},
		{Type: token.INSTANCEOF, Literal: "instanceof"},
		{Type: token.IDENT, Literal: "StreamedResponse"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		// $response->setCallback(function () use ($twig, $view) {
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "response"},
		{Type: token.OBJECT_OPERATOR, Literal: "->"},
		{Type: token.IDENT, Literal: "setCallback"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.FUNCTION, Literal: "function"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.USE, Literal: "use"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "twig"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "view"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		// $twig->display($view);
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "twig"},
		{Type: token.OBJECT_OPERATOR, Literal: "->"},
		{Type: token.IDENT, Literal: "display"},
		{Type: token.PARENTHESIS_OPENING, Literal: "("},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "view"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		// });
		{Type: token.CURLY_CLOSING, Literal: "}"},
		{Type: token.PARENTHESIS_CLOSING, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		// } else {
		{Type: token.CURLY_CLOSING, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.CURLY_OPENING, Literal: "{"},
		// $response;
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "response"},
		{Type: token.SEMICOLON, Literal: ";"},
		// }
		{Type: token.CURLY_CLOSING, Literal: "}"},
		// return $response;
		{Type: token.RETURN, Literal: "return"},
		{Type: token.VAR, Literal: "$"},
		{Type: token.IDENT, Literal: "response"},
		{Type: token.SEMICOLON, Literal: ";"},
		// }
		{Type: token.CURLY_CLOSING, Literal: "}"},
		// End
	}

	run(t, input, expectations)
}

func TestScanner_Next_FileStart(t *testing.T) {
	var input = `
/*
 * This file is part of the Silex framework.
 *
 * (c) Fabien Potencier <fabien@symfony.com>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Silex;

use Symfony\Component\Debug\ExceptionHandler as DebugExceptionHandler;
use Symfony\Component\Debug\Exception\FlattenException;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Event\GetResponseForExceptionEvent;
use Symfony\Component\HttpKernel\KernelEvents;

`
	expectations := []Token{
		{Type: token.COMMENT, Literal: `/*
 * This file is part of the Silex framework.
 *
 * (c) Fabien Potencier <fabien@symfony.com>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */`},
		// namespace Silex;
		{Type: token.NAMESPACE, Literal: "namespace"},
		{Type: token.IDENT, Literal: "Silex"},
		{Type: token.SEMICOLON, Literal: ";"},
		// use Symfony\Component\Debug\ExceptionHandler as DebugExceptionHandler;
		{Type: token.USE, Literal: "use"},
		{Type: token.IDENT, Literal: "Symfony"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Component"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Debug"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "ExceptionHandler"},
		{Type: token.AS, Literal: "as"},
		{Type: token.IDENT, Literal: "DebugExceptionHandler"},
		{Type: token.SEMICOLON, Literal: ";"},
		//use Symfony\Component\Debug\Exception\FlattenException;
		{Type: token.USE, Literal: "use"},
		{Type: token.IDENT, Literal: "Symfony"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Component"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Debug"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Exception"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "FlattenException"},
		{Type: token.SEMICOLON, Literal: ";"},
		//use Symfony\Component\EventDispatcher\EventSubscriberInterface;
		{Type: token.USE, Literal: "use"},
		{Type: token.IDENT, Literal: "Symfony"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Component"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "EventDispatcher"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "EventSubscriberInterface"},
		{Type: token.SEMICOLON, Literal: ";"},
		//use Symfony\Component\HttpFoundation\Response;
		{Type: token.USE, Literal: "use"},
		{Type: token.IDENT, Literal: "Symfony"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Component"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "HttpFoundation"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Response"},
		{Type: token.SEMICOLON, Literal: ";"},
		//use Symfony\Component\HttpKernel\Event\GetResponseForExceptionEvent;
		{Type: token.USE, Literal: "use"},
		{Type: token.IDENT, Literal: "Symfony"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Component"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "HttpKernel"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Event"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "GetResponseForExceptionEvent"},
		{Type: token.SEMICOLON, Literal: ";"},
		//use Symfony\Component\HttpKernel\KernelEvents;
		{Type: token.USE, Literal: "use"},
		{Type: token.IDENT, Literal: "Symfony"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "Component"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "HttpKernel"},
		{Type: token.BACKSLASH, Literal: "\\"},
		{Type: token.IDENT, Literal: "KernelEvents"},
		{Type: token.SEMICOLON, Literal: ";"},
	}

	run(t, input, expectations)
}
