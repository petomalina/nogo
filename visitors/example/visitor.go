package example

import (
	"errors"
	"go/ast"
)

type Visitor struct{}

func (v *Visitor) Run(f ast.Node) []error {
	walker := &Walker{
		errs: []error{},
	}

	ast.Walk(walker, f)

	return walker.errs
}

func (v *Visitor) Name() string {
	return "example"
}

type Walker struct {
	errs []error
}

func (w *Walker) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl:
		if n.Name.Name == "main" {
			w.errs = append(w.errs, errors.New("Found too many mains (at least one)"))
		}
	default:
		return w
	}

	return nil
}
