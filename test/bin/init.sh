#!/bin/bash
FILE=${1:?file}
echo '{"value":{"binary":true,"code":"'$(base64 -w 0 $FILE)'"}}' >$FILE.json
curl -XPOST http://localhost:${PORT:-8080}/init -d @$FILE.json
