#!/bin/bash
FILE=${1:?compiled file}
file --mime-type "$FILE"  | sed -e 's/x-mach-binary/x-executable/'
echo '{"name":"Mike"}' | $FILE 3>/tmp/$$
cat /tmp/$$
rm /tmp/$$
