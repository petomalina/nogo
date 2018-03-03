package govisitor

import (
	"errors"
	"go/ast"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ast/astutil"
)

var errs []error

type Visitor struct {
	f *ast.File
}

func (v *Visitor) Run(f ast.Node) []error {
	log.Debug("Running go-call antipattern visitor")

	v.f = f.(*ast.File)

	ast.Walk(v, f)

	return errs
}

func (v *Visitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.GoStmt:
		v.detectScope(n)
	default:
		return v
	}

	return nil
}

func (v Visitor) detectScope(n ast.Node) {
	nodes, _ := astutil.PathEnclosingInterval(v.f, n.Pos(), n.End())

	for _, node := range nodes {
		switch x := node.(type) {
		case *ast.FuncDecl:
			if ast.IsExported(x.Name.Name) {
				errs = append(errs, errors.New("go statement used in exported function"))
			}
		}
	}
}
