#!/bin/bash
FILE=${1:?file}
JSON=/tmp/json$$
if file -i $FILE | grep text/ >/dev/null
then echo '{"value":{"main":"main","code":'$(cat $FILE | jq -R -s .)'}}' >$JSON
else echo '{"value":{"binary":true,"code":"'$(base64 -w 0 $FILE)'"}}' >$JSON
fi
#cat $JSON | jq .
curl -H "Content-Type: application/json" -XPOST -w "%{http_code}\n"  http://localhost:${PORT:-8080}/init -d @$JSON 2>/dev/null
rm $JSON
