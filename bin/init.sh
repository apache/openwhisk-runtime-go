#!/bin/bash
FILE=${1:?file}
echo '{"value":{"binary":true,"code":"'$(base64 $FILE)'"}}' >$FILE.json
curl -XPOST http://localhost:${PORT:-8081}/init -d @$FILE.json
