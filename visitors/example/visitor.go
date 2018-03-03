package example

import (
	"go/ast"

	log "github.com/sirupsen/logrus"
)

type Visitor struct {
}

func (v *Visitor) Run(f *ast.File) []error {
	log.Debug("Running example pattern visitor")

	return nil
}
