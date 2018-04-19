package printer

import (
	"github.com/pmukhin/gophp/ast"
	"strings"
)

// Printer traverses ast and prints the resulting code
type Printer struct {
	blocks []string
}

// New is Printer constructor
func New() *Printer {
	p := new(Printer)
	p.blocks = make([]string, 0, 256)

	return p
}

func (p *Printer) Visit(_ ast.NodeType, node ast.Node) {
	p.blocks = append(p.blocks, node.String())
}

func (p *Printer) Result() string {
	return strings.Join(p.blocks, "\n")
}
