# gotmplx

* **.Env**
* Cmdline --var as **.Var**

```bash
# gotmplx --var MyVar1=value2 --var "MyVar2=value3" *.tmpl.txt
# gotmplx --var "MyVar1=value2" --eval "Test: {{ .Var.MyVar1 }}"
```

* Output to STDOUT

## --var

The `--var test=val2` will be available in the template as **.Var.Test**.

## --csv Name:File

The `--csv "nr1=test.csv"` will be available in the template as **.CSV.Nr1**.

You can use *-* as filename.
