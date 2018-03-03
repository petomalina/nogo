package example

import (
	"errors"
	"go/ast"

	log "github.com/sirupsen/logrus"
)

type Visitor struct{}

func (v *Visitor) Run(f ast.Node) []error {
	log.Debug("Running example pattern visitor")

	walker := &Walker{
		errs: []error{},
	}

	ast.Walk(walker, f)

	return walker.errs
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
