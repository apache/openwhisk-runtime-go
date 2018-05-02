#!/bin/bash
# executable, defaults to main
exec="${1:-main}"
# absolute path of taget dir or file
source="${2:-/src}"
source="$(readlink -f $source)"
dest="${3:-/out}"
dest="$(readlink -f $dest)"
# prepare a compilation dir
compiledir="$(mktemp -d)"
compilefile="$(mktemp)"
mkdir -p "$compiledir/src/action" "$compiledir/src/main" 
# capitalized main function name
main="$(tr '[:lower:]' '[:upper:]' <<< ${exec:0:1})${exec:1}"
# preparing for compilation
if test -d "$source"
# copy all the files unzipped
then cp -rf "$source"/* "$compiledir/src/"
     mkdir "$compiledir/src/action" 2>/dev/null
     cp "$source"/* "$compiledir/src/action/"
# if we have a single file action, copy it
else cp "$source" "$compiledir/src/action/action.go"
fi 
# prepare the main
cat <<EOF >$compiledir/src/main/main.go
package main

import (
	"os"
	"action"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
)

func main() {
	whisk.StartWithArgs(action.$main, os.Args[1:])
}
EOF
# build it
cd "$compiledir"
GOPATH="$GOPATH:$compiledir" go build -i action
GOPATH="$GOPATH:$compiledir" go build -o "$compilefile" main
# if output is a directory use executable name 
if test -d "$dest"
then dest="$dest/$exec"
fi
cp "$compilefile" "$dest"
chmod +x "$dest"
