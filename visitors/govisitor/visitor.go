package govisitor

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type Visitor struct{}

type Walker struct {
	errs []error
	f    *ast.File
	fset *token.FileSet
}

func (v *Visitor) Run(f ast.Node, fset *token.FileSet) []error {
	walker := &Walker{
		errs: []error{},
		f:    f.(*ast.File),
		fset: fset,
	}

	ast.Walk(walker, f)

	return walker.errs
}

func (v *Visitor) Name() string {
	return "go"
}

func (w *Walker) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.GoStmt:
		w.detectScope(n)
	default:
		return w
	}

	return nil
}

func (w *Walker) detectScope(n ast.Node) {
	nodes, _ := astutil.PathEnclosingInterval(w.f, n.Pos(), n.End())

	for _, node := range nodes {
		switch x := node.(type) {
		case *ast.FuncDecl:
			if ast.IsExported(x.Name.Name) {
				w.errs = append(w.errs, errors.New(fmt.Sprintf("%s: go statement used in exported function", w.fset.Position(x.Pos()))))
			}
		}
	}
}
