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
	l := sync.Mutex{}

	for name, node := range asts {
		wg.Add(1)
		go func(name string, node ast.Node) {
			log.Debug("Running visitor: '", v.Name(), "' on: '", name, "'")
			errs := v.Run(node)
			if len(errs) > 0 {
				l.Lock()
				results[name] = errs
				l.Unlock()
			}
			wg.Done()
		}(name, node)
	}

	wg.Wait()

	return results
}

func RunVisitorsInParallel(asts map[string]ast.Node, vs []Visitor) map[string]map[string][]error {
	results := map[string]map[string][]error{}

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for _, v := range vs {
		wg.Add(1)
		go func(v Visitor) {
			res := RunVisitorInParallel(asts, v)
			if len(res) > 0 {
				l.Lock()
				results[v.Name()] = res
				l.Unlock()
			}
			wg.Done()
		}(v)
	}

	wg.Wait()

	return results
}
