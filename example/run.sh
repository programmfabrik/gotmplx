#!/bin/sh

# stdin as csv
CMD1="cat example/sample1.csv | ./gotmplx --var some=something --var moar=data --csv one=- example/sample1.txt example/partial_tpl_1.txt"
echo "# stdin as csv\n"
echo "$CMD1\n"
eval $CMD1
# cat example/sample1.csv | ./gotmplx --var some=something --var moar=data --csv one=- example/sample1.txt example/partial_tpl_1.txt

# stdin as template
CMD2="echo \"{{ template \\\"partial_1\\\" . }}\" | ./gotmplx --var some=something --var moar=data --csv one=example/sample1.csv - example/partial_tpl_1.txt"
echo "\n\n# stdin as template\n"
echo "$CMD2\n"
eval $CMD2
# echo "{{ template \"partial_1\" . }}" | ./gotmplx --var some=something --var moar=data --csv one=example/sample1.csv - example/partial_tpl_1.txt

# eval as template
CMD3="./gotmplx --var some=something --var moar=data --csv one=example/sample1.csv --eval \"{{ template \\\"partial_1\\\" . }}\" example/partial_tpl_1.txt"
echo "\n\n# eval as template\n"
echo "$CMD3\n"
eval $CMD3
# ./gotmplx --var some=something --var moar=data --csv one=example/sample1.csv --eval "{{ template \"partial_1\" . }}" example/partial_tpl_1.txt