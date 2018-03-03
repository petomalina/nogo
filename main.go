package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/gelidus/nogo/utils"
	"github.com/gelidus/nogo/visitors/example"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var (
	visitors = map[string]utils.Visitor{
		"example": &example.Visitor{},
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

			pkgs, err := parser.ParseDir(fset, fileName, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			for name, node := range pkgs {
				astFileMap[name] = node
			}
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
		for fileName, ast := range astFileMap {
			for _, name := range patterns {
				log.Info("Running visitor: ", name, " for file: ", fileName)

				errs := visitors[name].Run(ast)
				if errs != nil && len(errs) > 0 {
					log.Warn("Errors were reported by the visitor: ", errs)
				}
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
