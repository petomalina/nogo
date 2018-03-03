package utils

import "go/ast"

type Visitor interface {
	Run(ast.Node) []error
}

type VisitorFunc func(n ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor { return f(n) }
