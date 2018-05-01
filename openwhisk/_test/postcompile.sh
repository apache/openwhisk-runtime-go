#!/bin/bash
FILE=${1:?compiled file}
file -i "$FILE"
echo '{"name":"Mike"}' | $FILE 3>/tmp/$$
cat /tmp/$$
rm /tmp/$$
