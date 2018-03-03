package utils

import "go/ast"

type Visitor interface {
	Run(ast.Node) []error
}
