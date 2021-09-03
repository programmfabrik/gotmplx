// +build go1.15

package main

// Copyright Programmfabrik GmbH
// All Rights Reserved

// The gotmplx command wires up variables into a go template and renders it as output

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"github.com/programmfabrik/go-csvx"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none-commit"
	date    = "2006-01-02 15:04:05Z07:00"
	builtBy = "unknown"

	rootCmd = &cobra.Command{
		Use:     "gotmplx TEMPLATE [PARTIAL_TEMPLATE]*",
		Short:   "gotmplx: Command line tool to render a go template",
		Version: fmt.Sprintf("%s %s %v %s", version, commit, date, builtBy),
		Example: "gotmplx --var some=something --csv one=example/sample1.csv example/sample1.txt example/partial_tpl_1.txt",
		Run:     render,
		PreRun:  parseVariables,
	}
	showVersion          bool
	vars                 []string
	csvs                 []string
	eval                 string
	templateEnvVariables map[string]interface{}
	templateVariables    map[string]interface{}
	templateCSVVariables map[string][]map[string]interface{}
)

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Version of gotmplx")
	rootCmd.Flags().StringArrayVarP(&vars, "var", "", []string{}, "Parse and use variable in template (--var myvar=value)")
	rootCmd.Flags().StringArrayVarP(&csvs, "csv", "", []string{}, "Parse and use CSV file rows in template (--csv key=file)")
	rootCmd.Flags().StringVarP(&eval, "eval", "e", "", "Parse this text instead of file argument (--eval \"{{ .Var.myvar }}\"")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

func parseVariables(cmd *cobra.Command, args []string) {

	envStr := os.Environ()
	templateEnvVariables = make(map[string]interface{})
	for _, v := range envStr {
		key, value, err := splitVarParam(v)
		if err != nil {
			fmt.Fprint(cmd.OutOrStderr(), err)
			os.Exit(1)
		}
		templateEnvVariables[key] = value
	}

	templateCSVVariables = make(map[string][]map[string]interface{})
	for _, v := range csvs {
		var (
			csvBytes []byte
			err      error
		)
		key, csvFileName, err := splitVarParam(v)
		if err != nil {
			fmt.Fprint(cmd.OutOrStderr(), err)
			os.Exit(1)
		}
		if csvFileName == "-" {
			csvBytes, err = ioutil.ReadAll(cmd.InOrStdin())
			if err != nil {
				fmt.Fprint(cmd.OutOrStderr(), errors.Wrap(err, "Could not read stdin"))
				os.Exit(1)
			}
		} else {
			csvBytes, err = ioutil.ReadFile(csvFileName)
			if err != nil {
				fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not read CSV file %s", key))
				os.Exit(1)
			}
		}

		csvp := csvx.CSVParser{
			Comma:            ',',
			Comment:          '#',
			TrimLeadingSpace: true,
			SkipEmptyColumns: true,
		}

		templateCSVVariables[key], err = csvp.Untyped(csvBytes)
		if err != nil {
			fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not parse CSV file %s", csvFileName))
			os.Exit(1)
		}
	}

	templateVariables = make(map[string]interface{})
	for _, v := range vars {
		key, value, err := splitVarParam(v)
		if err != nil {
			fmt.Fprint(cmd.OutOrStderr(), err)
			os.Exit(1)
		}
		templateVariables[key] = value
	}
}

func render(cmd *cobra.Command, args []string) {

	var (
		tpl *template.Template
		err error
	)

	if len(args) == 0 && eval == "" {
		cmd.Usage()
		os.Exit(1)
	}

	if eval != "" {
		tpl = template.New("eval").Funcs(sprig.FuncMap())
		tpl, err = tpl.Parse(eval)
		if err != nil {
			fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not parse inline template `%s`", eval))
			os.Exit(1)
		}
	}

	for _, arg := range args {
		var t *template.Template
		if arg == "-" {
			if tpl == nil {
				tpl = template.New("stdin").Funcs(sprig.FuncMap())
				t = tpl
			} else {
				t = tpl.New("stdin")
			}
			stdInBytes, err := ioutil.ReadAll(cmd.InOrStdin())
			if err != nil {
				fmt.Fprint(cmd.OutOrStderr(), errors.Wrap(err, "Could not read stdin"))
				os.Exit(1)
			}
			_, err = t.Parse(string(stdInBytes))
			if err != nil {
				fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not parse template from stdin: %s", string(stdInBytes)))
				os.Exit(1)
			}
		} else {
			if tpl == nil {
				tpl = template.New(filepath.Base(arg)).Funcs(sprig.FuncMap())
				t = tpl
			} else {
				t = tpl.New(filepath.Base(arg))
			}
			_, err = t.ParseFiles(arg)
			if err != nil {
				fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not parse template file %s", arg))
				os.Exit(1)
			}
		}
	}

	data := map[string]interface{}{
		"Env": templateEnvVariables,
		"Var": templateVariables,
		"CSV": templateCSVVariables,
	}

	err = tpl.Execute(cmd.OutOrStdout(), data)
	if err != nil {
		fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not execute template file %s with data %v", tpl.Name(), data))
		os.Exit(1)
	}
}

func splitVarParam(param string) (string, string, error) {
	parts := strings.Split(param, "=")
	if len(parts) < 2 {
		return "", "", errors.Errorf("Flag arguments should be `name=value`, given %s", param)
	}
	return parts[0], strings.Join(parts[1:], "="), nil
}
