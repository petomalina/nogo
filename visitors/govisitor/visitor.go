package govisitor

import (
	"errors"
	"go/ast"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ast/astutil"
)

type Visitor struct{}

type Walker struct {
	errs []error
	f    *ast.File
}

func (v *Visitor) Run(f ast.Node) []error {
	log.Debug("Running go-call antipattern visitor")

	walker := &Walker{
		errs: []error{},
		f:    f.(*ast.File),
	}

	ast.Walk(v, f)

	return walker.errs
}

func (v *Visitor) Name() string {
	return "go"
}

func (v *Visitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.GoStmt:
		w.detectScope(n)
	default:
		return v
	}

	return nil
}

func (w *Walker) detectScope(n ast.Node) {
	nodes, _ := astutil.PathEnclosingInterval(w.f, n.Pos(), n.End())

	for _, node := range nodes {
		switch x := node.(type) {
		case *ast.FuncDecl:
			if ast.IsExported(x.Name.Name) {
				w.errs = append(w.errs, errors.New("go statement used in exported function"))
			}
		}
	}
}
