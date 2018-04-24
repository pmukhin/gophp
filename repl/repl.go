package repl

import (
	"os"
	"fmt"
	"bufio"
	"github.com/pmukhin/gophp/parser"
	"github.com/pmukhin/gophp/scanner"
	error2 "github.com/pmukhin/gophp/error"
	"github.com/pmukhin/gophp/eval"
	"github.com/pmukhin/gophp/object"
)

var (
	evaluator eval.Evaluator
	ctx       object.Context
)

func Main() {
	out := os.Stdout
	in := os.Stdin

	// create interpreter instance
	evaluator = eval.New()

	// create context
	ctx = object.NewContext(nil)
	object.RegisterGlobals(ctx)

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
	runeLine := []rune(line)
	formatter := error2.NewFormatter("<console>", runeLine)

	p := parser.New(scanner.New(runeLine), formatter)
	program, e := p.Parse()

	if e != nil {
		return object.Null, e
	}

	return evaluator.Eval(program, ctx)
}
