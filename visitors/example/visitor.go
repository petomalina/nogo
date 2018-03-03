package example

import (
	"errors"
	"go/ast"

	"github.com/gelidus/nogo/utils"

	log "github.com/sirupsen/logrus"
)

var (
	errs = []error{}
)

type Visitor struct{}

func (v *Visitor) Run(f ast.Node) []error {
	log.Debug("Running example pattern visitor")

	ast.Walk(utils.VisitorFunc(walker), f)

	return errs
}

func walker(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl:
		if n.Name.Name == "main" {
			errs = append(errs, errors.New("Found too many mains (at least one)"))
		}
	default:
		return utils.VisitorFunc(walker)
	}

	return nil
}
