#!/bin/bash
export T=$PWD
while read -p "> " line
do eval ${line#$}
done
