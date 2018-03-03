package pointers

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
	return "pointers"
}

type Walker struct {
	errs []error
}

func (w *Walker) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.StarExpr:
		switch n.X.(type) {
		case *ast.InterfaceType:
			w.errs = append(w.errs, errors.New("Don't use pointers to interfaces"))
		case *ast.ArrayType:
			w.errs = append(w.errs, errors.New("Don't use pointers to slices"))
		case *ast.MapType:
			w.errs = append(w.errs, errors.New("Don't use pointers to maps"))
		case *ast.ChanType:
			w.errs = append(w.errs, errors.New("Don't use pointers to channels"))
		}
	default:
		return w
	}

	return nil
}
