#!/bin/bash
FILE=${1:?file}
JSON=/tmp/json$$
if file -i $FILE | grep shellscript >/dev/null
then echo '{"value":{"code":'$(cat $FILE | jq -R -s .)'}}' >$JSON
else echo '{"value":{"binary":true,"code":"'$(base64 -w 0 $FILE)'"}}' >$JSON
fi
curl -XPOST http://localhost:${PORT:-8080}/init -d @$JSON 2>/dev/null
rm $JSON
