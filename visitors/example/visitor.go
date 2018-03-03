package example

import (
	"errors"
	"go/ast"

	log "github.com/sirupsen/logrus"
)

var (
	errs = []error{}
)

type VisitorFunc func(n ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

type Visitor struct {
}

func (v *Visitor) Run(f *ast.File) []error {
	log.Debug("Running example pattern visitor")

	ast.Walk(VisitorFunc(walker), f)

	return errs
}

func walker(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl:
		if n.Name.Name == "main" {
			errs = append(errs, errors.New("Found too many mains (at least one)"))
		}
	default:
		return VisitorFunc(walker)
	}

	return nil
}
