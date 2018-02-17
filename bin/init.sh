#!/bin/bash
FILE=${1:?file}
echo '{"value":{"binary":true,"code":"'$(base64 $FILE)'"}}' >$FILE.json
curl -XPOST http://localhost:8080/init -d @$FILE.json
