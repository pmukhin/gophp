package interpret

import (
	"os"
	"io/ioutil"
	"github.com/pmukhin/gophp/scanner"
	"github.com/pmukhin/gophp/parser"
	"github.com/pmukhin/gophp/eval"
	"github.com/pmukhin/gophp/object"
)

func Main(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return runFile(f)
}

func runFile(file *os.File) error {
	sourceCode, err := readSourceCode(file)
	if err != nil {
		return err
	}
	// remove fucking shebang
	if string(sourceCode[:5]) == "<?php" {
		sourceCode = sourceCode[5:]
	}
	p := parser.New(scanner.New(sourceCode))
	code, e := p.Parse()
	// parse error
	if e != nil {
		return e
	}

	ctx := object.NewContext(nil, object.InternalFunctionTable)
	_, e = eval.Eval(code, ctx)

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
