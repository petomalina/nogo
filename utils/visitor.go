package utils

import (
	"go/ast"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Visitor interface {
	Run(ast.Node) []error
	Name() string
}

func RunVisitorInParallel(asts map[string]ast.Node, v Visitor) map[string][]error {
	results := map[string][]error{}

	wg := sync.WaitGroup{}
	wg.Add(len(asts))

	for name, node := range asts {
		go func(name string) {
			log.Debug("Running visitor: '", v.Name(), "' on: '", name, "'")
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

func RunVisitorsInParallel(asts map[string]ast.Node, vs []Visitor) map[string]map[string][]error {
	results := map[string]map[string][]error{}

	wg := sync.WaitGroup{}
	wg.Add(len(vs))

	for _, v := range vs {
		go func(v Visitor) {
			res := RunVisitorInParallel(asts, v)
			if len(res) > 0 {
				results[v.Name()] = res
			}
			wg.Done()
		}(v)
	}

	wg.Wait()

	return results
}
