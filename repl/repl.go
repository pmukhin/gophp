package repl

import (
	"os"
	"fmt"
	"bufio"
	"github.com/pmukhin/gophp/parser"
	"github.com/pmukhin/gophp/scanner"
	"github.com/pmukhin/gophp/eval"
	"github.com/pmukhin/gophp/object"
)

var ctx eval.Context

func Main() {
	out := os.Stdout
	in := os.Stdin

	// create context
	ctx = eval.NewContext(nil, eval.InternalFunctionTable)

	reader := bufio.NewReader(in)
	for {
		fmt.Fprint(out, "php> ")
		str, e := reader.ReadString('\n')
		if e != nil {
			fmt.Fprintln(out, e.Error())
			continue
		}
		_, e = runLine(str)
		if e != nil {
			fmt.Fprintln(out, e.Error())
			continue
		}
	}
}

func runLine(line string) (object.Object, error) {
	p := parser.New(scanner.New([]rune(line)))
	program, e := p.Parse()

	if e != nil {
		return object.Null, e
	}

	return eval.Eval(program, ctx)
}
