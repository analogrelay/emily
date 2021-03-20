package parser

import (
	"github.com/anurse/emily/compiler/ast"
	"github.com/anurse/emily/compiler/scanner"
)

func ParseFile(scanner *scanner.Scanner) (*ast.File, error) {
	var p parser
	p.scanner = scanner
	f := p.parseFile()
	return f, p.err()
}
