package main

import (
	"errors"
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
		log.Debug("Executing on file: ", fileName)

		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		log.Debug("File parsing successful: ", fileName)

		patterns := c.GlobalStringSlice("pattern")
		for _, name := range patterns {
			log.Info("Running visitor for pattern: ", name)

			errs := visitors[name].Run(f)
			if errs != nil && len(errs) > 0 {
				log.Warn("Errors were reported by the visitor: ", errs)
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
