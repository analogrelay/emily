package parser

import (
	"github.com/anurse/emily/compiler/ast"
	"github.com/anurse/emily/compiler/scanner"
)

type multiError []error

func (m multiError) Error() string {
	str := ""
	for i, v := range m {
		if i != 0 {
			str += ", "
		}
		str += v.Error()
	}
	return str
}

type parser struct {
	scanner *scanner.Scanner

	errors []error
}

func (p *parser) err() error {
	if len(p.errors) == 0 {
		return nil
	}
	return multiError(p.errors)
}

func (p *parser) parseFile() *ast.File {
	panic("not yet implemented")
}
