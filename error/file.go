package error

import (
	"fmt"
	"strings"
)

// NewFormatter ...
func NewFormatter(filename string, data []rune) *Formatter { return &Formatter{filename: filename, data: data} }

type Formatter struct {
	filename string
	data     []rune
}

func (p *Formatter) Format(message string, pos int) error {
	format := `ParseError: %s in %s:%d:%d

^^^^^^^^^^^^^^^^^^^^^^
%s
%s
`
	lineNum, linePos, line := p.getLine(pos)
	pointer := strings.Repeat(" ", linePos)
	pointer += "^"

	return fmt.Errorf(format, message, p.filename, lineNum, linePos, line, pointer)
}

func (p *Formatter) getLine(pos int) (int, int, string) {
	line := 0
	firstCharPos := 0
	for i := 0; i < pos; i++ {
		if p.data[i] == '\n' {
			firstCharPos = i + 1
			line++
		}
	}
	// no \n
	if line == 0 {
		return 1, pos, string(p.data)
	}

	end := firstCharPos
	for {
		if p.data[end] == '\n' {
			break
		}
		end++
	}

	return line, pos - firstCharPos, string(p.data[firstCharPos:end])
}
