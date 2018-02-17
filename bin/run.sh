#!/bin/bash
DEFAULT='{"name": "Mike"}'
JSON=${1:-$DEFAULT}
DATA='{"value":'$JSON'}'
curl -XPOST http://localhost:8080/run -d "$DATA"
