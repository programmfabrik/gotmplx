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
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gotmplx TEMPLATE [PARTIAL_TEMPLATE]*",
		Short:   "gotmplx: Command line tool to render a go template",
		Version: "1.0",
		Run:     render,
		PreRun:  parseVariables,
	}
	vars                 []string
	csvs                 []string
	eval                 string
	templateEnvVariables map[string]interface{}
	templateVariables    map[string]interface{}
	templateCSVVariables map[string][]map[string]interface{}
)

func init() {
	rootCmd.Flags().StringArrayVarP(&vars, "var", "", []string{}, "Parse and use variable in template (--var key=value)")
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
		parts := strings.Split(v, "=")
		templateEnvVariables[parts[0]] = strings.Join(parts[1:], "=")
	}

	templateCSVVariables = make(map[string][]map[string]interface{})
	for _, v := range csvs {
		parts := strings.Split(v, "=")
		csvFileName := strings.Join(parts[1:], "=")
		var (
			csvBytes []byte
			err error
		)
		if csvFileName == "-" {
			csvBytes, err = ioutil.ReadAll(cmd.InOrStdin())
			if err != nil {
				fmt.Fprint(cmd.OutOrStderr(), errors.Wrap(err, "Could not read stdin"))
				os.Exit(1)
			}
		} else {
			csvBytes, err = ioutil.ReadFile(csvFileName)
			if err != nil {
				fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not read CSV file %s", parts[1]))
				os.Exit(1)
			}
		}

		vars, err := CSVToMap(csvBytes, ',')
		if err != nil {
			fmt.Fprint(cmd.OutOrStderr(), errors.Wrapf(err, "Could not parse CSV file %s", parts[1]))
			os.Exit(1)
		}
		templateCSVVariables[parts[0]] = vars
	}

	templateVariables = make(map[string]interface{})
	for _, v := range vars {
		parts := strings.Split(v, "=")
		templateVariables[parts[0]] = strings.Join(parts[1:], "=")
	}
}

func render(cmd *cobra.Command, args []string) {

	var (
		tpl *template.Template
		err error
	)

	if len(args) == 0 && eval == "" {
		fmt.Fprintln(cmd.OutOrStderr(), "No file argument neither eval string has been defined")
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
