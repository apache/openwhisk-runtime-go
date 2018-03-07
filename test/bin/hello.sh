#!/bin/bash
echo '{"openwhisk":1}'
while read line
do
   name="$(echo $line | jq -r .name)"
   logger -s "name=$name" 
   hello="Hello, $name"
   logger -s "sent response"
   echo '{"hello":"'$hello'"}'
done

