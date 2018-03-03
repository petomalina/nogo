package utils

import (
	"go/ast"
	"go/token"
	"sync"

	log "github.com/sirupsen/logrus"
)

type SourceFile struct {
	Node ast.Node
	Fset *token.FileSet
}

type Visitor interface {
	Run(ast.Node, *token.FileSet) []error
	Name() string
}

func RunVisitorInParallel(asts map[string]SourceFile, v Visitor) map[string][]error {
	results := map[string][]error{}

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for name, node := range asts {
		wg.Add(1)
		go func(name string, source SourceFile) {
			log.Debug("Running visitor: '", v.Name(), "' on: '", name, "'")
			errs := v.Run(source.Node, source.Fset)
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

func RunVisitorsInParallel(asts map[string]SourceFile, vs []Visitor) map[string]map[string][]error {
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
