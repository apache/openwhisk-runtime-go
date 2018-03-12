#!/bin/bash
DEFAULT='{"name": "Mike"}'
JSON=${1:-$DEFAULT}
DATA='{"value":'$JSON'}'
curl -H "Content-Type: application/json" -XPOST http://localhost:${PORT:-8080}/run -d "$DATA" 2>/dev/null
echo ""
