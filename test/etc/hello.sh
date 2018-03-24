#!/bin/bash
while read line
do
   name="$(echo $line | jq -r .name)"
   if test "$name" == ""
   then exit
   fi
   echo "name=$name" 
   hello="Hello, $name"
   echo '{"hello":"'$hello'"}' >&3
done

