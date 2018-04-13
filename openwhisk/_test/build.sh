#!/bin/bash
cd "$(dirname $0)"

function build {
  if ! test -e $1
  then cp $1.src $1.go
       GOARCH=amd64 GOOS=linux go build -a -o $1 $1.go
       rm $1.go
  fi
}

function zipit {
  if ! test -e $1
  then
    mkdir $$
    cp $2 $$/$3
    zip -q -j $1 $$/$3
    rm -rf $$
  fi 
}

build exec
test -e exec.zip || zip -q -r exec.zip exec etc dir
build hi
zipit hi.zip hi main
build hello_message
zipit hello_message.zip hello_message main
zipit hello_message1.zip hello_message message
build hello_greeting
zipit hello_greeting.zip hello_greeting main
zipit hello_greeting1.zip hello_greeting greeting


