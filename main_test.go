package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var inputSamplePartialTmpl = `
{{ define "partial_1" }}
Environment
{{- range $k, $v := .Env }}
    Env {{ $k }} => {{ $v }}
{{- end }}
{{- end }}
`

var inputSampleTmpl = `
Sample template START
{{- template "partial_1" . }}
Variables
{{- range $k, $v := .Var }}
    Var {{ $k }} => {{ $v }}
{{- end }}
CSV data
{{- range $k, $v := .CSV }}
    CSV {{ $k }}
    {{- range $i, $v2 := $v }}
        Row {{ $i }}
        {{- range $k3, $v3 := $v2 }}
            Field {{ $k3 }} => {{ $v3 }}
        {{- end }}
    {{- end }}
{{- end }}
Sample template END
`

var outputSampleText = `
Sample template START
Environment
    Env env1 => val1
Variables
    Var moar => more=data
    Var some => something
CSV data
    CSV one
        Row 0
            Field active => true
            Field db => fylr-uni-muenster
            Field instance => uni-muenster
            Field url => https://uni-muenster.fylr.dev
        Row 1
            Field active => true
            Field db => fylr-uni-muenster-archiv
            Field instance => uni-muenster-archiv
            Field url => https://uni-muenster-archiv.fylr.dev
        Row 2
            Field active => true
            Field db => fylr-census
            Field instance => census
            Field url => https://census.fylr.dev
        Row 3
            Field active => true
            Field db => fylr-demo
            Field instance => demo
            Field url => https://demo.fylr.dev
        Row 4
            Field active => false
            Field db => fylr-gta-collections
            Field instance => gta-collections
            Field url => https://gta-collections.fylr.dev
        Row 5
            Field active => false
            Field db => fylr-unib-heidelberg
            Field instance => unib-heidelberg
            Field url => https://unib-heidelberg.fylr.dev
        Row 6
            Field active => false
            Field db => fylr-leon-testing
            Field instance => leon-testing
            Field url => https://leon-testing.fylr.dev
Sample template END
`

var inputSampleCSV = `
active,instance,db,url
bool,string,string,string
TRUE,uni-muenster,fylr-uni-muenster,https://uni-muenster.fylr.dev
TRUE,uni-muenster-archiv,fylr-uni-muenster-archiv,https://uni-muenster-archiv.fylr.dev
TRUE,census,fylr-census,https://census.fylr.dev
TRUE,demo,fylr-demo,https://demo.fylr.dev
FALSE,gta-collections,fylr-gta-collections,https://gta-collections.fylr.dev
FALSE,unib-heidelberg,fylr-unib-heidelberg,https://unib-heidelberg.fylr.dev
FALSE,leon-testing,fylr-leon-testing,https://leon-testing.fylr.dev
`

var tFolder string

func TestMain(t *testing.T) {
	tFolder = t.TempDir()
	t.Run("one", testTemplateEval)
	t.Run("two", testTemplateEvalCSVStdIn)
	t.Run("three", testTemplateDelimiters)
}

func testTemplateEval(t *testing.T) {

	csvFilePath := filepath.Join(tFolder, "test_sample.csv")
	err := ioutil.WriteFile(csvFilePath, []byte(inputSampleCSV), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tpl1FilePath := filepath.Join(tFolder, "partial1.txt")
	err = ioutil.WriteFile(tpl1FilePath, []byte(inputSamplePartialTmpl), 0644)
	if err != nil {
		t.Fatal(err)
	}

	inBuf := bytes.NewBufferString(inputSampleCSV)
	rootCmd.SetIn(inBuf)
	outBuf := bytes.NewBufferString("")
	rootCmd.SetOut(outBuf)
	rootCmd.SetArgs([]string{
		"--eval", inputSampleTmpl,
		"--var", "some=something",
		"--var", "moar=more=data",
		"--csv", fmt.Sprintf("one=%s", csvFilePath),
		tpl1FilePath,
	})

	// Clear / set env
	os.Clearenv()
	err = os.Setenv("env1", "val1")
	if err != nil {
		t.Fatal(err)
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(outBuf)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != outputSampleText {
		t.Fatalf("Expected:\n%s\nGot:\n%s", outputSampleText, string(out))
	}
}

func testTemplateEvalCSVStdIn(t *testing.T) {

	tpl1FilePath := filepath.Join(tFolder, "partial1.txt")
	err := ioutil.WriteFile(tpl1FilePath, []byte(inputSamplePartialTmpl), 0644)
	if err != nil {
		t.Fatal(err)
	}

	inBuf := bytes.NewBufferString(inputSampleCSV)
	rootCmd.SetIn(inBuf)
	outBuf := bytes.NewBufferString("")
	rootCmd.SetOut(outBuf)
	rootCmd.SetArgs([]string{
		"--eval", inputSampleTmpl,
		"--var", "some=something",
		"--var", "moar=more=data",
		"--csv", "one=-",
		tpl1FilePath,
	})

	// Clear / set env
	os.Clearenv()
	err = os.Setenv("env1", "val1")
	if err != nil {
		t.Fatal(err)
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(outBuf)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != outputSampleText {
		t.Fatalf("Expected:\n%s\nGot:\n%s", outputSampleText, string(out))
	}
}

func testTemplateDelimiters(t *testing.T) {

	inputTmpl := `
Environment
%%- range $k, $v := .Env %%
    Env %% $k %% => %% $v %%
%%- end %%
Variables
%%- range $k, $v := .Var %%
    Var %% $k %% => %% $v %%
%%- end %%
`

	outputText := `
Environment
    Env env1 => val1
    Env env2 => val2
Variables
    Var moar => more=data
    Var some => something
`

	csvs = nil
	outBuf := bytes.NewBufferString("")
	rootCmd.SetOut(outBuf)
	rootCmd.SetArgs([]string{
		"--template-delim-left", "%%",
		"--template-delim-right", "%%",
		"--eval", inputTmpl,
		"--var", "some=something",
		"--var", "moar=more=data",
	})

	os.Clearenv()
	err := os.Setenv("env1", "val1")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("env2", "val2")
	if err != nil {
		t.Fatal(err)
	}
	err = rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(outBuf)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != outputText {
		t.Fatalf("\nGot:\n%s\nExpected:\n%s", string(out), outputText)
	}
}
