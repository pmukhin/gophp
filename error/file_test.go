package error

import (
	"testing"
)

func TestPrinter_Make(t *testing.T) {
	data := []rune(`
<?php

$temporary = new final class extends Exception {}
try {
	throw $temporary	
} catch (Exception $exception) {
	println("can not do operation")
}
`)

	printer := Formatter{filename: "src/ImportantClass.php", data: data}
	err := printer.Format("unexpected token extends", 36)

	expected := `ParseError: unexpected token extends in src/ImportantClass.php:2:29

^^^^^^^^^^^^^^^^^^^^^^
$temporary = new final class extends Exception {}
                             ^
`
	if err.Error() != expected {
		t.Errorf("printer.Format() = [%s], expected [%s]", err, expected)
	}
}
