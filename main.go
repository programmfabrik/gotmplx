//go:build go1.15
// +build go1.15

package main

// Copyright Programmfabrik GmbH
// All Rights Reserved

// The gotmplx command wires up variables into a go template and renders it as output

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	ttemplate "text/template"

	"github.com/Masterminds/sprig"
	"github.com/programmfabrik/go-csvx"
	"github.com/spf13/cobra"
	"github.com/yudai/pp"
	"gopkg.in/yaml.v3"
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
	vars, csvs, ymls, jsons                                            []string
	eval, output                                                       string
	templateEnvVariables, templateVariables, templateYML, templateJSON map[string]any
	templateCSVVariables                                               map[string][]map[string]any

	dump, html bool

	templateDelimLeft, templateDelimRight string
	stdinBytes                            []byte
)

func init() {
	rootCmd.Flags().StringArrayVarP(&vars, "var", "", []string{}, "Parse and use variable in template (--var myvar=value)")
	rootCmd.Flags().StringArrayVarP(&csvs, "csv", "", []string{}, "Parse and use CSV file rows in template (--csv key=file)")
	rootCmd.Flags().StringArrayVarP(&ymls, "yml", "", []string{}, "Parse and use YML file in template (--yml key=file)")
	rootCmd.Flags().StringArrayVarP(&jsons, "json", "", []string{}, "Parse and use JSON file in template (--json key=file)")

	rootCmd.Flags().BoolVarP(&dump, "dump", "d", false, "Pretty print data passed to template to stdout")
	rootCmd.Flags().BoolVarP(&html, "html", "", false, "Render template as HTML (default is TEXT)")
	rootCmd.Flags().StringVarP(&eval, "eval", "e", "", "Parse this text instead of file argument (--eval \"{{ .Var.myvar }}\"")
	rootCmd.Flags().StringVarP(&output, "output", "o", "-", `Send output to file. Use for "-" (default)`)

	rootCmd.Flags().StringVarP(&templateDelimLeft, "template-delim-left", "l", "", "Use this as left delimiter in go templates")
	rootCmd.Flags().StringVarP(&templateDelimRight, "template-delim-right", "r", "", "Use this as right delimiter in go templates")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

// stdin reads all bytes from stdin. if called more than once, the stdin bytes
// from the first call are returns
func stdin() []byte {
	if stdinBytes == nil {
		var err error
		stdinBytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Could not read stdin: %s", err.Error())
		}
	}
	return stdinBytes
}

func parseVariables(cmd *cobra.Command, args []string) {

	envStr := os.Environ()
	templateEnvVariables = map[string]any{}
	for _, v := range envStr {
		key, value, ok := strings.Cut(v, "=")
		if ok {
			templateEnvVariables[key] = value
		}
	}

	templateYML = map[string]any{}
	for _, y := range ymls {
		var (
			ymlBytes []byte
			err      error
		)
		key, ymlFileName, ok := strings.Cut(y, "=")
		if !ok {
			log.Fatalf("Unable to split yml file %q", y)
		}
		if ymlFileName == "-" {
			ymlBytes = stdin()
		} else {
			ymlBytes, err = ioutil.ReadFile(ymlFileName)
			if err != nil {
				log.Fatalf("Could not read yml file %q: %s", key, err.Error())
			}
		}
		var d any
		err = yaml.Unmarshal(ymlBytes, &d)
		if err != nil {
			log.Fatalf("Could not parse yml file %q: %s", ymlFileName, err.Error())
		}
		templateYML[key] = d
	}

	templateJSON = map[string]any{}
	for _, j := range jsons {
		var (
			jsonBytes []byte
			err       error
		)
		key, jsonFileName, ok := strings.Cut(j, "=")
		if !ok {
			log.Fatalf("Unable to split json file %q", j)
		}
		if jsonFileName == "-" {
			jsonBytes = stdin()
		} else {
			jsonBytes, err = ioutil.ReadFile(jsonFileName)
			if err != nil {
				log.Fatalf("Could not read json file %q: %s", key, err.Error())
			}
		}
		var d any
		err = json.Unmarshal(jsonBytes, &d)
		if err != nil {
			log.Fatalf("Could not parse json file %q: %s", jsonFileName, err.Error())
		}
		templateJSON[key] = d
	}

	templateCSVVariables = map[string][]map[string]any{}
	for _, v := range csvs {
		var (
			csvBytes []byte
			err      error
		)
		key, csvFileName, ok := strings.Cut(v, "=")
		if !ok {
			log.Fatalf("Unable to split csv file %q", v)
		}
		if csvFileName == "-" {
			csvBytes = stdin()
		} else {
			csvBytes, err = ioutil.ReadFile(csvFileName)
			if err != nil {
				log.Fatalf("Could not read csv file %q: %s", key, err.Error())
			}
		}

		csvp := csvx.CSVParser{
			Comma:            ',',
			Comment:          '#',
			TrimLeadingSpace: true,
			SkipEmptyColumns: true,
		}

		templateCSVVariables[key], err = csvp.Typed(csvBytes)
		if err != nil {
			println(len(csvBytes))
			log.Fatalf("Could not parse csv file %q: %s", csvFileName, err.Error())
		}
	}

	templateVariables = map[string]any{}
	for _, v := range vars {
		key, value, ok := strings.Cut(v, "=")
		if !ok {
			log.Fatalf("Unable to split --var %q", v)
		}
		templateVariables[key] = value
	}
}

func render(cmd *cobra.Command, args []string) {

	if len(args) == 0 && eval == "" {
		cmd.Usage()
		os.Exit(1)
	}

	tplBytes := []byte{}
	if eval != "" {
		tplBytes = append(tplBytes, eval...)
	}

	for _, arg := range args {
		if arg == "-" {
			tplBytes = append(tplBytes, stdin()...)
		} else {
			fBytes, err := os.ReadFile(arg)
			if err != nil {
				log.Fatalf("Unable to read %q: %s", arg, err.Error())
			}
			tplBytes = append(tplBytes, fBytes...)
		}
	}

	data := map[string]any{
		"env":  templateEnvVariables,
		"var":  templateVariables,
		"csv":  templateCSVVariables,
		"yml":  templateYML,
		"json": templateJSON,
	}

	if dump {
		pp.Println(data)
		os.Exit(0)
	}

	var err error

	// HTML rendering
	if html {
		tpl := template.
			New("tmpl").
			Funcs(sprig.FuncMap()).
			Delims(templateDelimLeft, templateDelimRight)
		if len(tplBytes) > 0 {
			tpl, err = tpl.Parse(string(tplBytes))
			if err != nil {
				log.Fatalf("Could not parse template: %s", err.Error())
			}
		}
		err = tpl.Execute(os.Stdout, data)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		tpl := ttemplate.
			New("tmpl").
			Funcs(sprig.FuncMap()).
			Delims(templateDelimLeft, templateDelimRight)
		if len(tplBytes) > 0 {
			tpl, err = tpl.Parse(string(tplBytes))
			if err != nil {
				log.Fatalf("Could not parse template: %s", err.Error())
			}
		}

		var out io.Writer
		if output != "-" {
			out, err = os.Create(output)
			if err != nil {
				log.Fatalf("Unable to open %q for output: %s", output, err.Error())
			}
		} else {
			out = os.Stdout
		}

		err = tpl.Execute(out, data)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
