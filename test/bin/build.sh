#!/bin/sh
FILE=${1:?go file}
OUT=$(basename $FILE)
BIN=${OUT%%.go}
ZIP=${BIN}.zip
go build -i -o bin/$BIN $FILE
GOOS=linux GOARCH=amd64 go build -o exec $FILE
zip zip/$ZIP exec
rm exec
echo "built bin/$BIN zip/$ZIP"
