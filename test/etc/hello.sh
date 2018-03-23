#!/bin/bash
while read line
do
   name="$(echo $line | jq -r .name)"
   echo "name=$name" 
   hello="Hello, $name"
   echo '{"hello":"'$hello'"}' >&3
done

