#!/bin/bash
cd "$(dirname $0)"
SRC=${1:?source}
ID=${2:?numbe}
go get github.com/apache/incubator-openwhisk-client-go/whisk
rm -Rvf compile/$ID >/dev/null
rm -Rvf output/$ID >/dev/null
mkdir -p compile/$ID output/$ID
if test -d "$SRC"
then cp -r "$SRC" compile/$ID
else cp $SRC compile/$ID/exec
fi


