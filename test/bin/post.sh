#!/bin/bash
FILE=${1:?file}
curl -H "Content-Type: application/json" -XPOST -w "%{http_code}\n" http://localhost:${PORT:-8080}/init -d @$FILE 2>/dev/null
