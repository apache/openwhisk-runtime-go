#!/bin/bash
while read line
do
   name="$(echo $line | jq -r .name)" 
   if [ "$name" == "*" ]
   then echo "Goodbye!" >&2 
        exit 0
   fi
   echo msg="hello $name"
   echo '{"hello": "'$name'"}' >&3
done

