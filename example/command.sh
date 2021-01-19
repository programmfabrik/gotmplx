#!/bin/sh

# stdin as csv
cat example/sample1.csv | ./gotmplx --var some=something --var moar=data --csv one=- example/sample1.txt example/partial_tpl_1.txt

# stdin as template
echo "{{ template \"partial_1\" }}" | ./gotmplx --var some=something --var moar=data --csv one=example/sample1.csv - example/partial_tpl_1.txt

# eval as template
./gotmplx --var some=something --var moar=data --csv one=example/sample1.csv --eval "{{ template \"partial_1\" }}" example/partial_tpl_1.txt