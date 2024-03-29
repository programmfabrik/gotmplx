# gotmplx

gotmplx implements the go Template Engine and makes it available as a standalone application. This application extends the engine with various functions and supports injection of environment variables, csv data and flag values.

## Getting started

Since this repository provides versioned versions with binaries, you can decide whether to download such a binary or install the application using the *go install* command.

Example of using the *go install* command:

```bash
go install github.com/programmfabrik/gotmplx
```

## Available flags and what they do

| Short| Flag                     | Type     | Description |
|------|--------------------------|----------|-------------|
|      | `--csv`                  | []string | Parse and use csv file rows in template (--csv key=file) |
|      | `--yml`                  | []string | Parse and use yml file in template (--yml key=file) |
|      | `--json`                 | []string | Parse and use json file in template (--json key=file) |
|      | `--dump`                 | bool     | Dump parsed data which is passed to template. |
| `-o` | `--output`               | string   | Send output to file, use "-" for stdout (default) |
|      | `--html`                 | bool     | Render template as HTML, defaults to text rendering. |
| `-e` | `--eval`                 | string   | Parse this text instead of file argument (--eval "{{ .Var.myvar }}" |
| `-h` | `--help`                 |          | Help for gotmplx |
|      | `--var`                  | []string | Parse and use variable in template (--var myvar=value) |
| `-l` | `-template-delim-left`   | string   | Use this string as go template left delimiter |
| `-r` | `-template-delim-right`  | string   | Use this string as go template right delimiter |

## Examples

### Inject cli flag value into the template

```bash
./gotmplx --var "key=value" --eval "My key value: {{.var.key}}"
```

Result:

```txt
My key value: value
```

### Inject a csv file into the template

Create the csv file:

```bash
cat <<EOF > test.csv
name,id
henk,10
EOF
```

```bash
./gotmplx --csv "mycsvfile=test.csv" --eval "My csv value: {{.csv.mycsvfile}}"
```

Result:

```txt
My csv value: [map[id:10 name:henk]]
```

### Using an environment variable

Since gotmplx injects the environment into the `.Env` key, you can access any environment variable from a template.

```bash
./gotmplx --eval "My shell: {{.env.SHELL}}"
```

Result:

```txt
My shell: /bin/bash
```

### Using a file as template source

Create the template file:

```bash
cat <<EOF > my.tmpl.txt
Name: {{.Var.Name}}
ID: {{.Var.ID}}
Shell: {{.Env.SHELL}}
EOF
```

Execute gotmplx:

```bash
./gotmplx --var "Name=dummy" --var "ID=10" my.tmpl.txt
```

Result:

```yaml
Name: dummy
ID: 10
Shell: /bin/bash
```

### Using custom template functions

This application supports custom template functions using the [sprig-github](https://github.com/Masterminds/sprig) template function library. Visit the Godoc for [sprig-godoc](https://pkg.go.dev/github.com/Masterminds/sprig) to see all the available functions that are available.
