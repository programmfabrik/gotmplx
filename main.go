package main

import (
	"fmt"
	"log"
	"os"

	cli2 "github.com/urfave/cli/v2"
)

// Copyright Programmfabrik GmbH
// All Rights Reserved

var (
	version = "dev"
	commit  = "none-commit"
	date    = "2006-01-02 15:04:05Z07:00"
	builtBy = "unknown"
)

func main() {
	app := &cli2.App{
		Name:        "gotmplx",
		Usage:       `go run . --csv "key=example/sample.input.csv"  --var "key=value" --json "key=example/sample.input.json" --json 'key2={"hello": "world"}' example/sample.tmplt.json`,
		Version:     fmt.Sprintf("%s %s %v %s", version, commit, date, builtBy),
		Description: "",
		Copyright:   "Copyright @2021 Programmfabrik GmbH",
		Flags: []cli2.Flag{
			&cli2.StringSliceFlag{
				Name:  "var",
				Usage: "Parse and use variable in template (--var myvar=value)",
			},
			&cli2.StringSliceFlag{
				Name:    "csv",
				Usage:   "Parse and use CSV file rows in template (--csv key=file)",
				Aliases: []string{"c"},
			},
			&cli2.StringSliceFlag{
				Name:    "json",
				Usage:   "Parse and use JSON file rows in template (--json key=file)",
				Aliases: []string{"j"},
			},
			&cli2.StringFlag{
				Name:    "eval",
				Usage:   "Parse this text instead of file argument (--eval \"{{ .Var.myvar }}\"",
				Aliases: []string{"e"},
			},
			&cli2.StringFlag{
				Name:  "template-delim-left",
				Usage: "Use this as left delimiter in go templates",
			},
			&cli2.StringFlag{
				Name:  "template-delim-right",
				Usage: "Use this as right delimiter in go templates",
			},
		},
		Action: mainEntrypoint,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
