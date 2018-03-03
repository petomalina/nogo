package utils

import (
	"go/ast"
	"sync"
)

type Visitor interface {
	Run(ast.Node) []error
}

func RunInParallel(asts map[string]ast.Node, v Visitor) map[string][]error {
	results := map[string][]error{}

	wg := sync.WaitGroup{}
	wg.Add(len(asts))

	for name, node := range asts {
		go func(name string) {
			errs := v.Run(node)
			if len(errs) > 0 {
				results[name] = errs
			}
			wg.Done()
		}(name)
	}

	wg.Wait()

	return results
}
