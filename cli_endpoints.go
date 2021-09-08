package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/Masterminds/sprig"
	"github.com/urfave/cli/v2"
)

func mainEntrypoint(c *cli.Context) error {
	// extract csv data
	varData, err := stringSliceToMap(c.StringSlice("var"))
	if err != nil {
		return err
	}

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

	// extract env data
	envMap, err := stringSliceToMap(os.Environ())
	if err != nil {
		return err
	}

	templateVals := map[string]interface{}{
		"Var":  varData,
		"CSV":  csvData,
		"JSON": jsonData,
		"Env":  envMap,
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

	tmplt, err := template.New("stdin").Funcs(sprig.FuncMap()).Parse(templateData)
	if err != nil {
		return err
	}

	err = tmplt.Execute(os.Stdout, templateVals)
	if err != nil {
		return err
	}

	return nil
}
