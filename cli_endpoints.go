package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/Masterminds/sprig"
	golib "github.com/programmfabrik/go-lib"
	"github.com/urfave/cli/v2"
)

func mainEntrypoint(c *cli.Context) error {
	// extract csv data
	csvData, err := readData(c.StringSlice("csv"), &csvReader{})
	if err != nil {
		return err
	}

	// extract json data
	jsonData, err := readData(c.StringSlice("json"), &jsonReader{})
	if err != nil {
		return err
	}

	templateVals := map[string]interface{}{
		"Var":  golib.MapValues(c.StringSlice("var"), ""),
		"Env":  golib.MapValues(os.Environ(), ""),
		"CSV":  csvData,
		"JSON": jsonData,
	}

	// read template data
	templateData := ""
	if c.String("eval") != "" {
		templateData = c.String("eval")
	} else {
		if !c.Args().Present() {
			return errors.New("unable to find template file or eval argument")
		}

		fBytes, err := ioutil.ReadFile(c.Args().First())
		if err != nil {
			return fmt.Errorf("unable to read bytes from file %q: %w", c.Args().First(), err)
		}

		templateData = string(fBytes)
	}

	tmplt, err := template.New("stdin").
		Funcs(sprig.FuncMap()).
		Delims(c.String("template-delim-left"), c.String("template-delim-right")).
		Parse(templateData)
	if err != nil {
		return err
	}

	err = tmplt.Execute(os.Stdout, templateVals)
	if err != nil {
		return err
	}

	return nil
}
