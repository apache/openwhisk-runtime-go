#!/bin/bash
echo '{"openwhisk":1}'
while true
do read line
   hello="Hello, $(echo $line | jq -r .name)"
   echo '{"hello":"'$hello'"}'
done

