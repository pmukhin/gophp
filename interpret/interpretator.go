package interpret

import (
	"os"
	"io/ioutil"
	"github.com/pmukhin/gophp/scanner"
	"github.com/pmukhin/gophp/parser"
	"github.com/pmukhin/gophp/eval"
	"github.com/pmukhin/gophp/object"
	error2 "github.com/pmukhin/gophp/error"
)

var evaluator eval.Evaluator

func Main(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	evaluator = eval.New()

	// create context
	ctx := object.NewContext(nil)
	object.RegisterGlobals(ctx)

	return runFile(f, ctx)
}

func runFile(file *os.File, ctx object.Context) error {
	sourceCode, err := readSourceCode(file)
	if err != nil {
		return err
	}
	// remove fucking shebang
	if string(sourceCode[:5]) == "<?php" {
		sourceCode = sourceCode[5:]
	}

	formatter := error2.NewFormatter(file.Name(), sourceCode)
	p := parser.New(scanner.New(sourceCode), formatter)
	code, e := p.Parse()
	// parse error
	if e != nil {
		return e
	}

	_, e = evaluator.Eval(code, ctx)

	return e
}

func readSourceCode(file *os.File) ([]rune, error) {
	srcBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	srcStr := string(srcBytes)

	return []rune(srcStr), nil
}
