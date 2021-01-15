package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var inputSampleTmpl = `
Sample template START
Environment
{{- range $k, $v := .Env }}
    Env {{ $k }} => {{ $v }}
{{- end }}
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
    Var moar => data
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

func TestTemplate(t *testing.T) {

	// Temporarily flush csv file into disk
	tFolder := t.TempDir()
	csvFilePath := filepath.Join(tFolder, "test_sample.csv")
	err := ioutil.WriteFile(csvFilePath, []byte(inputSampleCSV), 0644)
	if err != nil {
		t.Fatal(err)
	}

	outBuf := bytes.NewBufferString("")
	rootCmd.SetOut(outBuf)
	rootCmd.SetArgs([]string{
		"--eval", inputSampleTmpl,
		"--var", "some=something",
		"--var", "moar=data",
		"--csv", fmt.Sprintf("one=%s", csvFilePath),
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
