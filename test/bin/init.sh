#!/bin/bash
FILE=${1:?file}
JSON=${2:-/tmp/json$$}
if file -i $FILE | grep text/ >/dev/null
then echo '{"value":{"main":"main","code":'$(cat $FILE | jq -R -s .)'}}' >/tmp/json$$
else echo '{"value":{"binary":true,"code":"'$(base64 -w 0 $FILE)'"}}' >/tmp/json$$
fi
#cat $JSON | jq .
curl -H "Content-Type: application/json" -XPOST -w "%{http_code}\n" http://localhost:${PORT:-8080}/init -d @$JSON 2>/dev/null
rm /tmp/json$$ 2>/dev/null
