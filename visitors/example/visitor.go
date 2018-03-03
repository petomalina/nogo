package example

import (
	"errors"
	"go/ast"

	log "github.com/sirupsen/logrus"
)

var (
	errs = []error{}
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

func (v *Walker) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl:
		if n.Name.Name == "main" {
			errs = append(errs, errors.New("Found too many mains (at least one)"))
		}
	default:
		return v
	}

	return nil
}
