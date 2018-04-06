#!/bin/bash
cd "$(dirname $0)"
test -e exec || GOARCH=amd64 GOOS=linux go build -o exec exec.go
test -e exec.zip || zip -q -r exec.zip exec etc dir
cd -
