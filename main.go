package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gelidus/nogo/utils"
	"github.com/gelidus/nogo/visitors/example"
	"github.com/gelidus/nogo/visitors/govisitor"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var (
	visitors = map[string]utils.Visitor{
		"example": &example.Visitor{},
		"go":      &govisitor.Visitor{},
	}
)

func main() {
	log.SetLevel(log.DebugLevel)

	app := cli.NewApp()
	app.Name = "nogo"
	app.Usage = ""

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "pattern, p",
			Usage: "Enables given pattern visitor",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() <= 0 {
			return errors.New("Missing file name")
		}

		// filename to parse
		fileName := c.Args()[0]
		info, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			return err
		}

		// fset for parse functions
		fset := token.NewFileSet()
		astFileMap := map[string]ast.Node{}

		if info.IsDir() {
			log.Debug("Executing on folder: ", fileName)

			filepath.Walk(fileName, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() != true && strings.HasSuffix(info.Name(), ".go") {
					f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
					if err != nil {
						return err
					}
					astFileMap[path] = f
				}
				return nil
			})
		} else {
			log.Debug("Executing on file: ", fileName)

			f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			astFileMap[fileName] = f
		}

		log.Debug("Parse successful for target: ", fileName)

		patterns := c.GlobalStringSlice("pattern")
		vs := []utils.Visitor{}
		for _, p := range patterns {
			vs = append(vs, visitors[p])
		}

		errs := utils.RunVisitorsInParallel(astFileMap, vs)

		if len(errs) > 0 {
			log.Warn("Errors were reported by the visitor: ", errs)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
