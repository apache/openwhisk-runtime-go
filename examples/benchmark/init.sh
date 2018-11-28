#!/bin/bash
INIT=${1:?action}
jq -n --rawfile file $INIT '{ "value": {"main":"main", "code":$file}}' >$INIT.json
curl -XPOST -H "Content-Type: application/json" http://localhost:8080/init -d @$INIT.json
