#!/bin/bash
FILE=${1:?file}
if file -i $FILE | grep shellscript
then echo '{"value":{"code":'$(cat $FILE | jq -R -s .)'}}' >$FILE.json
else echo '{"value":{"binary":true,"code":"'$(base64 -w 0 $FILE)'"}}' >$FILE.json
fi
curl -XPOST http://localhost:${PORT:-8080}/init -d @$FILE.json
