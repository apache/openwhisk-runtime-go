#!/bin/bash
while read line
do
   name="$(echo $line | jq -r .name)" 
   echo msg="hello $name"
   echo '{"hello": "'$name'"}' >&3
done

